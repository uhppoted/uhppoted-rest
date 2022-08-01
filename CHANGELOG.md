# CHANGELOG

## [v0.8.1] - 2022-08-01

### Changed
1. Reworked event struct in `get-status`, `get-event` and `get-events` to include:
   - event type code and description
   - event reason code and description
   - event direction code and description
2. Added (optional) protocol version to configuration.
3. Added (optional) translation locale to configuration.
4. Resolved INADDR_ANY to interface IPv4 address for controller listener address health check.


## [v0.8.0] - 2022-07-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.8.0

## [v0.7.3] - 2022-06-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.7.3

### [v0.7.2]

### Changed
1. Migrated to uhppoted-lib `config` command implementation
2. Reworked `get-events` to return `first`, `last` and `current` event indices.
3. Reworked `get-event` to return `first`, `last`, `current` and `next` events.

### [v0.7.1]

### Changed
1. Task list support:
   -  `PutTaskList`
2. Migrated to IUHPPOTED interface + implementation
