# v###

Nov 23 2022

- tested npx and npm within a package; 
- cleaned up symlink syncing and uninstalling

# v0.1.7

Nov 22 2022

- symlinking npm installs; 
- streaming unarchive; 
- removes NVM_BIN; 
- changes nvm-shim to use CombinedOutput for stderr messages;


# v0.1.6

Nov 21 2022

- Adds publish.sh script for tagging
- Ability to read from local .nvmrc with `nvm install` and `nvm use`
- Adds `nvm install --force` to forceably re-install an existing version
- Adds non-specific check in `nvm use` to see if we've already installed a fitting version
- `nvm list` is now ordered DESC
