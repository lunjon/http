## Release process

Follow these steps to create a new release.

### 1. Update version

Open [main.go](./cmd/httpreq/main.go) and set `version` to the new version.

### 2. Update the changelog.

Open [CHANGELOG.md](./CHANGELOG.md) and add a new release heading for the version.

### 3. Commit, tag and push

Run the commands below to commit, tag, and push

```shell
$ version=<version>
$ git add .
$ git commit -m "Release $version"
$ git tag -a $version
$ git push --follow-tags
```

### 4. Confirm release

Go to [travis-ci](https://travis-ci.org/lunjon/httpreq) and [github](https://github.com/lunjon/httpreq/releases) to confirm the release.

