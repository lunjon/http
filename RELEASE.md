# Release process

Follow these steps to create a new release for version `vX.Y.Z`.

**NOTE!** Include `v` in version.

## 1. Set version
Open [main.go](./main.go) and set the version.

## 2. Update the changelog
Open [CHANGELOG.md](./CHANGELOG.md) and add a new release heading for the version.

## 3. Run justfile target `release`
```shell
$ just release "vX.Y.Z"
```
