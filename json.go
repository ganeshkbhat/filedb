package gojson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "os"
)

// ReadFile reads a file from disk and returns its contents as a byte slice.
func ReadFile(filePath string) ([]byte, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

// WriteFile writes the given string data to the specified file path.
func WriteFile(filePath string, data string) error {
	// Convert the string data to a byte slice
	dataBytes := []byte(data)

	// Write the data to the file
	err := ioutil.WriteFile(filePath, dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// JsonMarshal marshals an object to a JSON byte slice.
func JsonMarshal(obj interface{}) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// JsonUnmarshal unmarshals a JSON byte slice into an object.
func JsonUnmarshal(jsonBytes []byte, obj interface{}) error {
	err := json.Unmarshal(jsonBytes, obj)
	if err != nil {
		return err
	}
	return nil
}

// YamlToJson converts a YAML byte slice to a JSON string.
func readJsonFile(filepath string) string {
	// Path to the YAML file
	jsonFilePath := filepath

	// Read the YAML file
	jsonData, err := ReadFile(jsonFilePath)
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonData)

	// Print the YAML string
	fmt.Println(jsonString)

	return jsonString
}

func Usage(filepath string, data string) {
	// // Example file path
	// filePath := "example.txt"
	// //
	// // Example data to write to the file
	// data := "Hello, world!"

	// Write the data to the file
	err := WriteFile(filepath, data)
	if err != nil {
		panic(err)
	}

	// Read the data from the file
	fileData, err := ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	// Print the data from the file
	fmt.Println(string(fileData))
}

