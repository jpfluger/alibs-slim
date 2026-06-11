# alibs-slim

🛑 **ARCHIVED PROJECT**  
This repository represents an older version of this project and is no longer actively maintained here. It has been left public for educational purposes.

🟢 **Note to prospective team members:**  
This repo is a time capsule — old, strange, and unapologetic. Want a seat at my table? Stow your AI toolbox. The whiteboard awaits.  Dissect the code. Ruthlessly roast my old design decisions. Then let’s talk real enterprise architecture, cybersecurity, and networking — the kind that survives in the trenches, not just white-room theory.

---

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
