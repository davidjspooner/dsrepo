package tfprovider

type Platform struct {
	OS   string
	Arch string
}

type Version struct {
	Version   string
	Protocols []string
	Platforms []Platform
}

type Provider struct {
	Namespace string
	Provider  string
	Versions  []Version
}

type key struct {
	Namespace string
	Provider  string
	Version   string
	OS        string
	Arch      string
}
