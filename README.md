# alibs-slim

Common shared libraries that serve as the plumbing for a myriad of [`golang`](https://go.dev/) binaries, with web-related constructs built around the [`Echo Framework`](https://echo.labstack.com/).

## aconns

See the connections [README](aconns/README.md).

## Versioning Strategy

Versioning is handled manually by maintainers after merge. Once changes are verified and merged into `main`, a version bump is applied to the core module and all subdirectory Go modules for consistency. (We aren't using git submodules.) Only actual changes to code receive a version bump. READMEs or testing modules do not receive a new version.

This ensures that consumers of `g-aconns` (which references multiple drivers) can depend on synchronized, tagged versions across modules.

By keeping each driver implementation in its own Go module, we avoid unnecessary dependency bloat — improving security scan accuracy and ensuring clean dependency trees.

> Note: At this point in time, this strategy seems reasonable. We're open to change, if another method becomes more suitable.

## License

`alibs-slim` is licensed under the [Apache License 2.0](LICENSE). For more details, please refer to the [LICENSE](LICENSE) file.
