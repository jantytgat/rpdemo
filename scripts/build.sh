#!/usr/bin/env bash

VERSION_COMMIT=$(git rev-list --abbrev-commit --abbrev=8 -1 HEAD)
GIT_REV_NAME=$(git name-rev HEAD --name-only)
if [[ $GIT_REV_NAME != "undefined" ]]
then
  VERSION_SEMVER=$(git name-rev HEAD --name-only)
else
  VERSION_SEMVER=$(git name-rev HEAD --name-only --tags)
fi

REGEX="^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-((0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(\+([0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?$"
if [[ $VERSION_SEMVER =~ $REGEX ]]
then
  VERSION_PRERELEASE="${BASH_REMATCH[6]}"
fi

if [[ VERSION_PRERELEASE == "" ]]
then
  echo "$VERSION_SEMVER did not match x.y.z-<prerelease>+<metadata>"
  VERSION_PRERELEASE=stable
fi

VERSION_DATE=$(date '+%Y%m%d')

CMD=$1
GOOS=$2
GOARCH=$3

FILENAME=$CMD
if [[ $GOOS == "windows" ]]
then
  FILENAME=$CMD".exe"
fi

BUILDVERSION="$VERSION_SEMVER+$VERSION_COMMIT.$VERSION_DATE"


# shellcheck disable=SC2097
GOOS=$GOOS \
GOARCH=$GOARCH \
go build \
  -ldflags="-w -s
  -X main.BuildVersion=$BUILDVERSION" \
  -o build/bin/"$VERSION_PRERELEASE"/"$VERSION_SEMVER"/"$GOOS"_"$GOARCH"/"$FILENAME" \
  cmd/"$CMD"/main.go