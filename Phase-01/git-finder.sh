#!/usr/bin/env bash
set -e

toplevel_path="$(git rev-parse --show-toplevel)"
branches=$(git branch --format='%(refname:short)') 
git grep "${1:-TODO}" ${branches} -- "${toplevel_path}"

