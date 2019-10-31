#!/bin/bash

function finish {
  # Your cleanup code here
  popd
}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

pushd $DIR/..

trap finish EXIT

SUBCOMMAND=$1
ADDTIONAL_OPTIONS=""
GO=${GO:-latest}

case "$SUBCOMMAND" in
    linux-i386)
        TARGET="linux/386"
        shift
        ;;
    linux-x86_64)
        TARGET="linux/amd64"
        shift
        ;;
    linux-arm-5)
        TARGET="linux/arm-5"
        shift
        ;;
    linux-arm-6)
        TARGET="linux/arm-6"
        shift
        ;;
    linux-arm-7)
        TARGET="linux/arm-7"
        shift
        ;;
    linux-arm64)
        TARGET="linux/arm64"
        shift
        ;;
    linux-mips)
        TARGET="linux/mips"
        ADDITIONAL_OPTIONS="--ldflags '-extldflags \"-static\"'"
        shift
        ;;
    linux-mipsle)
        TARGET="linux/mipsle"
        ADDITIONAL_OPTIONS="--ldflags '-extldflags \"-static\"'"
        shift
        ;;
    linux-mips64)
        TARGET="linux/mips64"
        ADDITIONAL_OPTIONS="--ldflags '-extldflags \"-static\"'"
        shift
        ;;
    linux-mips64le)
        TARGET="linux/mips64le"
        ADDITIONAL_OPTIONS="--ldflags '-extldflags \"-static\"'"
        shift
        ;;
    darwin-386)
        TARGET="darwin/386"
        shift
        ;;
    darwin-amd64)
        TARGET="darwin-10.10/amd64"
        shift
        ;;
    windows-386)
        TARGET="windows/386"
        shift
        ;;
    windows-amd64)
        TARGET="windows/amd64"
        shift
        ;;
    *)
        echo "Undefined architecture for cross-compile. Supported architectures: linux-i386, linux-x86_64, linux-arm-5, linux-arm-6, linux-arm-7, linux-arm64, linux-mips, linux-mipsle, linux-mips64, linux-mips64le, darwin-386, darwin-amd64, windows-386, windows-amd64"
        echo "Usage: ${0} {arch} [kcn|kpn...]"
        echo "    ${0} linux-386"
        echo "    ${0} darwin-386 kcn"
        echo "    ${0} windows-386 kcn kpn"
        exit 1
        ;;
esac

echo "make $*"
GOFLAGS= GO111MODULE=off BUILD_PARAM="xgo -- --go=${GO} --targets=${TARGET} -v" make $*
