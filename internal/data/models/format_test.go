package models

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDirectFormat(t *testing.T) {
	f, err := os.Open("./assets/format.direct.json")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	var format Format
	if err = json.NewDecoder(f).Decode(&format); err != nil {
		t.Error(err)
	}

	if format.Codename != "direct" || !format.IsDirect() {
		t.Error("Format must be Direct type")
	}
}
