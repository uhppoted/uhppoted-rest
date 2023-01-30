// Copyright 2023 uhppoted@twyst.co.za. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

/*
Package uhppoted-rest implements a REST server for the UHPPOTE TCP/IP Wiegand-26 access controllers.

The REST server wraps the low level UDP API implemented by uhppote-core in a REST layer, along with
some additional functionality to manage access control lists and events. The server include supports
for:

  - HTTP
  - HTTPS (optionally with client TLS authentication)
  - operation authorization
  - system health monitoring
*/
package rest
