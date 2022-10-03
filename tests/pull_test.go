package cli_test

import (
	"os"
	"path"
	"testing"

	"github.com/buildsec/vendorme/cmd/cli"
)

func TestExecErrorsWhenFileNotPresent(t *testing.T) {
	p := cli.PullCommand{
		VendorMeConfig: "defnitelyNotThere.yaml",
	}

	var args []string

	err := p.Exec(nil, args)

	if err == nil {
		t.Error("Should return error when file not found")
	}
}

func TestExecErrorsWithMalformedYaml(t *testing.T) {
	p := cli.PullCommand{
		VendorMeConfig: "bad_yaml.yaml",
	}

	var args []string

	err := p.Exec(nil, args)

	if err == nil {
		t.Error("Expect unmarshalling error here")
	}
}

func TestExecTektonYaml(t *testing.T) {
	tmp := t.TempDir()
	dir, derr := os.Getwd()
	if derr != nil {
		t.Error(derr)
	}
	os.Chdir(tmp)
	defer os.Chdir(dir)

	p := cli.PullCommand{
		VendorMeConfig: path.Join(dir, "tekton.yaml"),
	}

	var args []string

	err := p.Exec(nil, args)

	if err != nil {
		t.Error(err)
	}
}
