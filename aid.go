package main

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"
)

var AidTimer int64
var AidCounter uint32 = 0 // 后续限制最大生成3字节！

func Aid() []byte {
	now := 62135596800 + time.Now().Unix()
	if AidTimer == now {
		AidCounter++
	} else {
		AidTimer = now
		AidCounter = 0
	}
	res := [8]byte{
		byte(0xff & (AidTimer >> 32)),
		byte(0xff & (AidTimer >> 24)),
		byte(0xff & (AidTimer >> 16)),
		byte(0xff & (AidTimer >> 8)),
		byte(0xff & AidTimer),
		byte(0xff & (AidCounter >> 16)),
		byte(0xff & (AidCounter >> 8)),
		byte(0xff & AidCounter),
	}
	return res[:]
}

func Aid10(aid []byte) uint64 {
	return binary.BigEndian.Uint64(aid)
}

func Aid36(aid []byte) string {
	return strconv.FormatUint(Aid10(aid), 36)
}

func Aid16(aid []byte) string {
	return hex.EncodeToString(aid)
}

func AidDecoder(aid []byte) (uint64, uint32) {
	timer := binary.BigEndian.Uint64(append(make([]byte, 3), aid[0:5]...))
	counter := binary.BigEndian.Uint32(append(make([]byte, 1), aid[5:]...))
	return timer, counter
}

func Aid10Decoder(aid10 uint64) (uint64, uint32) {
	aid := make([]byte, 8)
	binary.BigEndian.PutUint64(aid, aid10)
	return AidDecoder(aid)
}

func Aid36Decoder(aid36 string) (uint64, uint32) {
	aid10, err := strconv.ParseUint(aid36, 36, 64)
	if err != nil {
		return 0, 0
	}
	return Aid10Decoder(aid10)
}

func Aid16Decoder(aid16 string) (uint64, uint32) {
	aid, err := hex.DecodeString(aid16)
	if err != nil {
		return 0, 0
	}
	return AidDecoder(aid)
}
