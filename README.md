![Go Report Card](https://goreportcard.com/badge/github.com/xyield/xrpl-encoder)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/xyield/xrpl-encoder)
![GitHub](https://img.shields.io/github/license/xyield/xrpl-encoder)
# XRPL Encoder


`xrpl-encoder` is a tool, written in Go, for encoding and decoding JSON or HEX data specific to XRP Ledger transactions. Whether you have a single input, a file, or an entire directory of files, this tool streamlines the process.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Command Line Flags](#command-line-flags)
  - [Interactive Mode](#interactive-mode)
- [Output](#output)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

- ### Make sure you have Go installed 
  Follow the official [Go installation guide](https://golang.org/doc/install)

- ### Dependencies
  This tool relies on the [`github.com/xyield/xrpl-go/binary-codec`](https://github.com/xyield/xrpl-go/binary-codec) package.
  
  However, Go will handle the dependencies for you thanks to the `go.mod` file in the repository.

## Installation

1. Clone this repository to your local machine.
2. Navigate to the project directory.
3. Build the application with the following command:

```bash
go build
```
## Usage

Files and folders to be encoded/decoded should be placed in the `process/` directory.

### Interactive Mode

If you run `xrpl-encoder` without any flags, it enters interactive mode. Here are the steps:

1. Run the tool:
```bash
./xrpl-encoder
```
2. Select an option:
```
WARNING: For very large data entries, you may overload your terminal 
when pasting with Direct Input (Option 1).
Consider using the File Input method (Option 2) for large datasets.

Choose input method:

1. Direct Input
2. File Input
3. Batch Processing (Directory Input)
4. Help
5. Exit
```
For both the File Input and Batch Processing options, you don't need to prepend your input with `process/`.
The tool automatically checks within this directory. 

For instance, if you have a file named example.json in the `process/` directory, 
you should enter `example.json` when prompted. 

If you have a folder with the path `process/examples/`, you should enter `examples`.
If you just want to process all files within the root of the `process/` directory, provide no argument.

### Command Line Flags

You can use the tool with the following command-line flags:

- `-d`: Directly provide HEX or JSON data as input. ```./xrpl-encoder -d 120007220000000024...```
- `-f`: Provide the path to a file containing HEX or JSON data.  ```./xrpl-encoder -f example.json```
- `-b`: Provide the path to a directory with multiple files. ```./xrpl-encoder -b examples```
- `-b`: Process all files in the root of the `process/` directory. ```./xrpl-encoder -b```
- `-h`: Show the help message. ```./xrpl-encoder -h```



## Output
The tool processes your input and provides either the encoded HEX or decoded JSON output. After processing, if you choose to save the output, the tool will use a default naming scheme or a custom name you provide, and store the file in the `process/outputs/` directory.

## Contributing
Your contributions are always welcome! Please feel free to submit pull requests, open issues, or provide feedback.

## License
This project is licensed under the MIT License. See [LICENSE](https://github.com/xyield/xrpl-encoder/LICENSE.txt) for more details.
