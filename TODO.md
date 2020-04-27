## v0.6.x

** IN PROGRESS **

- [x] grant
- [x] revoke
- [x] show
- [ ] rethink dates on grant/revoke ALL
- [ ] load-acl
- [ ] get-acl
- [ ] compare-acl
- [x] PUT card
- [x] DELETE card
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
