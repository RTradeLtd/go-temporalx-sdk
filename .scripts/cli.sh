#! /bin/bash

TARGETS="linux/amd64 linux/386 linux/arm linux/arm64 darwin/amd64 freebsd/amd64 freebsd/386 freebsd/arm freebsd/arm64 openbsd/amd64 openbsd/386 openbsd/arm openbsd/arm64"

# Cross-compile Temporal using gox, injecting appropriate tags.
go get -u github.com/mitchellh/gox

rm -rf release
mkdir -p release

gox -output="release/tex-cli-$(git describe --tags)-{{.OS}}-{{.Arch}}" \
    -ldflags "-X main.Version=$RELEASE -X main.PublicKey=$PUBLIC_KEY -s -w" \
    -osarch="$TARGETS" \
    ./cmd/tex

ls ./release/tex-cli* > files
for i in $(cat files); do
    sha256sum "$i" > "$i.sha256"
done
