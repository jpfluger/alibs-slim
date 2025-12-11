# Fork Notice

This directory contains a forked and modified version of the [badgerstore](https://github.com/alexedwards/scs/tree/master/badgerstore) package from the [SCS (Session Control System)](https://github.com/alexedwards/scs) project, originally authored by Alex Edwards under the MIT License.

## Why Fork and Copy?

The original `badgerstore` implementation relies on an outdated version of [Badger](https://github.com/dgraph-io/badger) (v1), which is no longer compatible with our project's requirements. We needed to update it to Badger v4 to leverage modern features, bug fixes, and performance improvements while ensuring compatibility with our Go 1.25 toolchain.

Initially, we considered maintaining a traditional fork and using Go module replacements in our `go.mod` file (e.g., via `replace` directives) to point to the forked repository. However, this approach led to cumbersome dependency management, especially in a larger monorepo setup like ours. To simplify integration and avoid "ugly" overrides in downstream projects, we opted to copy the modified source code directly into a subdirectory of our Apache v2-licensed project.

## Why Not Contribute Back Upstream?

We chose not to submit these changes as a pull request to the original SCS repository for a few key reasons:

- **Potential Breaking Changes**: Updating from Badger v1 to v4 involves API changes and behavioral differences that could break existing users of the original package who rely on older Badger versions. Badger v4 includes breaking API updates, requiring code refactors that could introduce incompatibilities without a major version bump in SCS.

- **Toolchain Compatibility**: Badger v4 requires Go 1.19+, which would force a minimum Go version increase in the upstream `go.mod`, breaking builds for users on older toolchains (e.g., Go 1.12–1.18). Our project compiles with Go 1.25, which may introduce subtle differences or require adjustments not yet tested in the upstream repository.

- **Project-Specific Customizations**: The modifications are tailored to our use case (e.g., integration with our custom adapter system in `aconns/aclient-badger`). Upstream contributions would need broader testing and alignment with SCS's goals, which might delay or complicate our development.

If the upstream project expresses interest in a v4 update, we're open to collaborating or providing our changes as a reference.

## License and Attribution

The code in this directory is derived from the original MIT-licensed work by Alex Edwards. As per the MIT License, we've preserved the original copyright notice and license terms in the `LICENSE` file. Our broader project is licensed under Apache v2, which is compatible with MIT, allowing this incorporation. No relicensing of the forked code has occurred—it remains under MIT.

If you're using this code, please review the original SCS documentation for best practices, as our fork maintains core functionality but includes project-specific refactorings.

For any questions or concerns, feel free to open an issue in our repository.