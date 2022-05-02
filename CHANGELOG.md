# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [v1.4.0] - 2022-05-02
### Changed
- Updated dependencies.

## [v1.3.1] - 2020-06-21
### Fixed
- #16 Removed file extension for Linux executable.
- #16 Removed incorrect escaping of paths on Linux to ensure files are deleted properly.

## [v1.3.0] - 2020-03-09
### Added
- Linux build

### Changed
- #13 Added `--older-than-units` flag to specify the time units to use.  Either (d)ays, (h)ours, or (m)inutes.  For backwards compatibility the default has been set to (d)ays so if the flag is not set, nothing breaks.

## [v1.2.0] - 2019-10-26
### Added
- #11 Optional `--remove-directories` flag to remove any empty subdirectories when doing a recursive search.

## [v1.1.0] - 2019-01-20
### Added
- #9 Optional `--case-insensitive` flag to look for file extensions in any case.

## [v1.0.0] - 2018-07-24
### Added
- #7 In addition to logging to stdout, logs will now also go to the Windows Event Log

## [v0.2.2] - 2017-10-10
### Added
- Logging additional information;  `path`, `version #`, `build #`, `older-than` flag, `recursive` flag.

## [v0.2.1] - 2017-10-09
### Added
- Logging message if there are no files to be deleted.

## [v0.2.0] - 2017-10-03
### Added
- Version flag to print version and build number.
- Accept a wildcard extension to delete all files within a specified path.

## [v0.1.0] - 2017-10-02
Initial release
