# CloudFoundry CLI Plugin - Doctor

This plugin is obviously inspired from [brew](http://brew.sh/) doctor :) It will scan your currently `target`ed cloudfoundry space to see if there are anomalies or useful action points that it can report back to you. Current functionality is only focused on apps and routes..

This plugin does *not* change any state or configuration, it merely just scans and gathers information than reports back anomalies.

## Installation

Install pre-built plugin from <https://plugins.cloudfoundry.org>:

```plain
cf install-plugin -r CF-Community "doctor"
```

Alternatively, build and install from source:

```plain
go get github.com/cloudfoundry/cli
go get github.com/emirozer/cf-doctor-plugin
cd $GOPATH/src/github.com/emirozer/cf-doctor-plugin
go build
cf install-plugin cf-doctor-plugin
```

## Usage

To triage the current space:

```plain
cf doctor
```

To triage all available spaces in current org:

```plain
cf doctor --all-spaces
```

## Sample output

![Screenshot](https://raw.github.com/emirozer/cf-doctor-plugin/master/docs/ndoc.png)
