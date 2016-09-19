# cMeter

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

# (deprecated) mockapi server stuff
mockapi:
  # address to host the mockapi http server
  addr: 'localhost:9090'

# container tracking stuff
tracking:
  # the label used to determine if a container should be tracked and metered
  tracking_label: 'com.example.track'

  # the label used as the service key value when sending container reports.
  key_label: 'com.example.track'

# stats collection stuff
collector:
  # The rate at which the collector polls for container stats.
  rate: 1800

# The reporting driver and driver parameters
# parameterless form
reporting: 'mock'

# or with driver parameters
reporting:
  mock:
    parameter1: 'foo'

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

#### Versioned build

Use `make` to create a versioned build:

> $ make compile

The default version is a semver-compatible string made up of the contents of the `/VERSION` file and the short form of the current git hash (e.g: `1.0.0-c63076f`). To override this default version, you may use the `BUILD_VERSION` environment variable to set it manually:

> $ BUILD_VERSION=7.7.7-lucky make compile

## Testing

There are currently no tests but this would likely use `make` or simply `go test ./...`

