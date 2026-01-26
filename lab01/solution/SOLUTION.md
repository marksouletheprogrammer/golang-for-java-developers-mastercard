# Lab 1 Solution

## How to Run

Run directly without building:
```bash
cd lab01/solution
go run main.go
```

Build and run the binary:
```bash
cd lab01/solution
go build
./hello-go
```

Build with custom output name:
```bash
go build -o my-program
./my-program
```

## Part 4: Adding External Dependencies

Add the logrus dependency:
```bash
go get github.com/sirupsen/logrus
```

After running `go get`, examine the files:

**go.mod** will show:
```
module hello-go

go 1.22

require github.com/sirupsen/logrus v1.9.3

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
```

Note the two `require` blocks:
- **Direct dependency**: `github.com/sirupsen/logrus` - the package you explicitly added
- **Indirect dependency**: `golang.org/x/sys` - a transitive dependency (logrus needs it)

**go.sum** will contain multiple entries:
- Each dependency gets two lines: one for the module hash (`h1:...`) and one for the go.mod hash (`/go.mod`)
- You'll see entries for logrus, golang.org/x/sys, and test dependencies
- This file ensures reproducible builds by verifying checksums

**Expected output:**
```
INFO[0000] Welcome to Go!
INFO[0000] Current timestamp: 2026-01-16 15:45:23 CST
INFO[0000] ISO 8601 format: 2026-01-16T15:45:23-06:00
```

## Part 5: go run vs go build

**go run**: Compiles and runs the program in one step. The binary is created in a temporary location and executed. Fast for development and testing.

**go build**: Compiles the program and creates a binary executable in the current directory. The binary can be distributed and run independently without Go installed. Use this for production deployments.

The binary name defaults to the directory name (hello-go) or can be specified with `-o` flag.
