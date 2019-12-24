package generate

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestApplicationSpec(t *testing.T) {
	spec := &ApplicationSpec{
		Account:         "admintest",
		Name:            "ft",
		CoinfigFileName: "genesis.json",
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
