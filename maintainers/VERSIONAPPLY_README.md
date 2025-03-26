# versionapply.sh - Best Practice Guide

> **MAINTAINERS ONLY**: These scripts are to assist maintainers on bumping version numbers for all sub-modules. This is performed after contributor merge requests are accepted by the maintainer.

The `versionapply.sh` script automates the process of updating version numbers, committing changes, and tagging the repository and its submodules. With options to specify branches, this guide outlines best practices for using `versionapply.sh` to manage versioning and ensure consistent commits and tags across project repositories.

> Always execute the script inside the `maintainers` directory.

## Prerequisites

- Ensure `git` is installed and configured for your project.
- You should have push permissions for the chosen branch.
- The `versionapply.sh` script and a `version` file containing the current version number should be in the root directory of your project.

## Overview of Best Practices

1. **Prepare the Environment and Repository**
    - Before running the script, ensure all local changes are committed or stashed to prevent conflicts.
    - Choose a specific branch for version management (e.g., a release or staging branch) to avoid interference with ongoing development work.

2. **Use Branch-Specific Updates**
    - Use the `-B` or `--branch` option to specify the branch you wish to update and push changes to. By default, the branch is set to `main`.
    - Example:
      ```bash
      ./versionapply.sh -G "Update dependencies" -T -B release
      ```
    - **Best Practice**: Define a branch strategy (e.g., main, release, staging) that aligns with your team’s workflow. For instance, use a `release` branch for production versions and `main` for general updates.

3. **Version Management with `go.mod` Updates**
    - By default, the script updates `go.mod` files in the specified submodules.
    - Use the `-N` or `--noupdate` flag to skip updating `go.mod` files if you only want to commit or tag changes.
    - **Best Practice**: Run the script with the default settings when deploying a new version to ensure dependencies are up-to-date. If the update is unnecessary, use `-N` to avoid redundant changes.

4. **Commit Changes Effectively**
    - Use the `-G` or `--git` option to commit changes with a specified message.
    - **Commit Example**:
      ```bash
      ./versionapply.sh -G "Updated dependencies and incremented version" -B release
      ```
    - **Best Practice**: Write clear, concise commit messages describing changes. Include details about updated dependencies or features that impact the version increment.

5. **Tag the Repository and Submodules**
    - Use the `-T` or `--tag` option to apply version tags to both the root repository and submodules.
    - Each tag will include the updated version number (from the `version` file) and will be pushed to the specified branch.
    - **Tag Example**:
      ```bash
      ./versionapply.sh -T -B release
      ```
    - **Best Practice**: Consistently use tagging to mark significant versions or releases. A robust tagging system helps with tracking versions and simplifies rollback if needed.

6. **Review Changes in Debug Mode**
    - Enable the `-D` or `--debug` flag to test the script without making any permanent changes.
    - **Example**:
      ```bash
      ./versionapply.sh -D -G "Simulated commit for version update"
      ```
    - **Best Practice**: Run the script with `--debug` first, especially if you’re running it on a new branch or for the first time in a release cycle. Debug mode helps verify changes without affecting the actual branch or repository.

7. **Combine Flags for Efficiency**
    - Combine options to perform multiple actions in one command. For example, to skip `go.mod` updates, commit changes with a message, and tag the repository:
      ```bash
      ./versionapply.sh -N -G "Minor updates for compatibility" -T -B release
      ```
    - **Best Practice**: Take advantage of combined flags to streamline version updates, especially when deploying minor releases or patches.

## Example Commands

- **Full Update with Tagging on a Release Branch**:
  ```bash
  ./versionapply.sh -G "Major update with new features" -T -B release
  ```
- **Commit Only (Skip `go.mod` Updates) on Staging Branch**:
  ```bash
  ./versionapply.sh -N -G "Minor fix" -B staging
  ```
- **Test Changes in Debug Mode**:
  ```bash
  ./versionapply.sh -D -G "Simulated commit"
  ```

## Additional Notes

- **Usage Help**: Run `./versionapply.sh -h` or `--help` to display usage instructions.
- **Version Number**: Ensure that the `version` file reflects the correct version before running the script. The script will automatically increment this version in tags and updates.

By following these best practices, you can manage versions, branches, and tags consistently and effectively, enhancing the maintainability and reliability of your project.