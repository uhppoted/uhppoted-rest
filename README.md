# uhppoted-rest

Wraps the `uhppote-core` device API in a REST API for use with access control systems based on the 
*UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards.

Supported operating systems:
- Linux
- MacOS
- Windows

## Raison d'être

`uhppoted-rest` implements a REST API that facilitates integration of the access control function with other 
systems (e.g. web servers, mobile applications) without requiring the device level functionality being built-in 
to the application.

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.6.1    | Maintenance release to update module dependencies                                         |
| v0.6.0    | Maintenance release to update module dependencies                                         |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules* |

## Installation

### Building from source

#### Dependencies

| *Dependency*                          | *Description*                                          |
| ------------------------------------- | ------------------------------------------------------ |
| [com.github/uhppoted/uhppote-core][1] | Device level API implementation                        |
| [com.github/uhppoted/uhppoted-api][2] | common API for external applications                   |
| golang.org/x/sys/windows              | Support for Windows services                           |
| golang.org/x/lint/golint              | Additional *lint* check for release builds             |

### Binaries

## uhppoted-rest

Usage: *uhppoted-rest \<command\> \<options\>*

Defaults to 'run' unless one of the commands below is specified: 

- daemonize
- undaemonize
- help
- version

Supported 'run' options:
- --console
- --debug

[1]: https://github.com/uhppoted/uhppote-core
[2]: https://github.com/uhppoted/uhppoted-api

