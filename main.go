package main

import (
	"dolittle.io/contracts-compatibility/artifacts"
	"dolittle.io/contracts-compatibility/dependencies/dotnet"
	"dolittle.io/contracts-compatibility/registries/docker"
	"dolittle.io/contracts-compatibility/registries/npm"
	"dolittle.io/contracts-compatibility/registries/nuget"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"os"
)

func main() {
	help := flag.Bool("h", false, "Print help information")
	output := flag.String("o", "markdown", "Output format [markdown,json]")
	username := flag.String("docker-username", "", "DockerHub username to use to authenticate with")
	password := flag.String("docker-password", "", "DockerHub password/PAT to use to authenticate with")
	flag.Parse()

	if *help {
		fmt.Println("Contracts Compatibility a tool to resolve compatible versions of the Dolittle Runtime and SDKs")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if !(*output == "markdown" || *output == "json") {
		fmt.Fprintln(os.Stderr, "Invalid output format ", *output, "specified. Only 'markdown' or 'json' is supported.")
		os.Exit(1)
	}

	graph, err := createGraph(*username, *password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed while creating dependency graph", err)
		os.Exit(2)
	}

	compatibility := ResolveCompatibilityFrom(graph)

	if *output == "markdown" {
		WriteTables(os.Stdout, compatibility)
	}
	if *output == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(compatibility)
	}
}

func createGraph(username, password string) (*artifacts.Graph, error) {
	var token docker.AuthToken
	var err error

	if username != "" && password != "" {
		token, err = docker.GetAuthenticatedUserAuthTokenFor("dolittle/runtime", username, password)
	} else {
		token, err = docker.GetAuthTokenFor("dolittle/runtime")
	}

	if err != nil {
		return nil, fmt.Errorf("could not get Docker Hub authentication token, %w", err)
	}

	resolver := docker.NewDependencyResolverFor(token, "dolittle/runtime", dotnet.NewDepsResolverFor("Dolittle.Runtime.Contracts"), "app/Dolittle.Runtime.Server.deps.json", "app/Server.deps.json")
	resolver.ResolveDependencyForVersion(semver.New("8.4.1"))
	return nil, errors.New("HELLO")

	graph := artifacts.CreateGraphFor(
		artifacts.NewReleaseListResolver(
			docker.NewReleaseListerFor(token, "dolittle/runtime"),
			docker.NewDependencyResolverFor(token, "dolittle/runtime", dotnet.NewDepsResolverFor("Dolittle.Runtime.Contracts"), "app/Dolittle.Runtime.Server.deps.json", "app/Server.deps.json"),
		),
		map[string]*artifacts.ReleaseListResolver{
			"DotNET": artifacts.NewReleaseListResolver(
				nuget.NewReleaseListerFor("Dolittle.SDK.Services"),
				nuget.NewDependencyResolverFor("Dolittle.SDK.Services", "Dolittle.Contracts"),
			),
			"JavaScript": artifacts.NewReleaseListResolver(
				npm.NewReleaseListerFor("@dolittle/sdk.services"),
				npm.NewDependencyResolverFor("@dolittle/sdk.services", "@dolittle/contracts"),
			),
		},
	)
	return graph, nil
}
