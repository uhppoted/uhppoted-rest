# TODO

- [x] Migrate to codeberg.org (cf. https://codeberg.org/uhppoted/uhppoted/issues/1)
- [x] 'set-firstcard' API (cf. https://codeberg.org/uhppoted/uhppoted/issues/2)

- [ ] add first-card privilege to `get-card` and `put-card` (cf. https://codeberg.org/uhppoted/uhppoted/issues/3)
    - [x] put-card
    - [x] get-card
    - [x] ~get-card-by-index~
    - [ ] update ACL functions
        - [x] config::withPIN
        - [x] config::withFirstCard
        - [x] getACL
        - [x] putACL
        - [ ] grant
        - [ ] revoke
        - [ ] show
        - [ ] compare

- [ ] clean up OpenAPI

### To Be Done

- (?) [Kiota](https://learn.microsoft.com/en-us/openapi/kiota/overview)
- [ ] Clean up REST error handling (complete 'mare, what was I thinking :-()
- [ ] Replace internal healthcheck with implementation from `uhppoted-lib`
- [ ] OpenAPI: use fixed map keys for doors, segments
- [ ] OpenAPI: fix all the new semantic/syntax errors
- [ ] GetEvent - range

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
