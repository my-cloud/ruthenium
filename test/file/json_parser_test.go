package file

import (
	"github.com/my-cloud/ruthenium/src/file"
	"github.com/my-cloud/ruthenium/test"
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
}

func Test_Parse_UnableToReadFile_ReturnsError(t *testing.T) {
	//// Arrange
	//d1 := []byte("hello\ngo\n")
	//err := os.WriteFile("/tmp/dat1", d1, 0644)
	//check(err)
	//
	//f, err := os.Create("/tmp/dat2")
	//check(err)
	//
	//defer f.Close()
	//
	//parser := file.NewJsonParser("")
	//var parsed interface{}
	//
	//// Act
	//err = parser.Parse(&parsed)
	//
	//// Assert
	//test.Assert(t, err != nil, "Error is nil whereas it should not.")
}
