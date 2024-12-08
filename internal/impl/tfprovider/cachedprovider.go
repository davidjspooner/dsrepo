package tfprovider

import "time"

type Platform struct {
	OS   string
	Arch string
}

type Version struct {
	Version   string
	Protocols []string
	Platforms []Platform
}

// Catalog represents a Terraform provider catalog.
type Provider struct {
	Namespace string
	Name      string
	Since     time.Time
	Versions  []Version
}
