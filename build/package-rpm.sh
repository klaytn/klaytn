#!/bin/bash

DAEMON_BINARIES=(kcn kpn ken kbn kscn kspn ksen)
BINARIES=(kgen homi)

set -e

function printUsage {
    echo "Usage: $0 [-b] <target binary>"
    echo "               -b: use baobab configuration."
    echo "  <target binary>: kcn | kpn | ken | kbn | kscn | kspn | ksen | kgen | homi"
    exit 1
}

# Parse options.
BAOBAB_FLAG=
BAOBAB_PREFIX=
while getopts "b" opt; do
	case ${opt} in
		b)
			echo "Using baobab configuration..."
			BAOBAB_FLAG=" --baobab"
			BAOBAB_PREFIX="-baobab"
			;;
	esac
done
shift $((OPTIND -1))

# Parse target
TARGET=$1
if [ -z "$TARGET" ]; then
    echo "specify target binary: ${DAEMON_BINARIES[*]} ${DAEMON[*]}"
    printUsage
fi

MYDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd $MYDIR/..
function finish {
  # Your cleanup code here
  popd
}
trap finish EXIT

KLAYTN_VERSION=$(go run build/rpm/main.go version)
KLAYTN_RELEASE_NUM=$(go run build/rpm/main.go release_num)
PLATFORM_SUFFIX=$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)

PACK_NAME=
PACK_VERSION=

# Search the target name in DAEMON_BINARIES
for b in ${DAEMON_BINARIES[*]}; do
    if [ "$TARGET" == "$b" ]; then
        PACK_NAME=${b}-${PLATFORM_SUFFIX}
        PACK_VERSION=${b}d${BAOBAB_PREFIX}-${KLAYTN_VERSION}
    fi
done

# Search the target name in BINARIES
for b in ${BINARIES[*]}; do
    if [ "$TARGET" == "$b" ]; then
        PACK_NAME=${b}-${PLATFORM_SUFFIX}
        PACK_VERSION=${b}${BAOBAB_PREFIX}-${KLAYTN_VERSION}
    fi
done

# If not found from both DAEMON_BINARIES and BINARIES, exit.
if [ -z "$PACK_NAME" ]; then
    echo "specify target binary: ${DAEMON_BINARIES[*]} ${DAEMON[*]}"
    printUsage
fi

# Go for packaging!
mkdir -p ${PACK_NAME}/rpmbuild/{SPECS,SOURCES,BUILDROOT}
go run build/rpm/main.go gen_spec $BAOBAB_FLAG --binary_type $TARGET > ${PACK_NAME}/rpmbuild/SPECS/${PACK_VERSION}.spec
git archive --format=tar.gz --prefix=${PACK_VERSION}/ HEAD > ${PACK_NAME}/rpmbuild/SOURCES/${PACK_VERSION}.tar.gz
echo "rpmbuild --buildroot ${MYDIR}/../${PACK_NAME}/rpmbuild/BUILDROOT -ba ${PACK_NAME}/rpmbuild/SPECS/${PACK_VERSION}.spec"
HOME=${MYDIR}/../${PACK_NAME}/ rpmbuild --buildroot ${MYDIR}/../${PACK_NAME}/rpmbuild/BUILDROOT -ba ${PACK_NAME}/rpmbuild/SPECS/${PACK_VERSION}.spec
