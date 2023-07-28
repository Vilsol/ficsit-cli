# ficsit-cli [![push](https://github.com/Vilsol/ficsit-cli/actions/workflows/push.yaml/badge.svg)](https://github.com/Vilsol/ficsit-cli/actions/workflows/push.yaml) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vilsol/ficsit-cli) ![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/vilsol/ficsit-cli) [![GitHub license](https://img.shields.io/github/license/Vilsol/ficsit-cli)](https://github.com/Vilsol/ficsit-cli/blob/master/LICENSE) ![GitHub all releases](https://img.shields.io/github/downloads/vilsol/ficsit-cli/total)

A CLI tool for managing mods for the game Satisfactory

## Installation

### Windows

Download the appropriate `.exe` for your CPU architecture.

* [AMD64 (64-bit)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_windows_amd64.exe)
* [386 (32-bit)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_windows_386.exe)
* [ARM64 (64-bit ARM)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_windows_arm64.exe)
* [ARMv7 (32-bit ARM)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_windows_armv7.exe)

### Linux

#### Arch

A package is published to AUR under the name [`ficsit-cli-bin`](https://aur.archlinux.org/packages/ficsit-cli-bin)

```shell
yay -S ficsit-cli-bin
```

#### Debian (inc. Ubuntu, Mint, PopOS!, etc)

Download the appropriate `.deb` for your CPU architecture.

* [AMD64 (64-bit)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_amd64.deb)
* [386 (32-bit)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_386.deb)
* [ARM64 (64-bit ARM)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_arm64.deb)
* [ARMv7 (32-bit ARM)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_armv7.deb)
* [PowerPC64](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_ppc64le.deb)

#### Fedora

Download the appropriate `.rpm` for your CPU architecture.

* [AMD64 (64-bit)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_amd64.rpm)
* [386 (32-bit)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_386.rpm)
* [ARM64 (64-bit ARM)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_arm64.rpm)
* [ARMv7 (32-bit ARM)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_armv7.rpm)
* [PowerPC64](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_ppc64le.rpm)

#### Alpine

Download the appropriate `.apk` for your CPU architecture.

* [AMD64 (64-bit)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_amd64.apk)
* [386 (32-bit)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_386.apk)
* [ARM64 (64-bit ARM)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_arm64.apk)
* [ARMv7 (32-bit ARM)](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_armv7.apk)
* [PowerPC64](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_linux_ppc64le.apk)

### macOS

Download the "all" build [here](https://github.com/Vilsol/ficsit-cli/releases/latest/download/ficsit_darwin_all).

## Usage

### Interactive CLI

To launch the interactive CLI, run the executable without any arguments.

### Command Line

Run `ficsit help` to see a list of available commands.

## Troubleshooting

* Profile and installation records are located in `%APPDATA%\ficsit\`
* Downloads are cached in `%LOCALAPPDATA%\ficsit\downloadCache\`

## Development

### Dependencies

* [Go 1.18](https://go.dev/doc/install)
* IDE of Choice. Goland or VSCode suggested.

### Code Generation

If you update any of the GraphQL queries, run this to update generated code:

```bash
(echo "y") | npx graphqurl https://api.ficsit.app/v2/query --introspect -H 'content-type: application/json' > schema.graphql
go generate -tags tools -x ./...
```

## Building

```bash
go build
```

Will produce `ficsit-cli.exe` in the repo root directory.

### Linting

If the CI gives you an error like this:

```bash
Error: File is not `gci`-ed with --skip-generated -s standard,default,prefix(github.com/satisfactorymodding/ficsit-cli),blank,dot (gci)
```

Install and run `gofumpt`.

```bash
go install mvdan.cc/gofumpt@latest
```

```bash
gofumpt -l -w .
```

You may also need to run GCI. If you're on Windows, watch out, `gci` is also a powershell shortened command, so run this from Bash. Install via:

```bash
go install github.com/daixiang0/gci@latest
```

Example command to run based on the above error:

```bash
gci.exe write --skip-generated -s standard -s default -s 'prefix(github.com/satisfactorymodding/ficsit-cli)' -s blank -s dot .
```
