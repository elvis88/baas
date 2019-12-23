package generate

import (
	"testing"
)

func TestFTDeploySpec(t *testing.T) {
	spec := &FTDeploySpec{
		Account: "admin",
		Name:    "ft",
	}

	if err := spec.Build(); err != nil {
		t.Error(err)
	}

	if err := spec.SetConfig("test"); err != nil {
		t.Error(err)
	}

	if content, err := spec.GetConfig(); err != nil {
		t.Error(err)
	} else if content != "test" {
		t.Error("mismatch")
	}

	// if err := spec.Remove(); err != nil {
	// 	t.Error(err)
	// }
}
