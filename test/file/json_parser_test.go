package file

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/file"
	"github.com/my-cloud/ruthenium/test"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_Parse_UnableToOpenFile_ReturnsError(t *testing.T) {
	// Arrange
	parser := file.NewJsonParser()
	var parsed interface{}

	// Act
	err := parser.Parse("", &parsed)

	// Assert
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
	if err != nil {
		expectedErrorMessage := "unable to open file"
		actualErrorMessage := err.Error()
		test.Assert(t, strings.Contains(actualErrorMessage, expectedErrorMessage), fmt.Sprintf("Wrong error message.\nExpected: %s\nActual:   %s", expectedErrorMessage, actualErrorMessage))
	}
}

func Test_Parse_UnableToUnmarshalBytes_ReturnsError(t *testing.T) {
	// Arrange
	jsonFile, _ := ioutil.TempFile("", "Test_Parse_UnableToUnmarshalBytes_ReturnsError.json")
	defer func() { _ = os.Remove(jsonFile.Name()) }()
	jsonData := []byte(`{`)
	_, _ = jsonFile.Write(jsonData)
	_ = jsonFile.Close()
	parser := file.NewJsonParser()
	var person interface{}

	// Act
	err := parser.Parse(jsonFile.Name(), &person)

	// Assert
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
	if err != nil {
		expectedErrorMessage := "unable to unmarshal"
		actualErrorMessage := err.Error()
		test.Assert(t, strings.Contains(actualErrorMessage, expectedErrorMessage), fmt.Sprintf("Wrong error message.\nExpected: %s\nActual:   %s", expectedErrorMessage, actualErrorMessage))
	}
}

func Test_Parse_ValidFile_OutputFilled(t *testing.T) {
	// Arrange
	jsonFile, _ := ioutil.TempFile("", "Test_Parse_ValidFile_OutputFilled.json")
	defer func() { _ = os.Remove(jsonFile.Name()) }()
	expectedPersonName := "John"
	expectedPersonAge := 30
	jsonDataString := fmt.Sprintf(`{"name":"%s","age":%d}`, expectedPersonName, expectedPersonAge)
	jsonData := []byte(jsonDataString)
	_, _ = jsonFile.Write(jsonData)
	_ = jsonFile.Close()
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	parser := file.NewJsonParser()
	var person Person

	// Act
	_ = parser.Parse(jsonFile.Name(), &person)

	// Assert
	test.Assert(t, person.Name == expectedPersonName || person.Age == expectedPersonAge, "Wrong field value.")
}
