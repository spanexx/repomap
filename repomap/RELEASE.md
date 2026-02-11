# How to Release Repomap

This guide explains how to build and publish the `@spanexx/repomap` package to the npm registry.

## Prerequisites
-   **Node.js & npm**: Ensure `npm` is installed and authenticated.
-   **Go**: Required to cross-compile the binaries (`go build`).
-   **Make**: Required to run the release script.

## Publishing Steps

1.  **Login to NPM**
    If you haven't already, login to your npm account:
    ```bash
    npm login
    ```

2.  **Verify Build (Optional)**
    You can manually verify that the binaries build successfully:
    ```bash
    make release
    ```
    This should populate the `build/` directory.

3.  **Publish**
    To publish the package with public access:
    ```bash
    npm publish --access public
    ```
    *Note: The `prepack` script in `package.json` will automatically run `make release` before publishing to ensure the binaries are up-to-date.*

## Versioning
To update the version:
1.  Update `version` in `package.json`.
2.  Update `const version` in `cmd/repomap/main.go`.
3.  Update `Makefile` version if hardcoded (currently uses `git describe`).
4.  Commit and tag the release in git:
    ```bash
    git tag v0.2.0
    git push origin v0.2.0
    ```
