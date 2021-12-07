package rekor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/in-toto/in-toto-golang/in_toto"
	rekorClient "github.com/sigstore/rekor/pkg/client"
	"github.com/sigstore/rekor/pkg/generated/client/entries"
	"github.com/tidwall/pretty"
	"github.com/trmiller/vendorme/cmd/cli/config"
	"github.com/trmiller/vendorme/util"
)

func Validate(vendorFile config.VendorFile, downloadedFile string) (err error) {
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
	}

	return nil
}
