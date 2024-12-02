package repository

type Config struct {
	Name, Type string
	Namespaces []string
	Upstream   struct {
		Url        string
		Credential string
	}
}

type Credential struct {
	Alias              string
	Username, Password string
}
