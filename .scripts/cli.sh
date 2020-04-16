#! /bin/bash

TARGETS="linux/amd64 linux/386 linux/arm linux/arm64 darwin/amd64 freebsd/amd64 freebsd/386 freebsd/arm freebsd/arm64 openbsd/amd64 openbsd/386 openbsd/arm openbsd/arm64"
DATE=`date +%Y.%m.%d.%H%M%S`
VERSION=`git describe --tags`

# Cross-compile Temporal using gox, injecting appropriate tags.
go get -u github.com/mitchellh/gox

rm -rf release
mkdir -p release

gox -output="release/tex-cli-$(git describe --tags)-{{.OS}}-{{.Arch}}" \
    -ldflags "-X main.Version=$VERSION -X main.CompileDate=$DATE -s -w" \
    -osarch="$TARGETS" \
    ./cmd/tex

ls ./release/tex-cli* > files
for i in $(cat files); do
    sha256sum "$i" > "$i.sha256"
done
