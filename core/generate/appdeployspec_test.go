package generate

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestApplicationDeploySpec(t *testing.T) {
	spec := &ApplicationDeploySpec{
		Account:         "admintest",
		Name:            "ft",
		CoinfigFileName: "config.yaml",
	}

	if err := spec.Build(); err != nil {
		t.Error(err)
	}

	copyto := func(fname, tname string) error {
		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(tname, bts, os.ModePerm)
	}

	if err := spec.SetConfig("test", copyto); err != nil {
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
