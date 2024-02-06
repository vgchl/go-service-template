

build:
	@./tools/build.sh

run:
	@go run .

lint:
	@./tools/build.sh lint

lint-fix:
	@./tools/build.sh lint --fix

install-githooks:
	ln -s tools/pre-commit.sh .git/hooks/pre-commit