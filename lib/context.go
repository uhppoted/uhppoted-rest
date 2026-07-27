package lib

type ContextKey int

const (
	Uhppote ContextKey = iota + 1
	DeviceID
	Devices
	CardNumber
	Door
	AuthorizedCards
	Compression
	WithPINs
	WithFirstCard
)

func (k ContextKey) String() string {
	return [...]string{
		"uhppote",
		"device-id",
		"devices",
		"card-number",
		"door",
		"authorized-cards",
		"compression",
		"with-PINs",
		"with-firstcard",
	}[k]
}
