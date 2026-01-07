# LeviLauncher CLI Usage Examples

## Prerequisites

Before using the launcher, ensure:
1. You own Minecraft Bedrock Edition
2. You're logged into Windows with a Microsoft account that has the Minecraft license
3. You've installed Minecraft from the Microsoft Store at least once (this validates your license with Windows)

The launcher verifies your Minecraft license during extraction using Windows authentication.

## Building the CLI Tool

On Windows with Go 1.24+ installed:

```bash
go build -o launcher.exe main.go
```

## List Available Versions

To see all available Minecraft Bedrock versions:

```bash
launcher.exe list
```

Output example:
```
Fetching available Minecraft versions...

Available Minecraft Bedrock Versions:
======================================

Release Versions:
  - 1.21.50.07
  - 1.21.44.01
  - 1.21.43.01
  ...

Preview Versions:
  - 1.21.60.23
  - 1.21.60.21
  ...
```

## Download a Specific Version

To download and extract a specific version:

```bash
launcher.exe download 1.21.50.07
```

This will:
1. Fetch the version list
2. Find the download URL for the specified version
3. Download the .msixvc file to the installer directory
4. Extract the game files to `versions/<version>/`
5. Ensure all required runtime dependencies are installed

Progress output example:
```
Starting download for version: 1.21.50.07
Found release version with URL: https://...
Downloading to: C:\...\installers\1.21.50.07.msixvc
Progress: 15.2% (128000000/842856960 bytes)
Progress: 30.5% (257000000/842856960 bytes)
...
Download complete: 842856960 bytes
Extracting files...

Successfully downloaded and extracted version 1.21.50.07 to folder: 1.21.50.07
```

## Directory Structure

After running the launcher, you'll find:

```
launcher.exe location/
├── data/                # Launcher data directory
│   ├── installers/      # Downloaded .msixvc files
│   │   └── 1.21.50.07.msixvc
│   └── versions/        # Extracted game files
│       └── 1.21.50.07/
│           └── Minecraft.Windows.exe
```

## Error Handling

The tool will display helpful error messages:

- If version not found:
  ```
  Version X.XX.XX.XX not found
  ```

- If Minecraft license not detected:
  ```
  Extraction failed: Minecraft authorization required
  This tool requires a valid Minecraft license to extract and register the game.
  Please ensure you:
    1. Have purchased Minecraft Bedrock Edition
    2. Are logged into Windows with a Microsoft account that owns Minecraft
    3. Have previously installed Minecraft from the Microsoft Store at least once
  ```

- If download fails:
  ```
  Failed to download: connection timeout
  ```

- If extraction fails:
  ```
  Failed to extract: ERR_APPX_INSTALL_FAILED
  ```

## Tips

1. **Resume Downloads**: The tool supports resuming interrupted downloads automatically
2. **Storage**: Make sure you have enough disk space (each version is ~800MB-1GB)
3. **Network**: A stable internet connection is recommended for large downloads
4. **Permissions**: The tool may need admin rights to create directories in some locations
