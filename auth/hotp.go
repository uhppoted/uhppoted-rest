package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/uhppoted/uhppoted-lib/kvs"
	"hash"
	"log"
	"math"
	"strconv"
	"strings"
)

type HOTP struct {
	window   uint64
	counters struct {
		*kvs.KeyValueStore
		filepath string
		log      *log.Logger
	}
}

const DIGITS = 6

func NewHOTP(window uint64, counters string, logger *log.Logger) (*HOTP, error) {
	v := func(value string) (interface{}, error) {
		return strconv.ParseUint(value, 10, 64)
	}

	hotp := HOTP{
		window: window,
		counters: struct {
			*kvs.KeyValueStore
			filepath string
			log      *log.Logger
		}{
			kvs.NewKeyValueStore("hotp:counters", v),
			counters,
			logger,
		},
	}

	if err := hotp.counters.LoadFromFile(counters); err != nil {
		log.Printf("WARN  %v", err)
	}

	return &hotp, nil
}

func (hotp *HOTP) Validate(uid, secret, otp string) error {
	otp = strings.TrimSpace(otp)
	if len(otp) != DIGITS {
		return fmt.Errorf("%s: Invalid OTP '%s' - expected %d digits", uid, otp, DIGITS)
	}

	counter, ok := hotp.counters.Get(uid)
	if !ok {
		counter = uint64(1)
	}

	for i := uint64(0); i < hotp.window; i++ {
		generated, err := generateHOTP(secret, counter.(uint64), DIGITS, sha1.New)
		if err != nil {
			return err
		}

		if subtle.ConstantTimeCompare([]byte(generated), []byte(otp)) == 1 {
			hotp.counters.Store(uid, counter.(uint64)+1, hotp.counters.filepath, hotp.counters.log)
			return nil
		}

		counter = counter.(uint64) + 1
	}

	return fmt.Errorf("%s: invalid OTP %s", uid, otp)
}

// Ref. https://github.com/pquerna/otp
func generateHOTP(secret string, counter uint64, digits int, algorithm func() hash.Hash) (passcode string, err error) {
	secret = strings.TrimSpace(secret)
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	bytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 8)
	mac := hmac.New(algorithm, bytes)
	binary.BigEndian.PutUint64(buffer, counter)

	mac.Write(buffer)
	sum := mac.Sum(nil)

	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	mod := int32(value % int64(math.Pow10(digits)))

	return fmt.Sprintf("%06d", mod), nil
}
