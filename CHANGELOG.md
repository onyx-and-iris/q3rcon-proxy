# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Before any major/minor/patch bump all unit tests will be run to verify they pass.

## [Unreleased]

-   [x]

## [0.6.0] - 2024-03-21

### Added

-   new environment variable `Q3RCON_DEBUG` for enabling debug logging. Defaults to 0.
-   rcon responses are now logged at debug level
-   invalid responses (rcon and query) now logged

### Changed

-   All packet header checking methods moved into Session struct.

### Fixed

-   a bug causing the proxy not to send back query responses

## [0.3.0] - 2024-03-08

### Added

-   outgoing rcon requests now logged at info level
-   new environment variable `Q3RCON_HOST` for specifying which ip to bind the proxy to. Defaults to `0.0.0.0`.

### Changed

-   now using [logrus][logrus] package for logging.

### Fixed

-   a `slice bounds out of range` error due to query packets being logged.

## [0.1.0] - 2024-01-27

### Added

-   only forward packets if the header matches q3 rcon/query.

## [0.0.1] - 2024-01-27

### Added

-   All source files for lilproxy including full commit history.

[logrus]: https://github.com/sirupsen/logrus
