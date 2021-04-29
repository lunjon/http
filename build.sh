#!/bin/bash

platforms=(
    "windows/amd64"
    "darwin/amd64"
    "linux/amd64"
)

echo "Building executable for:"
for platform in "${platforms[@]}"
do
    split=(${platform//\// })
    GOOS=${split[0]}
    GOARCH=${split[1]}
    
    output="bin/http.$GOOS.$GOARCH"
    if [ $GOOS = "windows" ]; then
        output+=".exe"
    fi

    echo "  - $platform"
    env \
        GOOS=$GOOS \
        GOARCH=$GOARCH \
        go build -o $output

    if [ $? -ne 0 ]; then
        echo 'Error...'
        exit 1
    fi
done
