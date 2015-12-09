#!/bin/bash

set -e

(cf uninstall-plugin "DoctorPlugin" || true) && go build -o doctor-plugin main.go && cf install-plugin doctor-plugin
