# asunder

> Asunder, Sweet and Other Distress

[![Latest Release](https://img.shields.io/github/release/x6r/asunder.svg)](https://github.com/x6r/asunder/releases)
[![Build Status](https://img.shields.io/github/workflow/status/x6r/asunder/build?logo=github)](https://github.com/x6r/asunder/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/x6r/asunder)](https://goreportcard.com/report/github.com/x6r/asunder)

asunder is a small pretty command-line TOTP manager.

![scrot](assets/scrot.png)

## Installation

### Binaries

Download a binary from the [releases](https://github.com/x6r/asunder/releases)
page.

### Build from source

Go 1.16 or higher required. ([install instructions](https://golang.org/doc/install.html))

    go install github.com/x6r/asunder@latest

## Usage

The first time you run asunder you will be asked to create a master password.
After that, everytime asunder is launched you will be prompted for that password.

Add entries:

    asunder add

Start asunder:

    asunder

Delete entries:

    asunder delete

## TODO

- [ ] Export database
- [ ] Import from other services

## License

[OSL-3.0](LICENSE)
