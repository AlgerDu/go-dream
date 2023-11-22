package dinfra

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strconv"
)

type (
	ID int64
)

var (
	decodeBase32Map [256]byte
	decodeBase58Map [256]byte

	ErrInvalidBase58 = errors.New("invalid base58")
	ErrInvalidBase32 = errors.New("invalid base32")
)

const (
	encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"
	encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

func init() {

	for i := 0; i < len(decodeBase58Map); i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}

	for i := 0; i < len(decodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

func (f ID) Int64() int64 {
	return int64(f)
}

func ParseInt64(id int64) ID {
	return ID(id)
}

func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

func ParseString(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return ID(i), err
}

func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

func ParseBase2(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 2, 64)
	return ID(i), err
}

func (f ID) Base32() string {

	if f < 32 {
		return string(encodeBase32Map[f])
	}

	b := make([]byte, 0, 12)
	for f >= 32 {
		b = append(b, encodeBase32Map[f%32])
		f /= 32
	}
	b = append(b, encodeBase32Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

func ParseBase32(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase32Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase32
		}
		id = id*32 + int64(decodeBase32Map[b[i]])
	}

	return ID(id), nil
}

func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

func ParseBase36(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 36, 64)
	return ID(i), err
}

func (f ID) Base58() string {

	if f < 58 {
		return string(encodeBase58Map[f])
	}

	b := make([]byte, 0, 11)
	for f >= 58 {
		b = append(b, encodeBase58Map[f%58])
		f /= 58
	}
	b = append(b, encodeBase58Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

func ParseBase58(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase58Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase58
		}
		id = id*58 + int64(decodeBase58Map[b[i]])
	}

	return ID(id), nil
}

func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

func ParseBase64(id string) (ID, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return -1, err
	}
	return ParseBytes(b)
}

func (f ID) Bytes() []byte {
	return []byte(f.String())
}

func ParseBytes(id []byte) (ID, error) {
	i, err := strconv.ParseInt(string(id), 10, 64)
	return ID(i), err
}

func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

func ParseIntBytes(id [8]byte) ID {
	return ID(int64(binary.BigEndian.Uint64(id[:])))
}
