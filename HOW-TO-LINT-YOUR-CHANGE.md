# How to lint your change

This document describe how to setup automatically or manually linting your change.

## Prerequisites
- `gofumpt` should be installed. `go install mvdan.cc/gofumpt@latest`
  - go version should be equal to or higher than v1.18.0
- `goimports` should be installed. To install it, run `go install golang.org/x/tools/cmd/goimports@latest`

## Setup Git Hook
This will apply code formating automatically when you commit to git repository. So you can pass the linting tests registered in Klaytn Circle CI.

- Copy and paste below pre-commit script to `klaytn/.git/hooks/pre-commit` file and change the file to be executable eg. `chmod +x pre-commit`

```go
#!/bin/sh
# Copyright 2012 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# git gofmt pre-commit hook
#
# To use, store as .git/hooks/pre-commit inside your repository and make sure
# it has execute permissions.
#
# This script does not handle file names that contain spaces.

gofiles=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
[ -z "$gofiles" ] && exit 0

unformatted=$(gofumpt -l $gofiles)
unimported=$(goimports -l $gofiles)

# Some files are not gofmt'd. Print message and fail.
echo >&2 "Go files must be formatted with gofmt. Below files are reformatted:"
for fn in $unformatted; do
	echo >&2 "$PWD/$fn"
	gofumpt -w $PWD/$fn
	git add $PWD/$fn
done

echo >&2 "Algin import packages. Below files are reformatted:"
for fn in $unimported; do
	echo >&2 "$PWD/$fn"
	goimports -w $PWD/$fn
	git add $PWD/$fn
done

exit 0
```

## Manually formatting
You can format codes manually by using below commands in the root of klaytn repository.

```bash
klaytn$ gofumpt -w .
klaytn$ goimports -w .
```