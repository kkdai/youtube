#!/bin/bash
STATUS=0
  checking() {
    local diff
    if ! diff="$(git diff -U1 --color --exit-code)"; then
      printf '\e[31mError: running `\e[1m%s\e[22m` results in modifications that you must check it again:\e[0m\n%s\n\n' "$*" "$diff" >&2
      git checkout -- .
      STATUS=1
    fi
  }
#          checking go fmt ./...
#          exit $STATUS



while read -r file linter msg; do
  IFS=: read -ra f <<<"$file"
  printf '::error file=%s,line=%s,col=%s::%s\n' "${f[0]}" "${f[1]}" "${f[2]}" "[$linter] $msg"
  STATUS=1
done < <(bin/golangci-lint run --out-format tab)
