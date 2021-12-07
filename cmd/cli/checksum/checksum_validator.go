package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/trmiller/vendorme/cmd/cli/config"
)

func Validate(vendorFile config.VendorFile, downloadedFile string) (err error) {
	f, err := os.Open(downloadedFile)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	sha256_hash := hex.EncodeToString(h.Sum(nil))

	if sha256_hash != vendorFile.Sha256 {
		color.Red("sha256 does not match for " + downloadedFile)
		return errors.New("Sha256 does not match")
	}

	return nil
}
