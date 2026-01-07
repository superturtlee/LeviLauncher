package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/corpix/uarand"
	"github.com/liteldev/LeviLauncher/internal/extractor"
	"github.com/liteldev/LeviLauncher/internal/mcservice"
	"github.com/liteldev/LeviLauncher/internal/utils"
	"github.com/liteldev/LeviLauncher/internal/vcruntime"
)

const (
	progressUpdateInterval = 5 * 1024 * 1024 // 5MB
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	extractor.Init()

	command := os.Args[1]
	switch command {
	case "list":
		handleList()
	case "download":
		if len(os.Args) < 3 {
			fmt.Println("Error: download command requires version argument")
			fmt.Println("Usage: launcher download <version>")
			os.Exit(1)
		}
		handleDownload(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("LeviLauncher CLI - Minecraft Bedrock Downloader")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  launcher list              - List available Minecraft versions")
	fmt.Println("  launcher download <ver>    - Download and extract a specific version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  launcher list")
	fmt.Println("  launcher download 1.21.50.07")
}

func handleList() {
	fmt.Println("Fetching available Minecraft versions...")
	versions := mcservice.FetchHistoricalVersions(false)

	if len(versions) == 0 {
		fmt.Println("Failed to fetch version list")
		os.Exit(1)
	}

	fmt.Println("\nAvailable Minecraft Bedrock Versions:")
	fmt.Println("======================================")

	// Parse and display versions
	if release, ok := versions["release"].([]interface{}); ok {
		fmt.Println("\nRelease Versions:")
		for _, v := range release {
			if vmap, ok := v.(map[string]interface{}); ok {
				if version, ok := vmap["version"].(string); ok {
					fmt.Printf("  - %s\n", version)
				}
			}
		}
	}

	if preview, ok := versions["preview"].([]interface{}); ok {
		fmt.Println("\nPreview Versions:")
		for _, v := range preview {
			if vmap, ok := v.(map[string]interface{}); ok {
				if version, ok := vmap["version"].(string); ok {
					fmt.Printf("  - %s\n", version)
				}
			}
		}
	}
}

func handleDownload(version string) {
	fmt.Printf("Starting download for version: %s\n", version)

	// Fetch version list to find download URL
	versions := mcservice.FetchHistoricalVersions(false)
	if len(versions) == 0 {
		fmt.Println("Failed to fetch version list")
		os.Exit(1)
	}

	var downloadURL string
	var versionType string = "release"
	var found bool

	// Search in release versions
	if release, ok := versions["release"].([]interface{}); ok {
		for _, v := range release {
			if vmap, ok := v.(map[string]interface{}); ok {
				if ver, ok := vmap["version"].(string); ok && ver == version {
					if urlVal, ok := vmap["url"].(string); ok {
						downloadURL = urlVal
						found = true
						break
					}
				}
			}
		}
	}

	// Search in preview versions if not found
	if !found {
		if preview, ok := versions["preview"].([]interface{}); ok {
			for _, v := range preview {
				if vmap, ok := v.(map[string]interface{}); ok {
					if ver, ok := vmap["version"].(string); ok && ver == version {
						if urlVal, ok := vmap["url"].(string); ok {
							downloadURL = urlVal
							versionType = "preview"
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found || downloadURL == "" {
		fmt.Printf("Version %s not found\n", version)
		os.Exit(1)
	}

	fmt.Printf("Found %s version with URL: %s\n", versionType, downloadURL)

	// Get installer directory
	installerDir, err := utils.GetInstallerDir()
	if err != nil {
		fmt.Printf("Failed to get installer directory: %v\n", err)
		os.Exit(1)
	}

	// Derive filename from URL
	filename := deriveFilename(downloadURL)
	destPath := filepath.Join(installerDir, filename)

	// Download the file
	fmt.Printf("Downloading to: %s\n", destPath)
	err = downloadFile(downloadURL, destPath)
	if err != nil {
		fmt.Printf("Failed to download: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Download completed successfully!")

	// Extract the downloaded file
	fmt.Println("Extracting files...")
	folderName := version
	if versionType == "preview" {
		folderName = version + "_preview"
	}

	ctx := context.Background()
	isPreview := versionType == "preview"
	errMsg := installExtractMsixvc(ctx, filename, folderName, isPreview)
	if errMsg != "" {
		fmt.Printf("Failed to extract: %s\n", errMsg)
		os.Exit(1)
	}

	fmt.Printf("\nSuccessfully downloaded and extracted version %s to folder: %s\n", version, folderName)
}

func deriveFilename(rawURL string) string {
	fname := "download.msixvc"
	if u, e := url.Parse(rawURL); e == nil {
		if v := u.Query().Get("filename"); v != "" {
			fname = ensureMsixvcFilename(v)
		} else {
			parts := strings.Split(u.Path, "/")
			if len(parts) > 0 && parts[len(parts)-1] != "" {
				fname = ensureMsixvcFilename(parts[len(parts)-1])
			}
		}
	}
	return fname
}

func ensureMsixvcFilename(name string) string {
	n := strings.TrimSpace(name)
	if n == "" {
		return "download.msixvc"
	}
	n = strings.ReplaceAll(n, "/", "_")
	n = strings.ReplaceAll(n, "\\", "_")
	lower := strings.ToLower(n)
	if !strings.HasSuffix(lower, ".msixvc") {
		n += ".msixvc"
	}
	return n
}

func downloadFile(rawURL, destPath string) error {
	// Remove filename parameter from URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Del("filename")
	u.RawQuery = q.Encode()
	cleanURL := u.String()

	// Create the directory if it doesn't exist
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file already exists and get its size
	var downloaded int64
	if fi, err := os.Stat(destPath); err == nil {
		downloaded = fi.Size()
		fmt.Printf("Resuming download from %d bytes\n", downloaded)
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", cleanURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", uarand.GetRandom())
	if downloaded > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", downloaded))
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Get total file size
	var total int64
	if resp.StatusCode == http.StatusPartialContent {
		total = downloaded + resp.ContentLength
	} else {
		total = resp.ContentLength
		downloaded = 0 // Reset if server doesn't support resume
	}

	// Open file for writing
	var flag int
	if downloaded > 0 && resp.StatusCode == http.StatusPartialContent {
		flag = os.O_WRONLY | os.O_APPEND
	} else {
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	
	out, err := os.OpenFile(destPath, flag, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download with progress indicator
	buf := make([]byte, 32*1024)
	lastPrint := int64(0)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				return fmt.Errorf("failed to write file: %w", werr)
			}
			downloaded += int64(n)
			
			// Print progress every 5MB
			if downloaded-lastPrint >= progressUpdateInterval || downloaded == total {
				percentage := float64(downloaded) / float64(total) * 100
				fmt.Printf("Progress: %.1f%% (%d/%d bytes)\n", percentage, downloaded, total)
				lastPrint = downloaded
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
	}

	fmt.Printf("Download complete: %d bytes\n", downloaded)
	return nil
}

func installExtractMsixvc(ctx context.Context, name string, folderName string, isPreview bool) string {
	n := strings.TrimSpace(name)
	if n == "" {
		return "ERR_MSIXVC_NOT_SPECIFIED"
	}
	inPath := n
	if !filepath.IsAbs(inPath) {
		if dir, err := utils.GetInstallerDir(); err == nil && dir != "" {
			// Ensure the filename has .msixvc extension before joining
			if !strings.HasSuffix(strings.ToLower(inPath), ".msixvc") {
				inPath += ".msixvc"
			}
			inPath = filepath.Join(dir, inPath)
		}
	}
	if !utils.FileExists(inPath) {
		return "ERR_MSIXVC_NOT_FOUND"
	}
	vdir, err := utils.GetVersionsDir()
	if err != nil || strings.TrimSpace(vdir) == "" {
		return "ERR_ACCESS_VERSIONS_DIR"
	}
	outDir := filepath.Join(vdir, strings.TrimSpace(folderName))
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "ERR_CREATE_TARGET_DIR"
	}

	rc, msg := extractor.Get(inPath, outDir)
	if rc != 0 {
		if strings.TrimSpace(msg) == "" {
			msg = "ERR_APPX_INSTALL_FAILED"
		}
		_ = os.RemoveAll(outDir)
		return msg
	}
	_ = vcruntime.EnsureForVersion(ctx, outDir)
	return ""
}
