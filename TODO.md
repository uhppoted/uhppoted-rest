## v0.6.x

** IN PROGRESS **

- [x] PUT card
- [x] DELETE card
- [x] grant
- [x] revoke
- [x] show
- [x] load-acl
- [x] get-acl
- [x] Update OpenAPI 'show' returned object to be array
- [x] Rethink dates on grant if current record has no permissions
- [x] Rethink dates on grant ALL
- [x] Rework to use uhppoted-api::Config
- [x] Move current 'REST' functions to 'device' package
- [x] build documentation
- [x] install documentation
- [x] Migrate current 'REST' functions to uhppoted-api
- [x] Use uhppoted-api event logging
- [x] Future proof API
- [ ] Command documentation
- [ ] Add from/to/doors validity check for parsing JSON cards
- [ ] Return original column name in uhppoted-api::parseHeader error

## TODO

- [ ] Trace/log requests/responses
- [ ] Request ID's
- [ ] Skeleton integration test
- [ ] Get events after XXX
- [ ] Client certificate revocation list
- [ ] uhppoted-rest: get-events date/id range

### Documentation

- [ ] godoc
- [ ] user manuals
- [ ] man/info page

### Other

1.  Integration tests
2.  Verify fields in listen events/status replies against SDK:
    - battery status can be (at least) 0x00, 0x01 and 0x04
3.  EventLogger 
    - MacOS: use [system logging](https://developer.apple.com/documentation/os/logging)
    - Windows: event logging
