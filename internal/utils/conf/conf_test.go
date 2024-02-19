package conf

import (
	"os"
	"testing"
)

func TestLoadConf(t *testing.T) {
	data, err := os.ReadFile("testdata/conf.md")
	if err != nil {
		t.Fatal(err)
	}

	// t.Log(string(data))

	config, err := LoadConf(string(data))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(config)
}
