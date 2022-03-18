package npm

type packageInfo struct {
	Id       string                        `json:"_id"`
	Versions map[string]packageVersionInfo `json:"versions"`
}

type packageVersionInfo struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}
