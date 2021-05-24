module github.com/uhppoted/uhppoted-rest

go 1.16

require (
	github.com/uhppoted/uhppote-core v0.6.13-0.20210524184639-f1352385886e
	github.com/uhppoted/uhppoted-api v0.6.13-0.20210524193322-a4dd68c940e2
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887
)

replace github.com/uhppoted/uhppote-core => ../uhppote-core

replace github.com/uhppoted/uhppoted-api => ../uhppoted-api
