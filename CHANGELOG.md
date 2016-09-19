# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]
### Removed
- the `api` command.
- `mockapi` configuration section.
- mock api server

## [0.0.1] - 2016-09-19
### Added
- select reporting driver in config file's `reporting` directive
- select containers driver in config file's `containers` directive
- `cmeterapi` reporting driver
- reporting driver factory registration
- containers driver factory registration
- `mock` reporting driver that just logs stuff out
- basic pipeline implementation for processing messages
- development config
- mock api server
- overridable and automatic build versioning support in Makefile
- embedded cAdvisor containers driver
- reporting driver abstraction
- containers driver abstraction
- basic README
- VERSION file
- this Changelog
- basic Makefile with `all`, `clean`, and `compile` targets