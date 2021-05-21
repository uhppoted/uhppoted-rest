module github.com/uhppoted/uhppoted-rest

go 1.16

require (
	github.com/uhppoted/uhppote-core v0.6.13-0.20210521170201-a1e7a69be646
	github.com/uhppoted/uhppoted-api v0.6.13-0.20210521171353-5f3529a78611
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887
)

replace github.com/uhppoted/uhppote-core => ../uhppote-core

replace github.com/uhppoted/uhppoted-api => ../uhppoted-api
