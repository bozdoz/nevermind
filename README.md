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

After `install`, and `use`, you might need to build the executable:

```bash
go build -o ~/.nevermind/bin ./nvm-shim
```

Make sure this directory is in your `PATH` (perhaps via `.bashrc`):

```bash
export PATH="$HOME/.nevermind/bin:$PATH"
```

If this is all set up, you should be able to run:

```bash
node -v
```

(WIP) `npm`, `npx`, and anything installed via `npm i -g`

### Remaining tasks

- Access to `npm` and `npx` executables
  - possible now by building nvm-shim to PATH as `node`, then symlinking `npm` and `npx` to the shim
- ~~Build script for nvm-shim~~, and symlinking node, npm, and npx
- github actions for building and generating tags and packages
  - no idea if this is what I want
- install script (bash?)
  - I want a way to automatically build nvm-shim, update PATH, create binary symlinks on installation
- progress bar on node download
- search for matching node download if only major or minor numbers are given
- ability to download latest LTS
- ~~documentation~~ Maybe done with godoc
- tests
- automatically call `install` when `use` doesn't match
- automatically call `use` after `install`
- make sure there's no infinite loops of `install` and `use`
- figure out go mod versioning??
- write install script for windows (extracting zip at minimum)
- figure out global installs (e.g. npm i -g yarn)
  - it goes to node/v/bin as a symlink
- optimize downloads and untar/ungzip with streams
  - ideally download to file & ungzip & sha at the same time with io.MultiWriter (am I crazy?)
