VERSION         :=      $(shell cat ./VERSION)
IMAGE_NAME      :=      willis7/contacts

# As a call to `make` without any arguments leads to the execution
# of the first target found I really prefer to make sure that this
# first one is a non-destructive one that does the most simple 
# desired installation. It's very common to people set it as `all`
# but it could be anything like `a`.
all: install

# Install just performs a normal `go install` which builds the source
# files from the package at `./` (I like to keep a `main.go` in the root
# that imports other subpackages). As I always commit `vendor` to `git`
# a `go install` will typically always work - except if there's an OS
# limitation in the build flags (e.g, a linux-only project).
install:
        go install -v

# keeping `./main.go` with just a `cli` and `./lib/*.go` with actual 
# logic, `tests` usually reside under `./lib` (or some other subdirectories).
# Here we could do something like `find . -name "*" -type d -exec ...` but IMO
# that's unnecessary. Just `cd`ing to what matters to you is fine - no need to
# handle the case of directories that you don't want to execute a command.
test:
        go test -v

# Just like `test`, formatting what matters. As `main.go` is in the root,
# `go fmt` the root package. Then just `cd` to what matters to you (`vendor`
# doesn't matter).
fmt:
        go fmt


# This target is only useful if you plan to also create a Docker image at
# the end. I have a separate `gist` with a sample Dockerfile tailored for
# golang that you can check out at <TODO>.
# I really like publishing a Docker image together with the GitHub release
# because Docker makes it very simple to someone run your binary without
# having to worry about the retrieval of the binary and execution of it
# - docker already provides the necessary boundaries.
image:
        docker build -t $(IMAGE_NAME) .


# This is pretty much an optional thing that I tend to always include.
# Goreleaser is a tool that allows anyone to integrate a binary releasing
# process to their pipelines. Here in this target With just a simple 
# `make release` you can have a `tag` created in GitHub with multiple
# builds if you wish. 
# See more at `gorelease` github repo.
release:
        git tag -a $(VERSION) -m "Release" || true
        git push origin $(VERSION)
        goreleaser --rm-dist

.PHONY: install test fmt release