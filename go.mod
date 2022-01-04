module github.com/uhppoted/uhppoted-rest

go 1.16

require (
	github.com/uhppoted/uhppote-core v0.7.2-0.20211231212401-366db0b80d0c
	github.com/uhppoted/uhppoted-lib v0.7.2-0.20220103212354-4beb6d44504b
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887
)

replace github.com/uhppoted/uhppoted-lib => ../uhppoted-lib
