package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type JsonParser struct{}

func NewJsonParser() *JsonParser {
	return &JsonParser{}
}

func (parser *JsonParser) Parse(path string, output interface{}) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}
	if err = jsonFile.Close(); err != nil {
		return fmt.Errorf("unable to close file: %w", err)
	}
	if err = json.Unmarshal(byteValue, &output); err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}
	return nil
}
