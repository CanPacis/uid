package uid

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
)

var encoding = base64.NewEncoding("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@*")

var ErrInvalidUIDFormat = errors.New("invalid uid format")

type UID [12]byte

func (uid UID) String() string {
	return encoding.EncodeToString(uid[:])
}

func (uid *UID) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return uid.UnmarshalBinary(src)
	case string:
		var err error
		*uid, err = Parse(src)
		return err
	default:
		return fmt.Errorf("uid scan: unsupported type %T", src)
	}
}

func (uid UID) Value() (driver.Value, error) {
	return uid.String(), nil
}

func (uid UID) MarshalBinary() ([]byte, error) {
	return uid[:], nil
}

func (uid *UID) UnmarshalBinary(src []byte) error {
	if len(src) != 12 {
		return fmt.Errorf("uid parse: %w (length %d)", ErrInvalidUIDFormat, len(src))
	}
	copy(uid[:], src)
	return nil
}

func (uid UID) MarshalText() ([]byte, error) {
	return []byte(uid.String()), nil
}

func (uid *UID) UnmarshalText(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("uid parse: %w (length %d)", ErrInvalidUIDFormat, len(data))
	}

	b, err := encoding.DecodeString(string(data))
	if err != nil {
		return fmt.Errorf("uid parse: %w", err)
	}
	copy(uid[:], b)
	return nil
}

func New() UID {
	uid := UID{}
	rand.Read(uid[:])
	return uid
}

func Parse(s string) (UID, error) {
	uid := UID{}
	return uid, uid.UnmarshalText([]byte(s))
}

func MustParse(s string) UID {
	uid, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return uid
}

func Validate(s string) error {
	if len([]byte(s)) != 16 {
		return fmt.Errorf("%w: invalid length", ErrInvalidUIDFormat)
	}

	_, err := encoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidUIDFormat, err)
	}
	return nil
}
