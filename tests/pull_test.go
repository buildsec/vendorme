package cli_test

import (
	"github.com/trmiller/vendorme/cmd/cli"
	"testing"
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
