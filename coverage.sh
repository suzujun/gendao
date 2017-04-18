#!/usr/bin/env bash

set -e
echo "" > coverage.txt
SECONDS=0

for d in $(go list ./... | grep -v vendor | grep -v src); do
    VALENCIA_MEDIA_API_TEST=1 go test -coverprofile=coverage.out $d
    if [ -f coverage.out ]; then
        cat coverage.out >> coverage.txt
        # grep -vE '_gen|_mock|_slice' coverage.out >> coverage.txt
        rm coverage.out
    fi
done

sed -i -e "1,1d" coverage.txt # 1行目は空行なので削除

echo "test elapsed time: ${SECONDS}sec"

