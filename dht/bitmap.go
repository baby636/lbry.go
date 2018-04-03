package dht

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"github.com/lbryio/errors.go"
	"github.com/lyoshenka/bencode"
)

type bitmap [nodeIDLength]byte

func (b bitmap) RawString() string {
	return string(b[:])
}

func (b bitmap) Hex() string {
	return hex.EncodeToString(b[:])
}

func (b bitmap) HexShort() string {
	return hex.EncodeToString(b[:4])
}

func (b bitmap) Equals(other bitmap) bool {
	for k := range b {
		if b[k] != other[k] {
			return false
		}
	}
	return true
}

func (b bitmap) Less(other interface{}) bool {
	for k := range b {
		if b[k] != other.(bitmap)[k] {
			return b[k] < other.(bitmap)[k]
		}
	}
	return false
}

func (b bitmap) Xor(other bitmap) bitmap {
	var ret bitmap
	for k := range b {
		ret[k] = b[k] ^ other[k]
	}
	return ret
}

// PrefixLen returns the number of leading 0 bits
func (b bitmap) PrefixLen() int {
	for i := range b {
		for j := 0; j < 8; j++ {
			if (b[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}
	return numBuckets
}

func (b *bitmap) UnmarshalBencode(encoded []byte) error {
	var str string
	err := bencode.DecodeBytes(encoded, &str)
	if err != nil {
		return err
	}
	if len(str) != nodeIDLength {
		return errors.Err("invalid node ID length")
	}
	copy(b[:], str)
	return nil
}

func (b bitmap) MarshalBencode() ([]byte, error) {
	str := string(b[:])
	return bencode.EncodeBytes(str)
}

func newBitmapFromBytes(data []byte) bitmap {
	if len(data) != nodeIDLength {
		panic("invalid bitmap of length " + strconv.Itoa(len(data)))
	}

	var bmp bitmap
	copy(bmp[:], data)
	return bmp
}

func newBitmapFromString(data string) bitmap {
	return newBitmapFromBytes([]byte(data))
}

func newBitmapFromHex(hexStr string) bitmap {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return newBitmapFromBytes(decoded)
}

func newRandomBitmap() bitmap {
	var id bitmap
	_, err := rand.Read(id[:])
	if err != nil {
		panic(err)
	}
	return id
}
