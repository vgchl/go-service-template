#!/usr/bin/env bash
set -E

trap ui_trapped_error ERR

function main {
    go="go"
    if ! has_local "go"; then
      go="docker compose run --rm go"
    fi
    buf="buf"
    if ! has_local "buf"; then
      buf="docker compose run --rm buf"
    fi

    clean
    ui_separator
    build
    ui_separator
    test
}

function clean {
  ui_header "Clean"
  rm -f service
  rm -rf ./proto/gen
  ui_done
}

function build {
    build_proto
    ui_separator
    build_go
}

function build_go {
    ui_header "Build Go"
    $go mod tidy
    $go build -o service
    ui_done
}

function build_proto {
    ui_header "Build Protobuf"
    echo "Formatting..."
    (cd proto && $buf format --write)
    echo "Linting..."
    (cd proto && $buf lint)
    echo "Generating source code..."
    (cd proto && $buf generate)
    echo "Checking breaking changes..."
    (cd proto && $buf breaking --against ".git#branch=main")
    echo "Testing..."
    $go test ./proto/...
    ui_done
}

function test {
    ui_header "Test Go"
    $go test ./...
    ui_done
}

function has_local {
    if ! command -v "$1" &> /dev/null
    then
        return 1
    fi
    return 0
}

function ui_header {
  echo -e "─── $1$(tput sgr0) ───"
}

function ui_separator {
  echo
}

function ui_done {
  echo -e "${COLOR_F_GREEN}Done${COLOR_RESET}"
}

function ui_trapped_error {
  echo -e "${COLOR_F_GREEN}Failed${COLOR_RESET}"
  exit 1
}

COLOR_F_GREEN=$(tput setaf 2)
COLOR_RESET=$(tput sgr0)

main
