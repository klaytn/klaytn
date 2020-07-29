#!/usr/bin/env bash
set -euox

#temp
CIRCLE_SHA1=42b609ff7773b623afcaba7724636e086e41bcdc
CIRCLE_PULL_REQUEST=https://github.com/ground-x/kaikas-pixelplex/pull/307
GITHUB_TOKEN=

if [[ -z ${CIRCLE_PULL_REQUEST+x} ]]; then
    echo "No pull request detected for commit ${CIRCLE_SHA1}"
    exit
fi

VERSION=$(go run build/rpm/main.go version)

SHORT_SHA1=${CIRCLE_SHA1:0:7}

PACKAGES="kcn kpn ken kscn kspn ksen kbn kgen homi"
BAOBAB_PACKAGES="kcn kpn ken"
PACKAGE_PREFIX="http://packages.klaytn.net/klaytn/${VERSION}"

LINUX_PACKAGE_LINKS=""
DARWIN_PACKAGE_LINKS=""
BAOBAB_LINUX_PACKAGE_LINKS=""
BAOBAB_DARWIN_PACKAGE_LINKS=""

for i in ${PACKAGES}
do
    LINUX_PACKAGE_LINKS+="<a href="${PACKAGE_PREFIX}/$i-${VERSION}-0-linux-amd64.tar.gz">$i</a>"
    LINUX_PACKAGE_LINKS+=" "
    DARWIN_PACKAGE_LINKS+="<a href="${PACKAGE_PREFIX}/$i-${VERSION}-0-darwin-10.10-amd64.tar.gz">$i</a>"
    DARWIN_PACKAGE_LINKS+=" "
done

for i in ${BAOBAB_PACKAGES}
do
    BAOBAB_LINUX_PACKAGE_LINKS+="<a href="${PACKAGE_PREFIX}/$i-${VERSION}-0-linux-amd64.tar.gz">$i</a>"
    BAOBAB_LINUX_PACKAGE_LINKS+=" "
    BAOBAB_DARWIN_PACKAGE_LINKS+="<a href="${PACKAGE_PREFIX}/$i-${VERSION}-0-darwin-10.10-amd64.tar.gz">$i</a>"
    BAOBAB_DARWIN_PACKAGE_LINKS+=" "
done

COMMENT_ROWS="rpm:
linux: ${LINUX_PACKAGE_LINKS}
darwin: ${DARWIN_PACKAGE_LINKS}
baobab-linux: ${BAOBAB_LINUX_PACKAGE_LINKS}
baobab-darwin: ${BAOBAB_DARWIN_PACKAGE_LINKS}
"

COMMENT_HEAD="Builds ready [${SHORT_SHA1}]"
COMMENT_BODY="<details><summary>${COMMENT_HEAD}</summary>${COMMENT_ROWS}</details>"

#POST_COMMENT_URI="https://api.github.com/repos/klaytn/klaytn/issues/${CIRCLE_PR_NUMBER}/comments"
POST_COMMENT_URI="https://api.github.com/repos/whoisxx/klaytn/issues/${CIRCLE_PR_NUMBER}/comments"

curl -i \
-H "Accept: application/json"
#-H "Content-Type: application/json"
-H "Content-Type: application/x-www-form-urlencoded"
-H "User-Agent: klaytnbot"
-H "Authorization: token ${GITHUB_TOKEN}"
-X POST --data "${COMMENT_BODY}" ${POST_COMMENT_URI}
