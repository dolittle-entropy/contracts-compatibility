package main

import (
	"dolittle.io/contracts-compatibility/artifacts"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	//token, err := docker.GetAuthTokenFor("dolittle/runtime")
	//if err != nil {
	//	fmt.Println("Token error", err)
	//	return
	//}
	//
	//cache, err := os.Create("graph.json")
	//if err != nil {
	//	fmt.Println("Failed to open graph.json file")
	//	return
	//}
	//
	//graph := artifacts.CreateGraphFor(
	//	artifacts.NewReleaseListResolver(
	//		docker.NewReleaseListerFor(token, "dolittle/runtime"),
	//		docker.NewDependencyResolverFor(token, "dolittle/runtime", dotnet.NewDepsResolverFor("Dolittle.Runtime.Contracts"), "app/Dolittle.Runtime.Server.deps.json", "app/Server.deps.json"),
	//	),
	//	map[string]*artifacts.ReleaseListResolver{
	//		"DotNET": artifacts.NewReleaseListResolver(
	//			nuget.NewReleaseListerFor("Dolittle.SDK.Services"),
	//			nuget.NewDependencyResolverFor("Dolittle.SDK.Services", "Dolittle.Contracts"),
	//		),
	//		"JavaScript": artifacts.NewReleaseListResolver(
	//			npm.NewReleaseListerFor("@dolittle/sdk.services"),
	//			npm.NewDependencyResolverFor("@dolittle/sdk.services", "@dolittle/contracts"),
	//		),
	//	},
	//)
	//
	//encoder := json.NewEncoder(cache)
	//encoder.SetIndent("", "  ")
	//encoder.Encode(graph)
	//cache.Close()

	cache, err := os.Open("graph.json")
	if err != nil {
		fmt.Println("Failed to open graph.json file", err)
		return
	}

	graph := &artifacts.Graph{}
	err = json.NewDecoder(cache).Decode(graph)
	if err != nil {
		fmt.Println("Failed to parse graph.json file", err)
		return
	}

	compatibility := ResolveCompatibilityFrom(graph)
	//encoder := json.NewEncoder(os.Stdout)
	//encoder.SetIndent("", "  ")
	//encoder.Encode(compatibility)
	WriteTables(os.Stdout, compatibility)
}
