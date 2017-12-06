package yamlssm_test

import (
	"fmt"
	"testing"

	"github.com/suzuken/yamlssm"
)

func Example() {
	type T struct {
		Key string
	}
	var t T

	// t.Key will be filled by your ssm value set by `your.ssm.key` automatically.
	yamlssm.Unmarshal([]byte(`key: ssm://your.ssm.key`), &t)

	// YAML file without ssm:// value is working as intended.
	yamlssm.Unmarshal([]byte(`key: value`), &t)

	fmt.Print(t.Key)
	// Output:
	// value
}

func TestNormalYAML(t *testing.T) {
	in := []byte(`key: normal`)

	type T struct {
		Key string
	}

	var tt T
	if err := yamlssm.Unmarshal(in, &tt); err != nil {
		t.Errorf("unmarshal error: %s", err)
	}

	if tt.Key != "normal" {
		t.Errorf("want %s, got %s", "normal", tt.Key)
	}
}

func TestReplace(t *testing.T) {
	in := []byte(`key: ssm://testkey`)

	type T struct {
		Key string
	}

	var tt T
	if err := yamlssm.Unmarshal(in, &tt); err != nil {
		t.Errorf("unmarshal error: %s", err)
	}

	if tt.Key != "normal" {
		t.Errorf("want %s, got %s", "normal", tt.Key)
	}
}
