# Golang for Java Developers

This is a workshop for Java developers to learn how to do cloud native Goland in a hands-on way. It contains 11 labs starting with the fundamentals and then moving on to cloud native patterns.

# Teacher Instructions

This is the main repo for Improving's Golang Workshop. Feel free to clone and push this to a different repo if needed to teach the course. 

The teacher must:

- Provide a PDF of the slides to the class. Master slides here: https://improving-my.sharepoint.com/:p:/p/mark_soule/IQAAdxenPQoQTLIy7a_chNJYAa6h3OOR0joghu_giMXQ3O8?e=0QBZOx
- Ensure that the class will be able to setup all prePrerequisites before the class. If not, adapt these labs accordingly. 
- Each lab has a `solution` folder in the `solution` branch. If any labs are modified, the solution should also be modified.
- The `common` folder is for common utils and students will not need to modify.

---

# Student Instructions

Students should clone this repository and checkout the `students-start-here` branch:

```
git checkout -b students-start-here
```

Ensure that all Prerequisites below are setup then await directions from your instructor.

Setup prequesites listed below then start with lab01. Some labs build off of previous labs and those labs will be indicated as such. Starter files are provided for all labs.

## Prerequisites

- IDE with golang support. Options include:
	- Intellij Ultimate + Golang plugin.
	- JetBrains GoLand (functionally the same as above).
	- Any fork of VSCode + Golang plugin.
- Git
- Make
- Golang Compiler version 1.22 or higher.
- Read access to the lab repository. Options discussed:
	- I publish to Improving owned public Github repository.
	- I zip up the files and send them to be hosted in your repository.
- Docker and Docker Compose.
- Access to Dockerhub or a suitable mirror.
	- Some labs require running docker containers for external systems (i.e. prometheus, database)
- Both HTTP and gRPC clients
	- Any suitable tools will work, including curl, Postman, grpcurl, etc.
- Golang Protocol Buffer Libraries
	- protoc-gen-go - installed via: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
	- protoc-gen-go-grpc - installed via: `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
- Protocol Buffer Compiler (protoc).
  - https://protobuf.dev/installation/

## Installing Go

1. Download from: https://go.dev/dl/
2. Run the installer for your operating system
3. Verify installation: `go version`

### Environment Setup

After installation, configure your Go environment variables:

- GOROOT: The location where Go is installed (usually set automatically by the installer). Verify with `go env GOROOT`
- GOPATH: Your Go workspace directory where downloaded modules and compiled binaries are stored. Default is `$HOME/go` (macOS/Linux) or `%USERPROFILE%\go` (Windows). Verify with `go env GOPATH`
- PATH: Ensure `$GOPATH/bin` is in your PATH to access installed Go tools
	- macOS/Linux: Add to shell profile: `export PATH=$PATH:$GOPATH/bin`
	- Windows (PowerShell): `[System.Environment]::SetEnvironmentVariable('PATH', "$env:PATH;$env:GOPATH\bin", 'User')`

### Windows-specific notes

- For Windows Subsystem for Linux (WSL): Install Go within your WSL distribution using the Linux installation instructions
- For native Windows (PowerShell/CMD): Use the Windows installer (.msi)

## Installing Protocol Buffer Compiler (protoc)

- macOS: `brew install protobuf`
- Linux: Download from https://github.com/protocolbuffers/protobuf/releases
- Windows (WSL): Follow Linux instructions within WSL
- Windows (native): Download Windows binary from releases page and add to PATH

Verify installation: `protoc --version`
