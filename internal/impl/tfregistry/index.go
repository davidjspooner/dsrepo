package tfregistry

type Platform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type Version struct {
	Version   string     `json:"version"`
	Protocols []string   `json:"protocols"`
	Platforms []Platform `json:"platforms"`
}

type Index struct {
	Versions []*Version `json:"versions"`
}
