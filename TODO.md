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
- [ ] Migrate current 'REST' functions to uhppoted-api
- [ ] build documentation
- [ ] install documentation

## TODO

### uhppoted-rest
- [ ] Get events after XXX
- [ ] Client certificate revocation list
- [ ] uhppoted-rest: get-events date/id range
- [ ] commonalise functionality with uhppoted-mqttd

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
