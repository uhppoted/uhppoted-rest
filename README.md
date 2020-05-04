# uhppoted-rest

Implements a REST API for use with access control systems based on the *UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards.

Supported operating systems:
- Linux
- MacOS
- Windows
- ARM7 _(e.g. RaspberryPi)_

## Raison d'être

`uhppoted-rest` implements a REST API that facilitates integration of the access control function with other systems (e.g. web servers, mobile applications) without requiring the device level functionality being built-in to the application.

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.6.1    | Maintenance release to update module dependencies                                         |
| v0.6.0    | Maintenance release to update module dependencies                                         |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules* |

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppoted-rest/releases):

- [v0.6.2 tar.gz](https://github.com/uhppoted/uhppoted-rest/releases/download/v0.6.2/uhppoted-rest_v0.6.2.tar.gz)
- [v0.6.2 zip](https://github.com/uhppoted/uhppoted-rest/releases/download/v0.6.2/uhppoted-rest_v0.6.2.zip)

The above archives contain the executables for all the operating systems - OS specific tarballs with all the _uhppoted_ components can be found in [uhpppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a directory of your choice. To install `uhppoted-rest` as a system service:
```
   cd <uhppote directory>
   sudo uhppoted-rest daemonize
```

`uhppoted-rest help` will list the available commands and associated options (documented below).

The `uhppoted-rest` service requires the following additional files:

- `uhppoted.conf`

### `uhppoted.conf`

`uhppoted.conf` is the communal configuration file shared by all the `uhppoted` project modules and is (or will 
eventually be) documented in [uhppoted](https://github.com/uhppoted/uhppoted). `uhppoted-rest` requires:
- the _REST_ section to define the configuration for the REST _httpd_ server
- the _devices_ section to resolve non-local controller IP addresses and door to controller door identities.

A sample [uhppoted.conf](https://github.com/uhppoted/uhppoted/blob/master/runtime/simulation/405419896.conf) file is included in the `uhppoted` distribution.

### Building from source

Assuming you have `Go` and `make` installed:

```
git clone https://github.com/uhppoted/uhppoted-rest.git
cd uhppoted-rest
make build
```

If you prefer not to use `make`:
```
git clone https://github.com/uhppoted/uhppoted-rest.git
cd uhppoted-rest
mkdir bin
go build -o bin ./...
```

The above commands build the `'uhppoted-rest` executable to the `bin` directory.

#### Dependencies

| *Dependency*                          | *Description*                                          |
| ------------------------------------- | ------------------------------------------------------ |
| [com.github/uhppoted/uhppote-core](https://github.com/uhppoted/uhppote-core) | Device level API implementation            |
| [com.github/uhppoted/uhppoted-api](https://github.com/uhppoted/uhppoted-api) | common API for external         |
| golang.org/x/sys/windows              | Support for Windows services                           |
| golang.org/x/lint/golint              | Additional *lint* check for release builds             |

## uhppoted-rest

Usage: *uhppoted-rest \<command\> \<options\>*

Defaults to `run` unless one of the commands below is specified: 

- `daemonize`
- `undaemonize`
- `help`
- `version`

Supported `run` options:
- `--console`
- `--debug`
