# Release process

Follow these steps to create a new release for version `vX.Y.Z`.

**NOTE!** Include `v` in version.

## 1. Set version
Open [main.go](./main.go) and set the version.

## 2. Update the changelog
Open [CHANGELOG.md](./CHANGELOG.md) and add a new release heading for the version.

## 3. Commit and push
```shell
$ git commit -am "Release vX.Y.Z"
$ git push
```

## 3. Create assets
```shell
$ ./build.sh
```

## 4. Use GitHub CLI in order to create release
```shell
$ gh release create vX.Y.Z bin/*
```
