## v0.6.x

** IN PROGRESS **

- [x] PUT card
- [x] DELETE card
- [x] grant
- [x] revoke
- [x] show
- [x] load-acl
- [ ] get-acl
- [ ] Update OpenAPI 'show' returned object to be array
- [ ] compare-acl
- [ ] rethink dates on grant/revoke ALL
- [ ] rethink dates on put/grant if current record has no permissions
- [ ] rework to use uhppoted-api::Config
- [ ] Move current 'REST' functions to 'device' package
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
