package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	dirPerm = 0755
)

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func readFileJSON(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(v)
}

// createDirIfNotExist creates the directory to path if it doesn't exist.
func createDirIfNotExist(path string) error {
	return os.MkdirAll(path, dirPerm)
}

// createFile creates the file if it doesn't exist.
func createFile(path string, output interface{}) (int, error) {
	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	switch data := output.(type) {
	case []byte:
		return file.Write(data)
	case string:
		return file.WriteString(data)
	default:
		return 0, fmt.Errorf("Unsupported format type=%T", data)
	}
}

// createFileIfNotExist creates the file if it doesn't exist.
func createFileIfNotExist(path string, output interface{}) (int, error) {
	if IsFileExist(path) {
		return 0, nil
	}
	return createFile(path, output)
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
