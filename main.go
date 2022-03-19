package main

import (
	"dolittle.io/contracts-compatibility/artifacts"
	"dolittle.io/contracts-compatibility/dependencies/dotnet"
	"dolittle.io/contracts-compatibility/registries/docker"
	"dolittle.io/contracts-compatibility/registries/npm"
	"dolittle.io/contracts-compatibility/registries/nuget"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	help := flag.Bool("h", false, "Print help information")
	output := flag.String("o", "markdown", "Output format [markdown,json]")
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

	graph, err := createGraph()
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

func createGraph() (*artifacts.Graph, error) {
	token, err := docker.GetAuthTokenFor("dolittle/runtime")
	if err != nil {
		return nil, fmt.Errorf("could not get Docker Hub authentication token, %w", err)
	}

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
