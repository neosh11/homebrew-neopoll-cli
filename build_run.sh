#!/bin/bash
mkdir -p .out
go build -o .out/neopoll .
.out/neopoll "$@"