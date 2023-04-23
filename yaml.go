package goyaml

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "os"
	"gopkg.in/yaml.v3"
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


// FromYaml unmarshals a YAML byte slice into an object.
func FromYaml(yamlBytes []byte, obj interface{}) error {
	err := yaml.Unmarshal(yamlBytes, obj)
	if err != nil {
		return err
	}
	return nil
}


// ToYaml marshals an object to a YAML byte slice.
func ToYaml(obj interface{}) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return yamlBytes, nil
}


// JsonToYaml converts a JSON byte slice to YAML format.
func JsonToYaml(jsonBytes []byte) ([]byte, error) {
	// Unmarshal the JSON data into an object
	var obj interface{}
	err := JsonUnmarshal(jsonBytes, &obj)
	if err != nil {
		return nil, err
	}

	// Marshal the object to YAML format
	yamlBytes, err := ToYaml(obj)
	if err != nil {
		return nil, err
	}

	return yamlBytes, nil
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
func YamlToJson(yamlData []byte) (string, error) {
	// Unmarshal the YAML data into an object
	var obj interface{}
	err := FromYaml(yamlData, &obj)
	if err != nil {
		return "", err
	}

	// Marshal the object to JSON format
	jsonBytes, err := JsonMarshal(obj)
	if err != nil {
		return "", err
	}

	// Convert the JSON byte slice to a string
	jsonString := string(jsonBytes)

	return jsonString, nil
}


// YamlToJson converts a YAML byte slice to a JSON string.
func UsageYamlToJson(filepath string) string {
	// Path to the YAML file
	yamlFilePath := filepath

	// Read the YAML file
	yamlData, err := ReadFile(yamlFilePath)
	if err != nil {
		panic(err)
	}

	// Convert the YAML data to a JSON string
	jsonString, err := YamlToJson(yamlData)
	if err != nil {
		panic(err)
	}

	// Print the JSON string
	fmt.Println(jsonString)

	return jsonString
}


// YamlToJson converts a YAML byte slice to a JSON string.
func UsageJsonToYaml(filepath string) string {
	// Path to the YAML file
	jsonFilePath := filepath

	// Read the YAML file
	jsonData, err := ReadFile(jsonFilePath)
	if err != nil {
		panic(err)
	}

	// Convert the JSON string back to YAML
	yamlBytes, err := JsonToYaml([]byte(jsonData))
	if err != nil {
		panic(err)
	}

	yamlString := string(yamlBytes)

	// Print the YAML string
	fmt.Println(yamlString)

	return yamlString
}


// YamlJsonToFro YAML <=> JSON
func UsageYamlJsonToFro() {
	//
	// Usage for Main
	//
	
	// YAML data
	yamlData := []byte(`name: John Smith
	age: 35
	address:
	street: "123 Main St"
	city: Anytown
	state: CA
	zip: 12345`)

	// Convert the YAML data to a JSON string
	jsonString, err := YamlToJson(yamlData)
	if err != nil {
		panic(err)
	}

	// Print the JSON string
	fmt.Println(jsonString)

	// Convert the JSON string back to YAML
	yamlBytes, err := JsonToYaml([]byte(jsonString))
	if err != nil {
		panic(err)
	}

	// Print the YAML string
	fmt.Println(string(yamlBytes))
}
