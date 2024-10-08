#!/bin/bash
set -ex
now="$(date)"
printf "Running at %s\n" "$now"

make search
make index
make meta
make output
cp -R output/* ../100kb-out/
cd ../100kb-out && git add . && git commit -m "$now" && git push