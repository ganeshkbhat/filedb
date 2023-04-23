package gotxt

import (
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


func TxtToBytes(byteSlice []byte) string {
	return string(byteSlice)
}


func BytesToTxt(txt string) []byte {
	return []byte(txt)
}

