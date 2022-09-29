//
// Copyright 2021 Tim Miller.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/buildsec/vendorme/cmd/cli/checksum"
	"github.com/buildsec/vendorme/cmd/cli/config"
	"github.com/buildsec/vendorme/cmd/cli/rekor"
	"github.com/buildsec/vendorme/util"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type PullCommand struct {
	VendorMeConfig string
}

func (c *PullCommand) Exec(ctx context.Context, directory []string) (err error) {
	if _, err := os.Stat(c.VendorMeConfig); err != nil {
		// File doesn't exist
		return err
	}

	vendorme_config, err := readConf(c.VendorMeConfig)

	if err != nil {
		return err
	}

	for _, vendorFile := range vendorme_config.Files {
		downloadedFile, err := util.DownloadFile(vendorFile.ReleaseFile, vendorFile.DestinationDir)
		if err != nil {
			return err
		}

		switch vendorFile.ValidationType {
		case config.Rekor:
			err = rekor.Validate(vendorFile, *downloadedFile)
		case config.Sha256:
			err = checksum.Validate(vendorFile, *downloadedFile)
		default:
			err = errors.New("unknown validation type " + vendorFile.ValidationType)
		}

		if err != nil {
			return err
		}

		color.Green("Sucessfully vendored & validated " + vendorFile.ReleaseFile)
	}

	return nil
}

func readConf(filename string) (*config.VendorConfig, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &config.VendorConfig{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
