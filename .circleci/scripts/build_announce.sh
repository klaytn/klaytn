#!/usr/bin/env bash
set -euox

VERSION=$(go run build/rpm/main.go version)~${CIRCLE_TAG##*-}

SHORT_SHA1=${CIRCLE_SHA1:0:7}
CIRCLE_PR=$(hub pr list -s open -L 10 -f "%I %sH %n" | grep $CIRCLE_SHA1)
CIRCLE_PR_NUMBER=${CIRCLE_PR%% *}

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
    BAOBAB_LINUX_PACKAGE_LINKS+="<a href="${PACKAGE_PREFIX}/$i-baobab-${VERSION}-0-linux-amd64.tar.gz">$i</a>"
    BAOBAB_LINUX_PACKAGE_LINKS+=" "
    BAOBAB_DARWIN_PACKAGE_LINKS+="<a href="${PACKAGE_PREFIX}/$i-baobab-${VERSION}-0-darwin-10.10-amd64.tar.gz">$i</a>"
    BAOBAB_DARWIN_PACKAGE_LINKS+=" "
done

COMMENT_ROWS="<ul><li>Linux: ${LINUX_PACKAGE_LINKS}</li><li>Darwin: ${DARWIN_PACKAGE_LINKS}</li><li>Baobab-linux: ${BAOBAB_LINUX_PACKAGE_LINKS}</li><li>Baobab-darwin: ${BAOBAB_DARWIN_PACKAGE_LINKS}</li></ul>"

COMMENT_HEAD="Builds ready [${SHORT_SHA1}]"
COMMENT_BODY="<details><summary>${COMMENT_HEAD}</summary>${COMMENT_ROWS}</details>"

POST_COMMENT_URI="https://api.github.com/repos/klaytn/klaytn/issues/${CIRCLE_PR_NUMBER}/comments"

curl -i --request POST ${POST_COMMENT_URI} \
-H 'Content-Type: application/json' \
-H 'User-Agent: klaytnbot' \
-H 'Authorization: token '"${GITHUB_TOKEN}"'' \
--data '{"body": "'"${COMMENT_BODY}"'"}'
