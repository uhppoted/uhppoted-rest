## v0.7.x

### IN PROGRESS

- [ ] Make compatible with updated `uhppoted-api` and `uhppote-core`
- [ ] Implement `get-time-profile`
- [ ] Implement `set-time-profile`
- [ ] Implement `clear-time-profiles`
- [ ] Implement `get-time-profiles`
- [ ] Implement `set-time-profiles`
- [ ] Implement time profiles for ACL API
- [ ] Check against time profile for open-door API
- [ ] Replace internal healthcheck with implementation from `uhppoted-api`

## TODO

- [ ] [retool](https://retool.com)
- [ ] Apply API actions to multiple devices
- [ ] Redesign API around d
- [ ] Revisit POST/PUT semantics (https://restfulapi.net/rest-put-vs-post)
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
