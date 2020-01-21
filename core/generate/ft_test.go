package generate

import (
	"testing"
)

func TestFTSpec(t *testing.T) {
	spec := NewFTSpec("testuser", "fttest", "ft")

	if err := spec.Build(); err != nil {
		t.Error(err)
	}

	if content, err := spec.GetConfig(); err != nil {
		t.Error(err)
	} else if content == "testconfig" {
		t.Error("mismatch")
	}

	if err := spec.SetConfig("testconfig"); err != nil {
		t.Error(err)
	}

	if content, err := spec.GetConfig(); err != nil {
		t.Error(err)
	} else if content != "testconfig" {
		t.Error("mismatch")
	}

	// if err := spec.Remove(); err != nil {
	// 	t.Error(err)
	// }

}

func TestFTDeploySpec(t *testing.T) {
	spec := NewFTDeploySpec("testuser", "ftdeploy", "ft", "ft")

	if err := spec.Build(); err != nil {
		t.Error(err)
	}

	if content, err := spec.GetConfig(); err != nil {
		t.Error(err)
	} else if content == "testconfig" {
		t.Error("mismatch")
	}

	if err := spec.SetConfig("testconfig"); err != nil {
		t.Error(err)
	}

	if content, err := spec.GetConfig(); err != nil {
		t.Error(err)
	} else if content != "testconfig" {
		t.Error("mismatch")
	}

	// if err := spec.Remove(); err != nil {
	// 	t.Error(err)
	// }

}
