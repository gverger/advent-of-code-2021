#!/usr/bin/env bash
set -e

new_dir=$1

mkdir $new_dir
cd $new_dir

# Gopls with vim don't like when the workspace and the folder with .git are not the same
# My personnal vim uses this file to define the root level of a project
touch .vim-go-workspace

go mod init github.com/gverger/advent2021/$new_dir
go mod edit -replace github.com/gverger/advent2021/utils=../utils

cat << EOF > main.go
package main

import "github.com/gverger/advent2021/utils"

func main() {
  utils.Main(run)
}

func run(lines []string) error {
  return nil
}
EOF

go mod tidy
