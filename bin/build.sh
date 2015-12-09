#!/bin/bash

if [[ "$(which gox)X" == "X" ]]; then
    echo "Please install gox. https://github.com/mitchellh/gox#readme"
    exit 1
fi

rm -f doctor_plugin*

gox -os linux -os windows -arch 386 --output="doctor_plugin_{{.OS}}_{{.Arch}}"
gox -os darwin -os linux -os windows -arch amd64 --output="doctor_plugin_{{.OS}}_{{.Arch}}"

rm -rf out
mkdir -p out
mv doctor_plugin* out/
