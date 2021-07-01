## v0.7.1

- [x] Migrate ACL dispatch to use IUHPPOTED
- [x] `SetTaskList`

## TODO

- [ ] Clean up REST error handling (complete 'mare, what was I thinking :-()
- [ ] Replace internal healthcheck with implementation from `uhppoted-lib`
- [ ] OpenAPI: break up spec into sections
- [ ] OpenAPI: use fixed map keys for doors, segments
- [ ] OpenAPI: fix all the new semantic/syntax errors

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
