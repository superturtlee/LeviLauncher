![LeviLauncher](https://socialify.git.ci/LiteLDev/LeviLauncher/image?custom_language=Go&description=1&font=Inter&forks=1&issues=1&language=1&logo=https%3A%2F%2Fgithub.com%2FLiteLDev%2FLeviLauncher%2Fblob%2Fmain%2Fbuild%2Fappicon.png%3Fraw%3Dtrue&name=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Auto)


<p align="center">
  <a href="https://discord.gg/v5R5P4vRZk"><img alt="Discord" src="https://img.shields.io/discord/849252980430864384?style=for-the-badge&logo=discord"></a>
  <a href="https://qm.qq.com/q/1z791rJgJG"><img alt="QQ 群 458083875" src="https://img.shields.io/badge/458083875-red?style=for-the-badge&logo=qq"></a>
</p>

<p align="center">
  <sup>🌐 语言: <a href="./README.md">English</a> • <b>中文</b></sup>
</p>

一个专为 Minecraft Bedrock Edition GDK 版本打造的命令行下载工具。

仅保留游戏下载与解压功能的精简版本。

## 项目状态

- ✅ 已简化为命令行工具，去除所有GUI功能。

## 适用范围

- 面向 Minecraft GDK（Windows）。需要合法的正版游戏副本。

## 主要功能

- ✅ 版本列表查询：`launcher list` 获取所有可用版本
- ✅ 版本下载与解压：`launcher download <版本号>` 下载并自动解压

## 系统要求

- 操作系统：Windows 10/11
- Go `1.24+`（仅用于编译）

## 使用方法

### 编译

```bash
go build -o launcher.exe main.go
```

### 查看可用版本

```bash
launcher.exe list
```

### 下载指定版本

```bash
launcher.exe download 1.21.50.07
```

## 命令说明

- `launcher list` - 列出所有可用的 Minecraft Bedrock 版本（包括正式版和预览版）
- `launcher download <版本号>` - 下载并解压指定版本到 versions 目录

## 目录结构

- `internal/`：核心功能模块（下载、解压、版本管理）
- `main.go`：命令行工具入口

## 许可证

Copyright © 2025 LeviMC, All rights reserved.

本项目的非闭源部分以 LGPL-3.0 许可发布；详情请参见仓库中的 [COPYING](COPYING) 和 [COPYING.LESSER](COPYING.LESSER) 文件。
