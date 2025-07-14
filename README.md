# nevermind

a node version manager

A play on `nvm`. Inspired by both [nvm-sh](https://github.com/nvm-sh/nvm) and [nvm-windows](https://github.com/coreybutler/nvm-windows), which were written in shell and go, respectively.

## Installation (WIP)

Currently in development.

Build the executable:

```bash
go generate ./...
```

Or directly:

```bash
go build -o ~/.nevermind/bin ./nvm-shim
```

Make sure this directory is in your `PATH` (perhaps via `.bashrc`):

```bash
export PATH="$HOME/.nevermind/bin:$PATH"
```

OR run directly:

Install:

```bash
go run ./nvm install 22
```

Use:

```bash
go run ./nvm use 22
```

This should have created a `~/.nevermind/config.json` file pointing to the version you've set to `use`, and installed that version to `~/.nevermind/node/`.

Run with `DEBUG=1` to output debug logs:

```bash
DEBUG=1 go run ./nvm install 22
```

If this is all set up, you should be able to run:

```bash
node -v
```

Check your config with:

```bash
go run ./nvm config
```

### Remaining tasks

- github actions for building and generating releases
  - https://github.com/softprops/action-gh-release
- install script (bash?)
  - I want a way to automatically build nvm-shim, update PATH, create binary symlinks on installation

nvm tasks:

- write nvm install script for windows (extracting zip at minimum)
- tests

nvm-shim tasks:

- tests
