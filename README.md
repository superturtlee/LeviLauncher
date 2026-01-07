![LeviLauncher](https://socialify.git.ci/LiteLDev/LeviLauncher/image?custom_language=Go&description=1&font=Inter&forks=1&issues=1&language=1&logo=https%3A%2F%2Fgithub.com%2FLiteLDev%2FLeviLauncher%2Fblob%2Fmain%2Fbuild%2Fappicon.png%3Fraw%3Dtrue&name=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Auto)

<p align="center">
  <a href="https://discord.gg/v5R5P4vRZk"><img alt="Discord" src="https://img.shields.io/discord/849252980430864384?style=for-the-badge&logo=discord"></a>
  <a href="https://qm.qq.com/q/1z791rJgJG"><img alt="QQ Group 458083875" src="https://img.shields.io/badge/458083875-red?style=for-the-badge&logo=qq"></a>
</p>

<p align="center">
  <sup>🌐 Language: <b>English</b> • <a href="./README.zh-CN.md">中文</a></sup>
</p>

A command-line tool for downloading Minecraft Bedrock Edition (GDK) on Windows.

This is a simplified version with only game download and extraction functionality.

## Project Status
- ✅ Simplified to CLI tool with all GUI features removed.

## Scope
- Targets Minecraft GDK (Windows). Requires a legitimate licensed game copy.
- **Important**: You must own Minecraft Bedrock Edition and be logged into Windows with a Microsoft account that has the license. The tool uses Windows authentication to verify ownership during extraction.

## Features

- ✅ Version list query: `launcher list` to get all available versions
- ✅ Version download and extraction: `launcher download <version>` to download and auto-extract

## Requirements
- OS: Windows 10/11
- **Minecraft Bedrock Edition license** - Must be logged into Windows with a Microsoft account that owns the game
- Go `1.24+` (for compilation only)

## Usage

### Build

```bash
go build -o launcher.exe main.go
```

### List Available Versions

```bash
launcher.exe list
```

### Download Specific Version

```bash
launcher.exe download 1.21.50.07
```

## Commands

- `launcher list` - List all available Minecraft Bedrock versions (both release and preview)
- `launcher download <version>` - Download and extract the specified version to the versions directory

## Structure
- `internal/`: Core functionality modules (download, extraction, version management)
- `main.go`: CLI tool entry point

## License

Copyright © 2025 LeviMC, All rights reserved.

This project is licensed under the LGPL-3.0 License for its non-closed source parts - see the [COPYING](COPYING) and [COPYING.LESSER](COPYING.LESSER) files for details.

