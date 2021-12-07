package config

const (
	Rekor  string = "rekor"
	Sha256 string = "sha256"
)

type VendorConfig struct {
	Files []VendorFile `yaml:"files"`
}

type VendorFile struct {
	ReleaseFile    string `yaml:"release_file"`
	RekorUUID      string `yaml:"rekor_uuid"`
	DestinationDir string `yaml:"destination_dir"`
	Version        string `yaml:"version"`
	ValidationType string `yaml:"validation_type"`
	Sha256         string `yaml:"sha256"`
}
