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

- [v0.6.1 tar.gz](https://github.com/uhppoted/uhppoted-rest/releases/download/v0.6.1/uhppoted-rest_v0.6.1.tar.gz)
- [v0.6.1 zip](https://github.com/uhppoted/uhppoted-rest/releases/download/v0.6.1/uhppoted-rest_v0.6.1.zip)

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

Usage: ```uhppoted-rest <command> <options>```

Supported commands:

- `help`
- `version`
- `run`
- `daemonize`
- `undaemonize`

Defaults to `run` if the command it not provided i.e. ```uhppoted-rest <options>``` is equivalent to ```uhppoted-rest run <options>```.

The OpenAPI specification for the [REST API](https://github.com/uhppoted/uhppoted-rest/blob/master/documentation/uhppoted-api.yaml) is included in the documentation folder.

### `run`

Runs the `uhppoted-rest` REST API server. Intended for use as a system service that runs in the background to handle REST requests. 

Command line:

` uhppoted-rest [--debug] [--console] [--config <file>] `

```
  --config      Sets the uhppoted.conf file to use for controller configurations. 
                Defaults to the communal uhppoted.conf file shared by all the uhppoted 
                modules.
  --console     Runs the REST API server as an application, logging events to the
                console.
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers
```

### `daemonize`

Registers the `uhppoted-rest` REST API server as a system service that will be started on system boot. The command creates the necessary system specific service configuration files and service manager entries.

Command line:

`uhppoted-rest daemonize `

### `undaemonize`

Unregisters the `uhppoted-rest` REST API server as a system service, but does not delete any created log or configuration files. 

Command line:

`uhppoted-rest undaemonize `


