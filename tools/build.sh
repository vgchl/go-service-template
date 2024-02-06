#!/usr/bin/env bash
set -eE

export DOCKER_SCAN_SUGGEST=false

function task_clean {
  ui_header "clean"
  rm -f mind_service
  rm -rf ./proto/gen
  ui_done
}

function task_build_proto {
  ui_header "build_proto"
  (
    # shellcheck disable=SC2164
    cd proto
    echo "Formatting..."
    $buf format --write
    echo "Linting..."
    $buf lint
    echo "Generating source code..."
    $buf generate
    echo "Checking breaking changes..."
    $buf breaking --against "../.git#branch=main,subdir=proto"
  )
  echo "Testing..."
  $go test ./proto/...
  ui_done
}

function task_build_go {
  ui_header "build_go"
  $go mod tidy
  version=$(git describe --tags --match="v[0-9]*.[0-9]*.[0-9]*" --exclude="v*[^0-9.]*" || echo "v0.0.0")
  $go build -ldflags="-X 'mind-service/app.Version=${version}'"
  ui_done
}

function task_build_docker {
  ui_header "build_docker"
  docker build -t mind-service .
  ui_done
}

function task_lint {
  ui_header "lint"
  lint="golangci-lint"
  if ! command -v "golangci-lint" &> /dev/null; then
    lint="docker compose run --no-TTY --rm golangci-lint"
  fi
  # shellcheck disable=SC2068
  $lint run $@
  ui_done
}

function task_test {
  ui_header "test"
  $go test ./...
  ui_done
}

### Utils ###

function ui_header {
  echo -e "─── $1 ───"
}

function ui_separator {
  echo
}

function ui_done {
  echo -e "${COLOR_F_GREEN}Done${COLOR_RESET}"
  ui_separator
}

function ui_trapped_error {
  echo -e "${COLOR_F_RED}Failed${COLOR_RESET}"
  exit 1
}

COLOR_F_GREEN=$(tput setaf 2)
COLOR_F_RED=$(tput setaf 1)
COLOR_RESET=$(tput sgr0)

trap ui_trapped_error ERR

function fallback {
  if command -v "$1" &> /dev/null; then
    echo "$1"
  else
    echo "$2"
  fi
}

go=$(fallback "go" "docker compose run --rm go")
buf=$(fallback "buf" "docker compose run --rm buf")
