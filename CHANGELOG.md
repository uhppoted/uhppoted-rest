# CHANGELOG

## Unreleased


## [0.9.0](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.9.0) - 2026-01-27

### Updated
1. Updated to Go 1.25.


## [0.8.11](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.11) - 2025-07-01

### Added
1. Added `get/set-antipassback` API function to retrieve or set a controller antipassback mode.

### Updated
1. Bumped to Go 1.24.


## [0.8.10](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.10) - 2025-01-30

### Added
1. ARMv6 build target (RaspberryPi ZeroW).


## [0.8.9](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.9) - 2024-09-06

### Added
1. TCP/IP support.

### Updated
1. Updated to Go 1.23.


## [0.8.8](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.8) - 2024-03-27

### Added
1. `restore-default-parameters` API function to reset a controller to the manufacturer default 
    configuration.
2. Added public Docker image to ghcr.io.

### Updated
1. Bumped Go version to 1.22.


## [0.8.7](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.7) - 2023-12-30

### Added
1. `set-door-passcodes` API function to set the supervisor passcodes for a door.

### Updated
1. Replaced `nil` event pointer in `get-status` with zero value.


## [0.8.6](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.6) - 2023-08-30

### Added
1. `activate-keypads` API function to activate/deactivate reader access keypads.


## [0.8.5](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.5) - 2023-06-13

### Added
1. `set-interlock` API function to set controller door interlock mode.

### Updated
1. Replaced card `From` and `To` field pointers with zero values.
2. Reworked `put-acl` to discard cards with invalid _start_ and _end_ dates.
3. Updated OpenAPI specification.


## [0.8.4](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.4) - 2023-03-17

### Added
1. `doc.go` package overview documentation.
2. Reworked logging to use uhppoted-lib/log package.
3. Added PIN support to get-card and put-card APIs.

### Updated
1. Fixed static-check lint errors and warnings.
2. Fixed MacOS daemonize to log to StdOut (was logging to StdErr).
3. Fixed Windows event logging.
4. Fixed daemonize configuration overwrite issue.


## [0.8.3](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.3) - 2022-12-16

### Added
1. Added ARM64 to release build artifacts

### Changed
1. Updated _systemd_ unit file to wait on `network-online.target`.
2. Updated service lockfile to use `flock` _syscall_.
3. Removed _zip_ files from release artifacts (no longer necessary)


## [0.8.2](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.2) - 2022-10-14

### Changed
1. Reworked RecordSpecialEvents to not use wrapped requests/responses
2. Added 'swipe open' and 'swipe close' event reasons to message internationalisation.
3. Included health-check interval in watchdog configuration. 


## [0.8.1](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.1) - 2022-08-01

### Changed
1. Reworked event struct in `get-status`, `get-event` and `get-events` to include:
   - event type code and description
   - event reason code and description
   - event direction code and description
2. Added (optional) protocol version to configuration.
3. Added (optional) translation locale to configuration.
4. Resolved INADDR_ANY to interface IPv4 address for controller listener address health check.


## [0.8.0](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.0) - 2022-07-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.8.0

## [0.7.3](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.7.3) - 2022-06-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.7.3

### [0.7.2](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.7.2)

### Changed
1. Migrated to uhppoted-lib `config` command implementation
2. Reworked `get-events` to return `first`, `last` and `current` event indices.
3. Reworked `get-event` to return `first`, `last`, `current` and `next` events.

### [0.7.1](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.7.1)

### Changed
1. Task list support:
   -  `PutTaskList`
2. Migrated to IUHPPOTED interface + implementation

## Older

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.7.0    | Added support for time profiles from the extended API                                     |
| v0.6.12   | Added handling for `nil` events in `GetStatus`                                            |
| v0.6.10   | Maintenance release for version compatibility with `uhppoted-app-wild-apricot`            |
| v0.6.8    | Maintenance release for version compatibility with `uhppoted-core` `v0.6.8`               |
| v0.6.7    | Implements `special-events` API to enable/disable door events                             |
| v0.6.5    | Maintenance release for version compatibility with `node-red-contrib-uhppoted`            |
| v0.6.4    | Maintenance release for version compatibility with `uhppoted-app-sheets`                  |
| v0.6.3    | Maintenance release to update module dependencies                                         |
| v0.6.2    | Implements access control list API                                                        |
| v0.6.1    | Maintenance release to update module dependencies                                         |
| v0.6.0    | Maintenance release to update module dependencies                                         |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules* |

