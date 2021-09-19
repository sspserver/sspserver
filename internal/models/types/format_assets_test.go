package types

import (
	"testing"
)

func Test_FormatsLoad(t *testing.T) {
	if len(MockFormats()) != 7 {
		t.Error("Invalid formats loading")
	}
}
