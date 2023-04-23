package convertors

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"bytes"
	"encoding/binary"

	"gopkg.in/yaml.v3"
)


/*
*
Struct for reading and working with CSV Imports
*/
type csvData struct {
	header []string
	rows   [][]string
}


func CsvRead(filename string) (*csvData, error) {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the CSV data
	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Return the CSV data
	return &csvData{header, rows}, nil
}


func CsvValidate(csvData *csvData) error {
	// Check if the CSV has a header row
	if len(csvData.header) == 0 {
		return errors.New("CSV file is missing header row")
	}

	// Check if all rows have the same number of columns as the header row
	for i, row := range csvData.rows {
		if len(row) != len(csvData.header) {
			return fmt.Errorf("CSV file row %d has incorrect number of columns", i+1)
		}
	}

	// Validation passed
	return nil
}


// // ReadFile reads a file from disk and returns its contents as a byte slice
// func ReadFile(filePath string) ([]byte, error) {
// 	return ioutil.ReadFile(filePath)
// }


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


func txtToBytes(byteSlice []byte) string {
	return string(byteSlice)
}


func bytesToTxt(txt string) []byte {
	return []byte(txt)
}


func objToBytes(byteSlice []byte) string {
	return string(byteSlice)
}


func bytesToObj(txt string) []byte {
	return []byte(txt)
}


func bytesToNum(byteSlice []byte) (uint32, error) {
	// byteSlice := []byte{0x39, 0x30, 0x00, 0x00} // equivalent to uint32(12345) in little-endian byte order
	var num uint32
	err := binary.Read(bytes.NewReader(byteSlice), binary.LittleEndian, &num)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}
	fmt.Println(num) // Output: 12345
	return num, nil
}

func bytesToInt(byteSlice []byte) (int, error) {
    var num uint32
    err := binary.Read(bytes.NewReader(byteSlice), binary.LittleEndian, &num)
    if err != nil {
        return 0, err
    }
    return int(num), nil
}


func numToBytes(num int) []byte {
	// num := uint32(12345)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		fmt.Println("Error:", err)
		return []byte{}
	}
	byteSlice := buf.Bytes()
	fmt.Println(byteSlice) // Output: [57 48 0 0]
	return byteSlice
}


/*
	type Person struct {
		Name    string
		Age     int
		Address string
	}

	func main() {
	    byteSlice := []byte(`{"name":"John","age":30,"address":"New York"}`)

	    var person Person
	    err := bytesToObject(byteSlice, &person)
	    if err != nil {
	        fmt.Println("Error:", err)
	        return
	    }

	    // use the person object
	    fmt.Printf("Name: %s, Age: %d, Address: %s\n", person.Name, person.Age, person.Address)
	}
*/
func bytesToObject(data []byte, v interface{}) (interface{}, error) {
	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}


/*
	type Person struct {
		Name    string
		Age     int
		Address string
	}

	func main() {
		person := Person{Name: "John", Age: 30, Address: "New York"}

		byteSlice, err := objectToBytes(person)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// use the byte slice
		fmt.Println("Byte slice:", byteSlice)
	}
*/
func objectToBytes(v interface{}) ([]byte, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}


// func main() {
// 	// Path to the YAML file
// 	yamlFilePath := "./example.yaml"

// 	// Convert the YAML file to JSON
// 	jsonData, err := ConvertYamlToJson(yamlFilePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Print the JSON data
// 	fmt.Println(string(jsonData))
// }

