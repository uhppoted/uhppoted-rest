module github.com/uhppoted/uhppoted-rest

go 1.16

require (
	github.com/uhppoted/uhppote-core v0.6.13-0.20210519162147-c4a6d5b3fe33
	github.com/uhppoted/uhppoted-api v0.6.13-0.20210519162421-ae61753cbb9e
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887
)

replace github.com/uhppoted/uhppoted-api => ../uhppoted-api
