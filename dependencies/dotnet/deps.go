package dotnet

type dotnetDependencies struct {
	Targets map[string]map[string]dotnetTarget `json:"targets"`
}

type dotnetTarget struct {
	Dependencies map[string]string `json:"dependencies"`
}
