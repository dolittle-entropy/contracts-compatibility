package docker

type dockerAuthToken struct {
	Token string `json:"token"`
}

type dockerTagList struct {
	Tags []string `json:"tags"`
}

type dockerManifest struct {
	Layers []dockerLayer `json:"layers"`
}

type dockerLayer struct {
	Digest string `json:"digest"`
	Size   uint   `json:"size"`
}
