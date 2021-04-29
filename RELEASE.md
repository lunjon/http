## Release process

Follow these steps to create a new release.

### 1. Set version

Open [main.go](./main.go) and set the version.

### 2. Update the changelog

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

### 4. Create GitHub release

1. Create executables: `./build.sh`.
1. Head over to https://github.com/lunjon/http/releases and create a new release.
1. Attach the executables.
