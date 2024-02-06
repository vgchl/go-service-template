#!/usr/bin/env bash
set -eE


function main {
    trap ui_trapped_error ERR

    go="go"
    if ! has_local "go"; then
      go="docker compose run --rm go"
    fi
    buf="buf"
    if ! has_local "buf"; then
      buf="docker compose run --rm buf"
    fi

    tasks="${@:0}"
    if [ "$tasks" == "" ]; then
      tasks="clean build lint test"
    fi

    if has_task "clean" $tasks; then
      task_clean; ui_separator
    fi
    if has_task "build" $tasks; then
      task_build; ui_separator
    fi
    if has_task "build-proto" $tasks; then
      task_build_proto; ui_separator
    fi
    if has_task "build-go" $tasks; then
      task_build_go; ui_separator
    fi
    if has_task "build-docker" $tasks; then
      task_build_docker; ui_separator
    fi
    if has_task "lint" $tasks; then
      task_lint; ui_separator
    fi
    if has_task "test" $tasks; then
      task_test; ui_separator
    fi
}

### Tasks ###

function task_clean {
  ui_header "Clean"
  rm -f mind_service
  rm -rf ./proto/gen
  ui_done
}

function task_build {
    task_build_proto
    ui_separator
    task_build_go
}

function task_build_proto {
    ui_header "Build Protobuf"
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
    ui_header "Build Go"
    $go mod tidy
    version=$(git describe --tags --match="v[0-9]*.[0-9]*.[0-9]*" --exclude="v*[^0-9.]*" || echo "v0.0.0")
    $go build -ldflags="-X 'mind-service/app.Version=${version}'"
    ui_done
}

function task_build_docker {
  ui_header "Build Docker"
  docker build -t mind-service .
  ui_done
}

function task_lint {
  ui_header "Lint"
  lint="golangci-lint"
  if ! command -v "golangci-lint" &> /dev/null; then
    lint="docker compose run --no-TTY --rm golangci-lint"
  fi
  $lint run
  ui_done
}

function task_test {
    ui_header "Test Go"
    $go test ./...
    ui_done
}

function has_task {
  for i in "${@:2}"; do
    if [ "$i" == "$1" ]; then
      return 0
    fi
  done
  return 1
}

### Utils ###

function has_local {
    if ! command -v "$1" &> /dev/null
    then
        return 1
    fi
    return 0
}

function ui_header {
  echo -e "─── $1 ───"
}

function ui_separator {
  echo
}

function ui_done {
  echo -e "${COLOR_F_GREEN}Done${COLOR_RESET}"
}

function ui_trapped_error {
  echo -e "${COLOR_F_RED}Failed${COLOR_RESET}"
  exit 1
}

COLOR_F_GREEN=$(tput setaf 2)
COLOR_F_RED=$(tput setaf 1)
COLOR_RESET=$(tput sgr0)

main "$@"
