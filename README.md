# XRPL Encoder

`xrpl-encoder` is a comprehensive tool for encoding and decoding JSON or HEX data specific to XRPL. Whether you have a single input, a file, or an entire directory of files, this tool streamlines the process.

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

Make sure you have Go installed on your system. Also, this tool relies on the [`github.com/xyield/xrpl-go/binary-codec`](https://github.com/xyield/xrpl-go/binary-codec) package, so ensure you have this dependency set up.

## Installation

1. Clone this repository to your local machine.
2. Navigate to the project directory.
3. Build the application with the following command:

```bash
go build
```
## Usage

### Command Line Flags

You can use the tool with the following command-line flags:

- `-data`: Directly provide HEX or JSON data as input.
- `-file`: Provide the path to a file containing HEX or JSON data.
- `-batch`: Provide the path to a directory containing multiple HEX or JSON files.
- `-help`: Show the help message.

### Interactive Mode

If you run `xrpl-encoder` without any flags, it enters interactive mode. Here are the steps:

1. Run the tool:
```bash
./xrpl-encoder
```

## Output
The tool processes your input and provides either the encoded HEX or decoded JSON output. After processing, it offers you an option to save the output to a file. If you choose to save, it will either use a default naming scheme or a custom name you provide.

## Contributing
Your contributions are always welcome! Please feel free to submit pull requests, open issues, or provide feedback.

## License
This project is licensed under the MIT License. See LICENSE for more details.
