package rekor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

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

		var subjects []string
		for _, s := range provenance.Subject {
			verification_string := fmt.Sprintf("%s:%s@sha256:%s", s.Name, vendorFile.Version, s.Digest["sha256"])
			subjects = append(subjects, verification_string)
		}

		if err := validateYaml(*fileContents, subjects); err != nil {
			imageError, ok := err.(*ImageValidationError)
			if ok {
				color.Red(fmt.Sprintf("Cannot locate ` %s ` in %s", imageError.image, downloadedFile))
			}
			return err
		}
	}

	return nil
}

type ImageValidationError struct {
	image string
}

func (e *ImageValidationError) Error() string {
	return fmt.Sprintf("Cannot validate %v", e.image)
}

func validateYaml(y string, subjects []string) error {
	for _, resourceYaml := range strings.Split(y, "---") {
		m := make(map[string]interface{})
		e := yaml.Unmarshal([]byte(resourceYaml), m)
		if e != nil {
			panic(e)
		}

		if err := walkMap(m, subjects); err != nil {
			return err
		}
	}
	return nil
}

func walkMap(m map[string]interface{}, subjects []string) error {
	for k, v := range m {
		switch x := v.(type) {
		case map[string]interface{}:
			walkMap(x, subjects)
		case []interface{}:
			for _, l := range v.([]interface{}) {
				if x, ok := l.(map[string]interface{}); ok {
					walkMap(x, subjects)
				}
			}
		case string:
			subject := v.(string)
			if k == "image" {
				if contains(subjects, subject) {
					color.Yellow(fmt.Sprintf("image: %v\n", subject))
				} else {
					return &ImageValidationError{image: subject}
				}
			}
		}
	}
	return nil
}

func contains(l []string, v string) bool {
	for _, s := range l {
		if v == s {
			return true
		}
	}
	return false
}
