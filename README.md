![build](https://github.com/uhppoted/uhppoted-rest/workflows/build/badge.svg)
![build](https://github.com/uhppoted/uhppoted-rest/workflows/ghcr/badge.svg)

# uhppoted-rest

Implements a REST API for use with access control systems based on the *UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards.

Supported platforms:
- Linux
- MacOS
- Windows
- RaspberryPi (ARM/ARM7/ARM6)

## Raison d'être

`uhppoted-rest` implements a REST API that facilitates integration of the access control function with other systems (e.g. web servers, mobile applications) without requiring the device level functionality being built-in to the application.

---
### Contents

- [Release Notes](#release-notes)
- [Installation](#installation)
   - [Docker](#docker)
   - [Building from source](#building-from-source)
- [Command line](#uhppoted-rest)
- [Configuration](#configuration)
   - [Authorization](#api-authorization)
       - [_users_ file](#users-file)
       - [_groups_ file](#groups-file)
       - ['Basic' authentication](#basic-authentication)
       - [HOTP](#hotp-authentication)
       - [`open-door`](#open-door-api)
---

## Release Notes

### Current Release

**[v0.8.10](https://github.com/uhppoted/uhppoted-rest/releases/tag/v0.8.10) - 2024-01-30**

1. ARMv6 build target (RaspberryPi ZeroW).
2. Added support for _auto-send_ interval.

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppoted-rest/releases):

The release tarballs contain the executables for all the operating systems - OS specific tarballs with all the _uhppoted_ components can be found in [uhpppoted](https://github.com/uhppoted/uhppoted/releases) releases.

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


### Docker

A public _Docker_ image is published to [ghcr.io](https://github.com/uhppoted?tab=packages&repo_name=uhppoted-rest). 

The image is configured to use the `/usr/local/etc/uhppoted/uhppoted.conf` file for configuration information.

#### `docker compose`

A sample Docker `compose` configuration is provided in the [`docker/compose`](docker) folder. 

To run the example, download and extract the [compose.zip](docker) scripts and supporting files into folder
of your choice and then:
```
cd <compose folder>
docker compose up
```

The REST server can be tested using _curl_, e.g.:
```
curl -X 'GET' 'http://127.0.0.1:8080/uhppote/device' -H 'accept: application/json' | jq .
```

The default image is configured for HTTP only. To enable HTTPS, enable it in the _uhppoted.conf_ file on the Docker volume 
mapped to _/usr/local/etc/uhppoted_ and copy the certificates and keys to the _/usr/local/etc/uhppoted/rest_ folder on 
the Docker volume mapped to _/usr/local/etc/uhppoted_, e.g.
```
docker cp ca.cert       restd:/usr/local/etc/uhppoted/rest
docker cp uhppoted.key  restd:/usr/local/etc/uhppoted/rest
docker cp uhppoted.cert restd:/usr/local/etc/uhppoted/rest
```

#### `docker run`

To start a REST server using Docker `run`:
```
docker pull ghcr.io/uhppoted/restd:latest
docker run --publish 8080:8080 --publish 8443:8443 --name restd --mount source=uhppoted,target=/var/uhppoted --rm ghcr.io/uhppoted/restd
```

The REST server can be tested using _curl_, e.g.:
```
curl -X 'GET' 'http://127.0.0.1:8080/uhppote/device' -H 'accept: application/json' | jq .
```

#### `docker build`

For inclusion in a Dockerfile:
```
FROM ghcr.io/uhppoted/restd:latest
```


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
go build -trimpath -o bin ./...
```

The above commands build the `'uhppoted-rest` executable to the `bin` directory.

#### Dependencies

| *Dependency*                                             | *Description*                              |
| -------------------------------------------------------- | ------------------------------------------ |
| [uhppote-core](https://github.com/uhppoted/uhppote-core) | Device level API implementation            |
| [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) | common API for external                    |
| golang.org/x/sys                                         | Support for Windows services               |

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

## Configuration

### API Authorization

The original design for the REST API expected requests on a secured connection from a gateway service that would manage authentication and authorization (and in general manage security for systems exposed to the wider Internet). The design allowed for mutual TLS authentication where the gateway and REST server were not contained within a secured/private network.

However, it is generally simpler and more convenient for mobile client applications and IoT systems to be able to access the REST API directly, which adds a requirement for per API authentication. 

For the REST server, per API authentication is (optionally) implemented using the _Authorization_ request header. Currently, two methods are supported:

  - [RFC 7617](https://tools.ietf.org/html/rfc7617) 'Basic' authentication
  - Custom 'Bearer' authorization with HOTP tokens 

with JWT and API authorization tokens are earmarked for implementation at a later stage. 

Permissions for API calls are managed by the _groups_ and _users_ files:

  1. A _user_, comprising a user ID, password and/or HOTP key and card number(s), is assigned to a group (or groups)
  2. A _group_ comprises a set of permissions that allow access to an API

API authorization is disabled by default and can be enabled by updating the _uhppoted_ configuration:

    rest.auth.enabled = true
    # rest.auth.users = <'users' file>
    # rest.auth.groups = <'groups' file>
    # rest.auth.hotp.range = <window for HOTP counters>
    # rest.auth.hotp.counters = <HOTP counters file>
    
The _users_, _groups_, and _hotp.counters_ files default, respectively, to:

    - <workdir>/rest/users
    - <workdir>/rest/groups
    - <workdir>/rest/counters

The default HOTP counter window is 8 (i.e. a request authorized using HOTP can present an OTP based on a counter value at most 8 larger than the current stored counter for the user ID).

#### _users_ file

The _users_ file comprises a set of `name-value` pairs, each on single line with the name/value separated by two or more spaces. The value is a JSON encoded _user_ object, comprising:

| _field_  | _type_                      | _description_                         |
| -------- | --------------------------- | ------------------------------------- |
| password | hexadecimal string          | SHA256 encoded password, optional     |
| hotp     | hexadecimal string          | HOTP secret, optional                 |
| cards    | list of regular expressions | List regular expressions matching cards for which user is authorized.        |
| groups   | list of strings             | List groups to which user is assigned |

Sample _users_ file:

```
gateway    { "password":"4ea5ee68fea05586106890ded5733820bb77d919cda27bc4b8139b7cd33b8889", "groups": [ "system" ], "cards": [ ".*" ] }
qwerty     { "password":"7bba6743c0ddb67462771e1f74950bf9863f24b8b73087cf88b8b9b47917649c", "hotp": "DFIOJ3BJPHPCRJBT", "groups": [ "users" ], "cards": [ "1928374646" ] }
```

#### _groups_ file

The _groups_ file comprises a set of `name-value` pairs, each on single line with the name/value separated by two or more spaces. The value is a comma separated set of regular expressions that define the permissions for the group.

A permission is formatted as `resource:action`, where:

- _resource_ is a regular expression matching the request URL(s) to which the group has authorization
- _action_ is the HTTP method for the URL e.g. some groups may have GET permission but not POST permission

Sample _groups_ file:

```
system    *:*
users     /uhppote/device/[0-9]+/door/[0-9]/swipes:post
```

#### 'Basic' Authentication

'Basic' authentication ([RFC 7617](https://tools.ietf.org/html/rfc7617)) provides only the most elementary security and should only be used on trusted networks and/or on connections secured with TLS.

The _Authorization_ request header is formatted as `Authorization: Basic user:password`, where `user:password` is encoded as a Base64 string.

e.g. for user "_qwerty_" and password "_uiop_", the _Authorization_ request header is:

    Authorization: Basic cXdlcnR5OnVpb3A=

#### HOTP Authentication

'HOTP' authentication ([RFC 4226](https://tools.ietf.org/html/rfc4226)) is provided as a lightweight authentication mechanism for IoT systems that provides somewhat better security than 'Basic' authentication.

The _Authorization_ request header is formatted as `Authorization: Bearer HOTP:user:OTP`, where `HOTP:user:OTP` is encoded as a Base64 encoded string.

e.g. for user "_qwerty_" and OTP "_763927_", the _Authorization_ request header is:

    Authorization: Bearer SE9UUDpxd2VydHk6NzYzOTI3

HOTP implements a counter based OTP as implemented by e.g. Google Authenticator. 

#### `open-door` API

The `open-door` API implements two additional security checks:

- the presented card number is verified against the list of cards for the user in the _Authentication_ header
- the presented card number must have access rights for the requested door
