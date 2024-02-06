SHELL=/usr/bin/env bash

build: install clean build_proto build_go build_docker lint test

clean:
	@. tools/build.sh; task_clean

build_proto:
	@. tools/build.sh; task_build_proto

build_go:
	@. tools/build.sh; task_build_go

build_docker:
	@. tools/build.sh; task_build_docker

lint:
	@. tools/build.sh; task_lint

lint-fix:
	@. tools/build.sh; task_lint --fix

test:
	@. tools/build.sh; task_test

run:
	@. tools/build.sh; task_run

install:
	@ln -sf ../../tools/pre-commit.sh .git/hooks/pre-commit