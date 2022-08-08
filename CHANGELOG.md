# [1.1.0] - 2022-8-8 [PR: #4](https://github.com/dolittle/contracts-compatibility/pull/4)
## Summary

Makes the tool work with Runtimes after version `8.0.0`. There were two changes outside the tool that made it stop working for newer versions.

### Added

- The option of passing in a username and password/PAT for Docker Hub. This is useful to avoid rate-limiting when debugging the tool locally.

### Fixed

- For the Runtime image, the tool now looks for a dependency to `Dolittle.Contracts` instead of `Dolittle.Runtime.Contracts`. We merged these to packages into the former recently, so now we need to rely on this one. This still works for older versions since we published them together previously.
- For the Runtime image, the tool no longer accepts "fat manifests" from the Docker Hub API. Since we started publishing "multi-platform images", these "fat manifest" files started returning a different format we didn't support. With this change, Docker Hub returns the `Linux amd64` manifest by default, thereby fixing the problem for now.


# [1.0.1] - 2022-3-21 [PR: #2](https://github.com/dolittle/contracts-compatibility/pull/2)
## Summary

Generates Markdown tables with newest (coolest) versions on the top.

### Changed

- The generated Markdown tables now show the newest versions on the top.


# [1.0.0] - 2022-3-20 [PR: #1](https://github.com/dolittle/contracts-compatibility/pull/1)
## Summary

This is the first release of the Contracts Compatibility tool. A tool which resolves compatible versions of the Dolittle Runtime and SDKs based on the current public releases of the artefacts - by fetching the releases from their respective registries. See the README file for more details.


