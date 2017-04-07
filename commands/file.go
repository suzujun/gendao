package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

const (
	dirPerm  = 0755
	filePerm = 0644
	buffSize = 1024
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
	buff, err := readFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(buff, v)
}

// createDirIfNotExist creates the directory to path if it doesn't exist.
func createDirIfNotExist(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	os.MkdirAll(path, dirPerm)
	return nil
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
		v := reflect.ValueOf(output)
		return 0, fmt.Errorf("Unsupported format type=%s", v.Type())
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
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
