package models

type Secret struct {
	Name        string            `json:"name" yaml:"name"`
	VersionName string            `json:"version-name" yaml:"version-name"`
	ShortName   string            `json:"short-name" yaml:"short-name"`
	Labels      map[string]string `json:"labels" yaml:"labels"`
	SecretValue string            `json:"secret-value" yaml:"secret-value"`
}
