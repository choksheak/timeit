#/bin/bash

echo Running go fmt
go fmt github.com/choksheak/timeit
go fmt github.com/choksheak/timeit/timeme

if [ !`which golint` ]; then
  go get -u github.com/golang/lint/golint
fi

echo Running golint
golint github.com/choksheak/timeit
golint github.com/choksheak/timeit/timeme

echo Building executables
go install github.com/choksheak/timeit
go install github.com/choksheak/timeit/timeme
