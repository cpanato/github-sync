package config

type Config struct {
	Users        []User       `yaml:"users"`
	Repositories []Repository `yaml:"repositories"`
}

type User struct {
	Username string `yaml:"username"`
	Email    string `yaml:"email"`
	Role     string `yaml:"role"`
}

type Collaborator struct {
	Username   string `yaml:"username"`
	Email      string `yaml:"email"`
	Permission string `yaml:"permission"`
}
type Repository struct {
	Name          string         `yaml:"name"`
	Collaborators []Collaborator `yaml:"collaborators"`
}
