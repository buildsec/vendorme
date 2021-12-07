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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/in-toto/in-toto-golang/in_toto"
	rekorClient "github.com/sigstore/rekor/pkg/client"
	"github.com/sigstore/rekor/pkg/generated/client/entries"
	"github.com/tidwall/pretty"
	"github.com/trmiller/vendorme/util"
	"gopkg.in/yaml.v3"
)

type PullCommand struct {
	VendorMeConfig string
}

type VendorConfig struct {
	Files []VendorFile `yaml:"files"`
}

type VendorFile struct {
	ReleaseFile    string `yaml:"release_file"`
	RekorUUID      string `yaml:"rekor_uuid"`
	DestinationDir string `yaml:"destination_dir"`
	Version        string `yaml:"version"`
}

func (c *PullCommand) Exec(ctx context.Context, directory []string) (err error) {
	if _, err := os.Stat(c.VendorMeConfig); err != nil {
		// File doesn't exist
		return err
	}

	config, err := readConf(c.VendorMeConfig)

	if err != nil {
		return err
	}

	for _, vendorFile := range config.Files {
		downloadedFile, err := util.DownloadFile(vendorFile.ReleaseFile, vendorFile.DestinationDir)
		if err != nil {
			return err
		}

		err = validate(ctx, vendorFile, *downloadedFile)

		if err != nil {
			return err
		}
	}

	return nil
}

func validate(ctx context.Context, vendorFile VendorFile, downloadedFile string) (err error) {
	//TODO: don't do this every time, also config this to allow pointing to your own rekor
	client, err := rekorClient.GetRekorClient("https://rekor.sigstore.dev/")
	if err != nil {
		return err
	}

	req := entries.NewGetLogEntryByUUIDParams()
	req.SetEntryUUID(vendorFile.RekorUUID)

	resp, err := client.Entries.GetLogEntryByUUID(req)

	if err != nil {
		return err
	}

	for k, entry := range resp.Payload {
		if k != vendorFile.RekorUUID {
			// This shouldn't really happen, the uuid id you asked for isn't what you got back
			continue
		}

		decoded, err := base64.StdEncoding.DecodeString(string(entry.Attestation.Data))

		if err != nil {
			return err
		}

		real_json := fmt.Sprintf("%s", pretty.Pretty(decoded))

		var provenance in_toto.ProvenanceStatement

		if err := json.Unmarshal([]byte(real_json), &provenance); err != nil {
			fmt.Println("Error unmarshalling predicate")
			return err
		}

		fileContents, err := util.ReadFile(downloadedFile)

		for _, s := range provenance.Subject {
			verification_string := fmt.Sprintf("%s:%s@sha256:%s", s.Name, vendorFile.Version, s.Digest["sha256"])
			if !strings.Contains(*fileContents, verification_string) {
				color.Red(fmt.Sprintf("Cannot locate ` %s ` in %s", verification_string, downloadedFile))
				return errors.New("Cannot validate " + verification_string)
			}
		}

		color.Green("Sucessfully vendored & validated " + vendorFile.ReleaseFile)
	}

	return nil
}

// I don't like this here.  TODO: Move elsewhere.
func readConf(filename string) (*VendorConfig, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &VendorConfig{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
