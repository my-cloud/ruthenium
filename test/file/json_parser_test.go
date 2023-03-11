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
	//path := filepath.Join(os.TempDir(), "/unreadable.json")
	//f, _ := os.Create(path)
	//_ = os.WriteFile(path, d1, 0644)
	////
	////parser := file.NewJsonParser()
	////var parsed interface{}
	////
	////// Act
	////err = parser.Parse("/tmp/dat1", &parsed)
	////
	////// Assert
	////test.Assert(t, err == nil, "Error is not nil whereas it should be.")
	//_ = f.Close()
	//_ = os.Remove(path)
}

func Test_Parse_UnableToCloseFile_ReturnsError(t *testing.T) {
}

func Test_Parse_UnableToUnmarshalBytes_ReturnsError(t *testing.T) {
}
