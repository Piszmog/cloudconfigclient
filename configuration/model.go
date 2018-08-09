package configuration

type Source struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	State           string           `json:"state"`
	PropertySources []PropertySource `json:"propertySources"`
}

type PropertySource struct {
	Name   string            `json:"name"`
	Source map[string]string `json:"source"`
}
