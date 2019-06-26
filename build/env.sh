#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
orgdir="$workspace/src/github.com/klaytn"
if [ ! -L "$orgdir/klaytn" ]; then
    mkdir -p "$orgdir"
    cd "$orgdir"
    ln -s ../../../../../. klaytn
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$orgdir/klaytn"
PWD="$orgdir/klaytn"

# Launch the arguments with the configured environment.
exec "$@"
