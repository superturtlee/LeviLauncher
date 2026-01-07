package mcservice

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/liteldev/LeviLauncher/internal/extractor"
	"github.com/liteldev/LeviLauncher/internal/utils"
	"github.com/liteldev/LeviLauncher/internal/vcruntime"
)

type VersionStatus struct {
	Version      string `json:"version"`
	IsInstalled  bool   `json:"isInstalled"`
	IsDownloaded bool   `json:"isDownloaded"`
	Type         string `json:"type"`
}

func InstallExtractMsixvc(ctx context.Context, name string, folderName string, isPreview bool) string {
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
		// CLI mode: errors are returned as string codes for the caller to handle and display
		if strings.TrimSpace(msg) == "" {
			msg = "ERR_APPX_INSTALL_FAILED"
		}
		_ = os.RemoveAll(outDir)
		return msg
	}
	_ = vcruntime.EnsureForVersion(ctx, outDir)
	//_ = preloader.EnsureForVersion(ctx, outDir)
	//_ = peeditor.EnsureForVersion(ctx, outDir)
	//_ = peeditor.RunForVersion(ctx, outDir)
	// CLI mode: success is indicated by returning an empty string
	return ""
}

func ResolveDownloadedMsixvc(version string, versionType string) string {
	dir, err := utils.GetInstallerDir()
	if err != nil || strings.TrimSpace(dir) == "" {
		return ""
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		lower := strings.ToLower(name)
		if !strings.HasSuffix(lower, ".msixvc") {
			continue
		}
		ext := filepath.Ext(name)
		if strings.ToLower(ext) == ".msixvc" {
			name = name[:len(name)-len(ext)]
		} else {
			name = strings.TrimSuffix(name, ".msixvc")
		}
		b := strings.TrimSpace(name)
		v := strings.TrimSpace(version)
		bl := strings.ToLower(b)
		vl := strings.ToLower(v)
		if vl == bl {
			return name
		}
	}
	return ""
}

func DeleteDownloadedMsixvc(version string, versionType string) string {
	name := strings.TrimSpace(ResolveDownloadedMsixvc(version, versionType))
	if name == "" {
		return "ERR_MSIXVC_NOT_FOUND"
	}
	dir, err := utils.GetInstallerDir()
	if err != nil || strings.TrimSpace(dir) == "" {
		return "ERR_ACCESS_INSTALLERS_DIR"
	}
	path := filepath.Join(dir, name+".msixvc")
	if !utils.FileExists(path) {
		alt := filepath.Join(dir, name)
		if utils.FileExists(alt) {
			path = alt
		}
	}
	if !utils.FileExists(path) {
		return "ERR_MSIXVC_NOT_FOUND"
	}
	if err := os.Remove(path); err != nil {
		return "ERR_WRITE_TARGET"
	}
	return ""
}

func GetInstallerDir() string {
	dir, err := utils.GetInstallerDir()
	if err != nil {
		return ""
	}
	return dir
}

func GetVersionsDir() string {
	dir, err := utils.GetVersionsDir()
	if err != nil {
		return ""
	}
	return dir
}

func GetVersionStatus(version string, versionType string) VersionStatus {
	status := VersionStatus{Version: version, Type: versionType, IsInstalled: false, IsDownloaded: false}
	if name := ResolveDownloadedMsixvc(version, versionType); strings.TrimSpace(name) != "" {
		status.IsDownloaded = true
	}
	return status
}

func GetAllVersionsStatus(versionsList []map[string]interface{}) []VersionStatus {
	var results []VersionStatus
	for _, versionData := range versionsList {
		version, ok := versionData["version"].(string)
		if !ok {
			version, ok = versionData["short"].(string)
			if !ok {
				continue
			}
		}
		versionType, ok := versionData["type"].(string)
		if !ok {
			versionType = "release"
		}
		status := GetVersionStatus(version, versionType)
		results = append(results, status)
	}
	return results
}
