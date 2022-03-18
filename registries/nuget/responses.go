package nuget

type catalogRegistration struct {
	Pages []catalogPage `json:"items"`
}

type catalogPage struct {
	Packages []catalogPackage `json:"items"`
}

type catalogPackageLink struct {
	Entry string `json:"catalogEntry"`
}

type catalogPackage struct {
	Entry catalogPackageEntry `json:"catalogEntry"`
}

type catalogPackageEntry struct {
	Id               string                          `json:"id"`
	Version          string                          `json:"version"`
	DependencyGroups []catalogPackageDependencyGroup `json:"dependencyGroups"`
}

type catalogPackageDependencyGroup struct {
	Dependencies []catalogPackageDependency `json:"dependencies"`
}

type catalogPackageDependency struct {
	Id              string `json:"id"`
	DependencyRange string `json:"range"`
}
