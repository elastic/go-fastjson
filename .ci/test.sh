#!/usr/bin/env bash
set -euxo pipefail

## beginning go_import_path
mkdir -p "${GOPATH}/src/go.elastic.co/fastjson"
cp -Rf . "${GOPATH}/src/go.elastic.co/fastjson"
cd $GOPATH/src/go.elastic.co/fastjson
go get
cd -
## end go_import_path

# Run the tests
set +e
export OUT_FILE="build/test-report.out"
mkdir -p build
go test -v 2>&1 | tee ${OUT_FILE}
status=$?

go get -v -u github.com/jstemmer/go-junit-report
go-junit-report > "build/junit-${GO_VERSION}.xml" < ${OUT_FILE}

exit ${status}