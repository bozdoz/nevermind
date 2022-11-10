# nevermind

A play on `nvm`.  Inspired by both [nvm-sh](https://github.com/nvm-sh/nvm) and [nvm-windows](https://github.com/coreybutler/nvm-windows), which were written in shell and go, respectively.

## Installation (WIP)

Currently in development.  

Requirements:

1. Docker
2. VSCode

Getting Started:

1. Use/Open in VSCode devcontainer (.devcontainer directory) extension
2. Run the go package download prompts
3. Call the nvm scripts directly with `go run`:

```bash
go run ./nvm install 16.0.0
```

```bash
go run ./nvm use 16.0.0
```

This should have created a `~/.nevermind/config.json` file pointing to the version you've set to `use`, and installed that version to `~/.nevermind/node/`.

Run with `DEBUG=1` to output debug logs:

```bash
DEBUG=1 go run ./nvm install 16.0.0
```

After `install`, and `use`, you can call the executable:

```bash
go run ./nvm-shim -v
```

```bash
go run ./nvm-shim -e "console.log('hello from node')"
```

### Remaining tasks

- Access to `npm` and `npx` executables
- Build script for nvm-shim, and symlinking node, npm, and npx
- build script in general (nothing is built!)
- github actions for building and generating tags and packages
- install script (bash?)
- progress bar on node download
- search for matching node download if only major or minor numbers are given
- ability to download latest LTS
- documentation
- tests
- automatically call `install` when `use` doesn't match
- automatically call `use` after `install`
- make sure there's no infinite loops of `install` and `use`
- download timeouts
- handle 404 not found from installs
- figure out go mod versioning
- write install script for windows (extracting zip at minimum)