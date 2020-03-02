module github.com/uhppoted/uhppoted-rest

go 1.14

require (
	github.com/uhppoted/uhppote-core v0.0.0-20200228192138-00c62a4d6ea3
	github.com/uhppoted/uhppoted-api v0.0.0-20200302181311-56c5fea77afc
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527
)

replace (
	github.com/uhppoted/uhppote-core => ../uhppote-core
	github.com/uhppoted/uhppoted-api => ../uhppoted-api
)
