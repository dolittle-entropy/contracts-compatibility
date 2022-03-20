# Dolittle Contracts Compatibility tool

The `contracts-compatibility` checks what versions of the Dolittle [Runtime](https://github.com/dolittle/Runtime) and
SDKs that are compatible with each other. And can generate a Markdown table or a JSON structure. Currently, the tool
checks the [DotNet SDK](https://github.com/dolittle/DotNET.SDK) and the
[JavaScript SDK](https://github.com/dolittle/JavaScript.SDK). It was built to generate the compatibility table on
[dolittle.io](https://dolittle.io/docs/reference/runtime/compatibility/), and for other internal uses.

Usage:
```shell
contracts-compatibility -h
Contracts Compatibility a tool to resolve compatible versions of the Dolittle Runtime and SDKs
Usage:
  -h    Print help information
  -o string
        Output format [markdown,json] (default "markdown")
```

Running the tool fetches all currently released versions from the public repositories, and resolves the compatible
versions.

> **Warning**: for the Runtime, the tool fetches the Docker manifests from Docker Hub without authentication. Since
> there are quite a few released versions, this quickly eats into your IPs rate limits. So be careful about running it
> a lot of times. Read more on the
> [Docker Hub rate limit documentation](https://docs.docker.com/docker-hub/download-rate-limit/#definition-of-limits).