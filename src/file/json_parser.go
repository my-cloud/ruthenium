package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type JsonParser struct {
	path string
}

func NewJsonParser(path string) *JsonParser {
	return &JsonParser{path}
}

func (parser *JsonParser) Parse(any interface{}) error {
	jsonFile, err := os.Open(parser.path)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}
	if err = jsonFile.Close(); err != nil {
		return fmt.Errorf("unable to close file: %w", err)
	}
	if err = json.Unmarshal(byteValue, &any); err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}
	return nil
}
