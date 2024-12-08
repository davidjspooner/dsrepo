package repository

import "github.com/davidjspooner/dsrepo/internal/access"

type UserAlias string

type Config struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`
	Namespaces []string `yaml:"namespaces"`
	Local      struct {
		Path string `yaml:"path"`
		API  string `yaml:"api"`
	} `yaml:"local"`
	Upstream struct {
		Url        string    `yaml:"url"`
		Credential UserAlias `yaml:"credential"`
	} `yaml:"upstream"`
	Policies access.PolicyList             `yaml:"policies"`
	Roles    access.RoleList               `yaml:"roles"`
	Users    map[UserAlias]access.RoleName `yaml:"users"`
}

type Credential struct {
	Alias       UserAlias `yaml:"alias"`
	Type        string    `yaml:"type"`
	Key, Secret string    `yaml:"key"`
}
