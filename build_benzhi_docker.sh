#!/bin/sh
set -e
docker build -f benzhi.Dockerfile -t go-data-bug734:local .
