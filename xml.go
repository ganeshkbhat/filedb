package goxml

import (
	"encoding/json"
	"encoding/xml"
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


func JsonToMap(jsonBytes []byte) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(jsonBytes, &jsonMap)
	return jsonMap, err
}


func MapToJson(jsonMap map[string]interface{}) ([]byte, error) {
	return json.Marshal(jsonMap)
}


func MapToXML(xmlMap map[string]interface{}, rootName string, indent string) ([]byte, error) {
	xmlBytes, err := xml.MarshalIndent(xmlMap, "", indent)
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + "<" + rootName + ">\n" + string(xmlBytes) + "\n</" + rootName + ">"), nil
}


func XmlToMap(xmlBytes []byte) (map[string]interface{}, error) {
	var xmlMap map[string]interface{}
	err := xml.Unmarshal(xmlBytes, &xmlMap)
	return xmlMap, err
}


func XmlToJSON(xmlStr string) ([]byte, error) {
	xmlMap, err := XmlToMap([]byte(xmlStr))
	if err != nil {
		return nil, err
	}

	jsonBytes, err := MapToJson(xmlMap)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}


func JsonToXML(jsonStr string, rootName string, indent string) ([]byte, error) {
	jsonBytes := []byte(jsonStr)

	jsonMap, err := JsonToMap(jsonBytes)
	if err != nil {
		return nil, err
	}

	xmlBytes, err := MapToXML(jsonMap, rootName, indent)
	if err != nil {
		return nil, err
	}

	return xmlBytes, nil
}


func JsonBytesToXML(jsonBytes []byte, rootName string, indent string) ([]byte, error) {
	jsonMap, err := JsonToMap(jsonBytes)
	if err != nil {
		return nil, err
	}

	xmlBytes, err := MapToXML(jsonMap, rootName, indent)
	if err != nil {
		return nil, err
	}

	return xmlBytes, nil
}

