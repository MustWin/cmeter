# cMeter

![cMeter logo](https://github.com/MustWin/cmeter/blob/master/docs/logo/cmeter-logo-title.png)

[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

cMeter (Container Meter) provides container hosts with a container metering solution.

## Installation

> $ go get github.com/MustWin/cmeter

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

## Bugs and Feedback

If you see a bug or have a suggestion, feel free to open an issue [here](https://github.com/MustWin/cmeter/issues).

## Contributions

PR's welcome! There are no strict style guidelines, just follow best practices and try to keep with the general look & feel of the code present. All submissions should at least be `go fmt -s` and have a test to verify *(if applicable)*.

For details on how to extend and develop cSense, see the [dev documentation](docs/development/).

## License

The MIT License (MIT)
Copyright (c) 2016 MustWin, LLC

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.