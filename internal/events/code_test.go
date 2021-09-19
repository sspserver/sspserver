package events

import (
	"bytes"
	"testing"
)

func Test_CodeCompression(t *testing.T) {
	dataTest := []byte("Test")
	code1 := CodeObj(dataTest, nil)
	code2 := code1.Compress()
	if !bytes.Equal(code2.Decompress().Data(), dataTest) {
		t.Errorf("Invalid compression: %s", code2.Error())
	}
}
