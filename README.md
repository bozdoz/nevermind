# nevermind

a node version manager

A play on `nvm`.  Inspired by both [nvm-sh](https://github.com/nvm-sh/nvm) and [nvm-windows](https://github.com/coreybutler/nvm-windows), which were written in shell and go, respectively.

## Installation (WIP)

Currently in development.  

Requirements:

1. Docker
2. VSCode

Getting Started:

1. Use/Open in VSCode devcontainer (.devcontainer directory) extension
2. Run the go package download prompts

You might need to build the executable:

```bash
go build -o ~/.nevermind/bin ./nvm-shim
```

Or with `go generate`:

```bash
go generate ./...
```

Make sure this directory is in your `PATH` (perhaps via `.bashrc`):

```bash
export PATH="$HOME/.nevermind/bin:$PATH"
```

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

If this is all set up, you should be able to run:

```bash
node -v
```

### Remaining tasks

- github actions for building and generating releases
  - no idea if this is what I want
- install script (bash?)
  - I want a way to automatically build nvm-shim, update PATH, create binary symlinks on installation
- test local npm installs in a project

nvm tasks:
- write nvm install script for windows (extracting zip at minimum)
- tests

nvm-shim tasks:
- tests
