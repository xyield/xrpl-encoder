package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	binarycodec "github.com/xyield/xrpl-go/binary-codec"
)

var (
	dataInput  = flag.String("data", "", "Directly provide HEX or JSON data as input.")
	fileInput  = flag.String("file", "", "Provide the path to a file containing HEX or JSON data.")
	helpFlag   = flag.Bool("help", false, "Show help message")
	batchInput = flag.String("batch", "", "Provide the path to a directory containing multiple HEX or JSON files.")
)

func processBatch(directory string) {
	fmt.Println("Processing directory:", directory)

	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No files found in the directory.")
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			fmt.Println("Processing file:", file.Name()) // Print the name of each file being processed

			filePath := filepath.Join(directory, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue // changed from return to continue, so that one file error doesn't stop the whole batch
			}
			processInput(string(content))
		}
	}
}

func displayHelp() {
	fmt.Println(`
Usage: xrpl-encoder [OPTIONS]

Options:
  -data   Directly provide HEX or JSON data as input.
  -file   Provide the path to a file containing HEX or JSON data.
  -batch  Provide the path to a directory containing multiple HEX or JSON files.

To use the tool in interactive mode, just run it without any flags.
    `)
}

func main() {
	flag.Parse()

	// Check for help flag immediately
	if *helpFlag {
		displayHelp()
		return
	}

	// Check if data input flag is provided
	if *dataInput != "" {
		processInput(*dataInput)
		return
	}

	// Check if file input flag is provided
	if *fileInput != "" {
		content, err := os.ReadFile(*fileInput)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		processInput(string(content))
		return
	}

	// Check if batch input flag is provided
	if *batchInput != "" {
		processBatch(*batchInput)
		return
	}

	// If no flags provided, enter the interactive menu loop
	for {
		// Display menu and get user choice
		choice := displayMenu()

		switch choice {
		case 1, 2, 3, 4:
			handleChoice(choice)
		case 5:
			fmt.Println("\nExiting the tool. Goodbye!")
			return
		default:
			fmt.Println("\nInvalid choice!")
		}
	}
}

func displayMenu() int {
	fmt.Println("\nWARNING: For very large data entries, you may overload your terminal when pasting with Direct Input (Option 1).")
	fmt.Println("Consider using the File Input method (Option 2) for large datasets.")
	fmt.Println("\nChoose input method:\n")
	fmt.Println("1. Direct Input")
	fmt.Println("2. File Input")
	fmt.Println("3. Batch Processing (Directory Input)")
	fmt.Println("4. Display Help")
	fmt.Println("5. Exit")

	var choice int
	fmt.Scanln(&choice)

	return choice
}

func pauseAndReturnToMenu() {
	fmt.Println("\nPress Enter to return to the main menu.")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func handleChoice(choice int) {
	switch choice {
	case 1:
		// Direct Input logic
		fmt.Println("\nPlease paste your JSON or HEX data and press Enter twice:", "\n\nInput Data:\n") //nolint
		inputData := readMultiLineInput()
		processInput(inputData)
		pauseAndReturnToMenu()

	case 2:
		// File Input logic
		fmt.Println("\nEnter the path to your input file:")
		var filePath string
		fmt.Scanln(&filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		processInput(string(content))
		pauseAndReturnToMenu()

	case 3:
		// Batch Processing logic
		fmt.Println("Please provide the directory path for batch processing:")
		var dirPath string
		fmt.Scanln(&dirPath)
		processBatch(dirPath)
		pauseAndReturnToMenu()

	case 4:
		displayHelp()
		pauseAndReturnToMenu()
	}
}

func processInput(inputData string) {
	// Trim the data
	inputData = strings.TrimSpace(inputData)

	if len(inputData) == 0 {
		fmt.Println("Error: No input data provided.")
		return
	}

	// Check if the data is surrounded by quotes or backticks and trim them
	firstChar := inputData[0:1]
	lastChar := inputData[len(inputData)-1:]

	switch {
	case (firstChar == "`" && lastChar == "`") || (firstChar == "\"" && lastChar == "\""):
		inputData = inputData[1 : len(inputData)-1]
	case firstChar == "`" || firstChar == "\"":
		inputData = inputData[1:]
	case lastChar == "`" || lastChar == "\"":
		inputData = inputData[:len(inputData)-1]
	}

	outputFileContent := ""

	_, err := hex.DecodeString(inputData)
	if err != nil {
		var jsonInput map[string]interface{}
		err := json.Unmarshal([]byte(inputData), &jsonInput)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		// Process the JSON fields to convert them to the correct types for encoding (e.g. int64 to int)
		jsonInput = processJSONFields(jsonInput).(map[string]interface{})

		encoded, err := binarycodec.Encode(jsonInput)
		if err != nil {
			fmt.Println("error during encoding:", err)
			return
		}

		fmt.Println("\nEncoded Tx Hex:\n\n", encoded)
		outputFileContent = encoded // Save the encoded hex for output

		decoded, err := binarycodec.Decode(encoded)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		jsonOutput, err := json.MarshalIndent(decoded, " ", " ")
		if err != nil {
			fmt.Println("error during JSON conversion:", err)
			return
		}

		fmt.Println("\n\nChecking if Re-decoded Tx Hex matches the original Tx JSON...")
		time.Sleep(1 * time.Second)

		var original map[string]interface{}
		var reDecoded map[string]interface{}
		json.Unmarshal([]byte(inputData), &original)
		json.Unmarshal(jsonOutput, &reDecoded)

		if reflect.DeepEqual(original, reDecoded) {
			fmt.Println("\nSUCCESS ---> Re-decoded Tx JSON matches the original Tx JSON")
		} else {
			fmt.Println("\nFAIL ---> Re-decoded Tx Hex does not match original Tx JSON\nNote: Some fields in the raw JSON won't be encoded because they don't exist in the binary-codec definitions.json, or they are supposed to be omitted from the binary encoding. This is expected behavior.")
			fmt.Println("\nRe-decoded Tx JSON:\n", string(jsonOutput))
		}
	} else {
		hexEncodedTx := inputData
		decoded, err := binarycodec.Decode(hexEncodedTx)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		jsonOutput, err := json.MarshalIndent(decoded, "", "  ")
		if err != nil {
			fmt.Println("error during JSON conversion:", err)
			return
		}

		fmt.Println("\nDecoded Tx Json:\n\n", string(jsonOutput))
		outputFileContent = string(jsonOutput) // Save the decoded JSON for output

		encoded, err := binarycodec.Encode(decoded)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("\n\nChecking if Re-encoded Tx JSON matches the original Tx Hex...")
		time.Sleep(1 * time.Second)

		if encoded == hexEncodedTx {
			fmt.Println("\nSUCCESS ---> Re-encoded Tx JSON matches the original Tx Hex")
		} else {
			fmt.Println("\nFAIL ---> Re-encoded Tx JSON does not match original Tx Hex")
			fmt.Println("\nRe-encoded Tx Hex:\n", encoded)
		}
	}

	if shouldSave, customName := askForFileOutput(); shouldSave {
		writeOutputToFile(outputFileContent, customName)
	}

}

func readMultiLineInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from input:", err)
	}

	return strings.Join(lines, "\n")
}

func processJSONFields(input interface{}) interface{} {
	switch v := input.(type) {
	case float64:
		return int(v)
	case map[string]interface{}:
		for key, value := range v {
			// Check if the key matches one of the target field names
			if key == "Indexes" || key == "Hashes" || key == "Amendments" || key == "Nftokenoffers" {
				// Attempt to convert the value to []interface{}
				if list, ok := value.([]interface{}); ok {
					v[key] = convertInterfaceSliceToStringSlice(list)
				}
			} else {
				// Continue processing nested values
				v[key] = processJSONFields(value)
			}
		}
	case []interface{}:
		for i, value := range v {
			v[i] = processJSONFields(value)
		}
	}
	return input
}

func convertInterfaceSliceToStringSlice(slice []interface{}) []string {
	var stringSlice []string
	for _, val := range slice {
		strVal, ok := val.(string)
		if ok {
			stringSlice = append(stringSlice, strVal)
		}
	}
	return stringSlice
}

func writeOutputToFile(output, customName string) {
	filename := customName

	// Determine the file extension: .json for JSON content and .txt for others
	extension := ".txt"
	if isJSON(output) {
		extension = ".json"
	}

	// If the user hasn't provided a custom filename, use the default naming logic
	if filename == "output" {
		i := 1
		for fileExists(filename + extension) {
			filename = fmt.Sprintf("output%d", i)
			i++
		}
	} else if fileExists(filename + extension) {
		// If the custom name exists, append a number to it until we find an available filename
		i := 1
		for fileExists(filename + fmt.Sprintf("_%d", i) + extension) {
			i++
		}
		filename = filename + fmt.Sprintf("_%d", i)
	}

	// Create the file with the determined name
	filename = filename + extension
	err := os.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Output saved to", filename)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func askForFileOutput() (bool, string) {
	fmt.Println("\nWould you like to save the output to a file? (y/n) or (y filename): ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)

	if strings.ToLower(answer) == "n" {
		return false, ""
	}

	parts := strings.SplitN(answer, " ", 2)
	if len(parts) > 1 && strings.ToLower(parts[0]) == "y" {
		return true, parts[1]
	}

	return strings.ToLower(parts[0]) == "y", "output"
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
