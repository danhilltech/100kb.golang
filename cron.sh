#!/bin/bash
set -x
now="$(date)"
printf "Running at %s\n" "$now"

make search
make index
make meta