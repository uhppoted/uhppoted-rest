module github.com/uhppoted/uhppoted-rest

go 1.14

require (
	github.com/uhppoted/uhppote-core v0.6.1
	github.com/uhppoted/uhppoted-api v0.6.1
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527
)

replace (
	github.com/uhppoted/uhppoted-api => ../uhppoted-api
)