# cMeter

![cMeter logo](https://github.com/MustWin/cmeter/blob/master/docs/logo/cmeter-logo-title.png)

cMeter (Container Meter) provides tenant container hosts with a container metering solution. 

## Installation

> $ go get github.com/MustWin/cmeter

### Tools

These are tools used for development.

#### Vendoring

cMeter uses `gvt` to track and vendor dependencies, install with:

> $ go get github.com/FiloSottile/gvt

For details on usage, see [the project's github](https://github.com/FiloSottile/gvt).

## Usage

> $ cmeter agent <config_path>

**Example** - with development config:

> $ cmeter agent ./config.dev.yml


## Configuration

A configuration file is *required* for cMeter but environment variables can be used to override configuration. A configuration file can be specified as a parameter or with the `CMETER_CONFIG_PATH` environment variable. 

All configuration environment variables are prefixed by `CMETER_`

A development configuration file is included: `/config.dev.yml` and a `/config.local.yml` has already been added to gitignore to be used for local testing or development.


```yaml
# configuration schema version number, only `0.1`
version: 0.1

# log stuff
logging:
  # minimum event level to log: `error`, `warn`, `info`, or `debug`
  level: 'debug'
  # log output format: `text` or `json`
  formatter: 'text'
  # custom fields to be added and displayed in the log
  fields:
    customfield1: 'value'

# container tracking stuff
tracking:
  # the label used to determine if a container should be tracked and metered
  tracking_label: 'com.example.track'

# stats collection stuff
collector:
  # The rate at which the collector polls for container stats.
  rate: 1800

# The reporting driver and driver parameters
# parameterless form
reporting: 'mock'

# or with driver parameters
reporting:
  ctoll:
    apikey: '2390511a-870d-11e6-ae22-56b6b6499611'
    endpoint: 'localhost:9090'

# Similar to the reporting driver section
containers: 'embedded'

```

Both `reporting` and `containers` only allow specification of *one* driver per configuration. Anymore will cause a validation error when the application starts.

## Building

One caveat with building currently is that because of the cAdvisor dependency for the containers driver, `cgo` *cannot* be disabled; the build will fail. So no `CGO_ENABLED=0` builds.

#### Dev/Local build

Use `go` and build from the root of the project:

> $ go build

Please note the version number displayed will be the value of `main.DEFAULT_VERSION`

#### Local versioned build

Use `make` to create a versioned build:

> $ make compile

The default version is a semver-compatible string made up of the contents of the `/VERSION` file and the short form of the current git hash (e.g: `1.0.0-c63076f`). To override this default version, you may use the `BUILD_VERSION` environment variable to set it manually:

> $ BUILD_VERSION=7.7.7-lucky make compile

#### Dist build

This is primarily meant to be used when building the docker image. Distribution builds are versioned like the local versioned builds and the build specifically targets `linux`

> $ make dist

#### Docker Image

Building a Docker Image is a two-step process because of the CGO requirement and the desire to keep a small image size. First we build the distribution binary:

> $ make dist

And then we can make the image:

> $ make image

The default image repo used is that of the Makefile's `DOCKER_REPO` variable. The image tag is the `BUILD_VERSION` variable and can be overridden as noted in the *"Local versioned build"* section above.

## Testing

Use `make` to run tests:

> $ make test

You can also use `go test` directly for any package without additional bootstrapping:

> $ go test ./agent/
