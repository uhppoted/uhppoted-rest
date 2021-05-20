module github.com/uhppoted/uhppoted-rest

go 1.16

require (
	github.com/uhppoted/uhppote-core v0.6.13-0.20210520192929-2298aeaedba0
	github.com/uhppoted/uhppoted-api v0.6.13-0.20210520193707-b8f6e75a5502
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887
)

replace github.com/uhppoted/uhppote-core => ../uhppote-core

replace github.com/uhppoted/uhppoted-api => ../uhppoted-api
