package sfj_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/pchchv/sfj"
)

var (
	expectedGeneration []byte
	subStruct          = true
	insecure           = false
	pkg                = "example"
	headerMap          = map[string]string{}
	server             = "https://api.github.com"
	lines              = []string{
		"/",
		"/repos/:user/:repo pchchv sfj",
	}
)

func init() {
	var err error
	expectedGeneration, err = os.ReadFile("expected_out.txt")
	if err != nil {
		panic(err)
	}
}

func TestDo(t *testing.T) {
	file, err := sfj.Do(pkg, server, lines, headerMap, insecure, subStruct)
	if err != nil {
		t.Fatalf("No error expected but got: %e", err)
	}

	if !bytes.Equal(file, expectedGeneration) {
		t.Errorf("Results should be equal, but they differs")
		t.Error(string(file))
		os.WriteFile("out.txt", file, 0644)
	}
}
