@echo Running go fmt
@go fmt github.com\choksheak\timeit
@go fmt github.com\choksheak\timeit\timeme

@REM Install golint
@REM go get -u github.com/golang/lint/golint

@echo Running golint
@golint github.com\choksheak\timeit
@golint github.com\choksheak\timeit\timeme

@echo Building executables
@go install github.com\choksheak\timeit
@go install github.com\choksheak\timeit\timeme
