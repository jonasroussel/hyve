package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	"github.com/google/uuid"
)

func GenerateSessionID() (uint, error) {
	bigID, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return 0, err
	}

	return uint(bigID.Uint64()), nil
}

func EncodeSessionID(id uint) string {
	idSize := binary.Size(id)

	idBytes := make([]byte, idSize)
	if idSize == 8 {
		binary.BigEndian.PutUint64(idBytes, uint64(id))
	} else {
		binary.BigEndian.PutUint32(idBytes, uint32(id))
	}

	idB64 := make([]byte, base64.StdEncoding.EncodedLen(len(idBytes)))
	base64.StdEncoding.Encode(idB64, idBytes)

	return string(idB64)
}

func DecodeSessionID(idB64 string) (uint, error) {
	idBytes, err := base64.StdEncoding.DecodeString(idB64)
	if err != nil {
		return 0, err
	}

	var id uint
	if len(idBytes) == 8 {
		id = uint(binary.BigEndian.Uint64(idBytes))
	} else {
		id = uint(binary.BigEndian.Uint32(idBytes))
	}

	return id, nil
}

func main() {
	id, err := GenerateSessionID()
	if err != nil {
		panic(err)
	}

	fmt.Println(uuid.New().Version())
	fmt.Println(id)

	idB64 := EncodeSessionID(id)
	fmt.Println(idB64)

	idDecoded, err := DecodeSessionID(idB64)
	if err != nil {
		panic(err)
	}

	fmt.Println(idDecoded)
}
