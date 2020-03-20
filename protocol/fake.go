package protocol

import "crypto/rand"

var is_test = false
var fakeRand []byte
var fakeSessionID uint32
var initializationVector []byte

func readRand(bs []byte) {
	if is_test {
		copy(bs, fakeRand)
	} else {
		rand.Read(bs)
	}
}

func SetFake(rand []byte, sessionID uint32) {
	fakeRand = make([]byte, len(rand))
	copy(fakeRand, rand)
	fakeSessionID = sessionID

	is_test = true
}

func SetInitializationVector(iv []byte) {
	initializationVector = iv
}

func makeInitializationVector(bs []byte) []byte {
	if is_test {
		if len(initializationVector) == len(bs) {
			copy(bs, initializationVector)
			return bs
		}
	}

	rand.Read(bs)
	return bs
}
