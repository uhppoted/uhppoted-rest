# CHANGELOG

## Unreleased

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
