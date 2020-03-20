package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"strconv"
	"strings"
)

type RAKPAlgorithmAuthMethod uint8

const (
	RAKPAlgorithmAuth_None        RAKPAlgorithmAuthMethod = 0
	RAKPAlgorithmAuth_HMAC_SHA1   RAKPAlgorithmAuthMethod = 1
	RAKPAlgorithmAuth_HMAC_MD5    RAKPAlgorithmAuthMethod = 2
	RAKPAlgorithmAuth_HMAC_SHA256 RAKPAlgorithmAuthMethod = 3
)

func ParseAuthMethod(s string) (RAKPAlgorithmAuthMethod, error) {
	switch strings.ToLower(s) {
	case "", "none":
		return RAKPAlgorithmAuth_None, nil
	case "hmac-sha1":
		return RAKPAlgorithmAuth_HMAC_SHA1, nil
	case "hmac-md5":
		return RAKPAlgorithmAuth_HMAC_MD5, nil
	case "hmac-sha256":
		return RAKPAlgorithmAuth_HMAC_SHA256, nil
	default:
		return RAKPAlgorithmAuth_None, errors.New(s + " is unsupported confidentialityMethod")
	}
}

func (m RAKPAlgorithmAuthMethod) String() string {
	switch m {
	case RAKPAlgorithmAuth_None:
		return "none"
	case RAKPAlgorithmAuth_HMAC_SHA1:
		return "hmac-sha1"
	case RAKPAlgorithmAuth_HMAC_MD5:
		return "hmac-md5"
	case RAKPAlgorithmAuth_HMAC_SHA256:
		return "hmac-sha256"
	default:
		return "unknow(" + strconv.FormatInt(int64(m), 10) + ")"
	}
}

var RAKPAlgorithmAuthLengths = []int{
	0, // RAKPAlgorithmAuth_None
	sha1.Size,
	md5.Size,
	sha256.Size}

var RAKPAlgorithmAuthMethods = []func(key, data []byte) []byte{
	// RAKPAlgorithmAuth_None
	func(key, data []byte) []byte {
		return nil
	},

	// RAKPAlgorithmAuth_HMAC_SHA1
	func(key, data []byte) []byte {
		mac := hmac.New(sha1.New, key)
		mac.Write(data)
		return mac.Sum(nil)
	},

	// RAKPAlgorithmAuth_HMAC_MD5
	func(key, data []byte) []byte {
		mac := hmac.New(md5.New, key)
		mac.Write(data)
		return mac.Sum(nil)
	},

	// RAKPAlgorithmAuth_HMAC_SHA256
	func(key, data []byte) []byte {
		mac := hmac.New(sha256.New, key)
		mac.Write(data)
		return mac.Sum(nil)
	},
}

func (self RAKPAlgorithmAuthMethod) Gen(key, data []byte) []byte {
	return RAKPAlgorithmAuthMethods[self](key, data)
}

func (self RAKPAlgorithmAuthMethod) Size() int {
	return RAKPAlgorithmAuthLengths[self]
}

type RAKPAlgorithmIntegrityMethod uint8

const (
	RAKPAlgorithmIntegrity_None            RAKPAlgorithmIntegrityMethod = 0
	RAKPAlgorithmIntegrity_HMAC_SHA1_96    RAKPAlgorithmIntegrityMethod = 1
	RAKPAlgorithmIntegrity_HMAC_MD5_128    RAKPAlgorithmIntegrityMethod = 2
	RAKPAlgorithmIntegrity_MD5_128         RAKPAlgorithmIntegrityMethod = 3
	RAKPAlgorithmIntegrity_HMAC_SHA256_128 RAKPAlgorithmIntegrityMethod = 4
)

func ParseIntegrityMethod(s string) (RAKPAlgorithmIntegrityMethod, error) {
	switch strings.ToLower(s) {
	case "", "none":
		return RAKPAlgorithmIntegrity_None, nil
	case "hmac-sha1-96":
		return RAKPAlgorithmIntegrity_HMAC_SHA1_96, nil
	case "hmac-md5-128":
		return RAKPAlgorithmIntegrity_HMAC_MD5_128, nil
	case "md5-128":
		return RAKPAlgorithmIntegrity_MD5_128, nil
	case "hmac-sha256-128":
		return RAKPAlgorithmIntegrity_HMAC_SHA256_128, nil
	default:
		return RAKPAlgorithmIntegrity_None, errors.New(s + " is unsupported confidentialityMethod")
	}
}

func (m RAKPAlgorithmIntegrityMethod) String() string {
	switch m {
	case RAKPAlgorithmIntegrity_None:
		return "none"
	case RAKPAlgorithmIntegrity_HMAC_SHA1_96:
		return "hmac-sha1-96"
	case RAKPAlgorithmIntegrity_HMAC_MD5_128:
		return "hmac-md5-128"
	case RAKPAlgorithmIntegrity_MD5_128:
		return "md5-128"
	case RAKPAlgorithmIntegrity_HMAC_SHA256_128:
		return "hmac-sha256-128"
	default:
		return "unknow(" + strconv.FormatInt(int64(m), 10) + ")"
	}
}

// truncatingMAC wraps around a hash.Hash and truncates the output digest to
// a given size.
type truncatingMAC struct {
	length int
	hmac   hash.Hash
}

func (t truncatingMAC) Write(data []byte) (int, error) {
	return t.hmac.Write(data)
}

func (t truncatingMAC) Sum(in []byte) []byte {
	out := t.hmac.Sum(in)
	return out[:len(in)+t.length]
}

var RAKPAlgorithmIntegrityKeychangeLengths = []int{0, 12, 16, 16, 16}
var RAKPAlgorithmIntegritySignLengths = []int{0, 20, 20, 20, 20}
var RAKPAlgorithmIntegrityMethods = []func(key, data []byte) []byte{
	// RAKPAlgorithmIntegrity_None
	func(key, data []byte) []byte {
		return data
	},

	// RAKPAlgorithmIntegrity_HMAC_SHA1_96
	func(key, data []byte) []byte {
		var h = hmac.New(sha1.New, key)
		h.Write(data)
		return h.Sum(nil)
	},

	// RAKPAlgorithmIntegrity_HMAC_MD5_128
	func(key, data []byte) []byte {
		var h = hmac.New(md5.New, key)
		h.Write(data)
		return h.Sum(nil)
	},

	// RAKPAlgorithmIntegrity_MD5_128
	func(key, data []byte) []byte {
		var h = md5.New()
		h.Write(data)
		return h.Sum(nil)
	},

	// RAKPAlgorithmIntegrity_HMAC_SHA256_128
	func(key, data []byte) []byte {
		var h = hmac.New(sha256.New, key)
		h.Write(data)
		return h.Sum(nil)
	},
}

func (self RAKPAlgorithmIntegrityMethod) Gen(key, data []byte) []byte {
	return RAKPAlgorithmIntegrityMethods[self](key, data)
}

func (self RAKPAlgorithmIntegrityMethod) KeyExchangeSize() int {
	return RAKPAlgorithmIntegrityKeychangeLengths[self]
}

func (self RAKPAlgorithmIntegrityMethod) SignSize() int {
	return RAKPAlgorithmIntegritySignLengths[self]
}

type RAKPAlgorithmConfidentialityMethod uint8

const (
	RAKPAlgorithmEncryto_None        RAKPAlgorithmConfidentialityMethod = 0
	RAKPAlgorithmEncryto_AES_CBC_128 RAKPAlgorithmConfidentialityMethod = 1
	RAKPAlgorithmEncryto_XRC4_128    RAKPAlgorithmConfidentialityMethod = 2
	RAKPAlgorithmEncryto_XRC4_40     RAKPAlgorithmConfidentialityMethod = 3
)

func ParseConfidentialityMethod(s string) (RAKPAlgorithmConfidentialityMethod, error) {
	switch strings.ToLower(s) {
	case "", "none":
		return RAKPAlgorithmEncryto_None, nil
	case "aes-cbc-128":
		return RAKPAlgorithmEncryto_AES_CBC_128, nil
	case "xrc4-128":
		return RAKPAlgorithmEncryto_XRC4_128, nil
	case "xrc4-40":
		return RAKPAlgorithmEncryto_XRC4_40, nil
	default:
		return RAKPAlgorithmEncryto_None, errors.New(s + " is unsupported confidentialityMethod")
	}
}

func (m RAKPAlgorithmConfidentialityMethod) String() string {
	switch m {
	case RAKPAlgorithmEncryto_None:
		return "none"
	case RAKPAlgorithmEncryto_AES_CBC_128:
		return "aes-cbc"
	case RAKPAlgorithmEncryto_XRC4_128:
		return "xrc4-128"
	case RAKPAlgorithmEncryto_XRC4_40:
		return "xrc4-40"
	default:
		return "unknow(" + strconv.FormatInt(int64(m), 10) + ")"
	}
}

var RAKPAlgorithmConfidentialityEncryptMethods = []func(key []byte) func(in, out []byte) ([]byte, error){
	func(key []byte) func(in, out []byte) ([]byte, error) {
		return func(in, out []byte) ([]byte, error) {
			copy(out, in)
			return out[:len(in)], nil
		}
	},

	func(key []byte) func(in, out []byte) ([]byte, error) {
		block, err := aes.NewCipher(key[:16])
		return func(in, out []byte) ([]byte, error) {
			if err != nil {
				return nil, err
			}
			padLen := (len(in) + 1) % 16
			if padLen != 0 {
				padLen = 16 - padLen
			}
			for i := 1; i <= padLen; i++ {
				in = append(in, byte(i))
			}
			in = append(in, byte(padLen))

			fmt.Println(len(out), (16 + len(in)))
			if len(out) != (16 + len(in)) {
				for i := len(out); i <= (16 + len(in)); i++ {
					fmt.Println("=====", i, len(out), (16 + len(in)))
					out = append(out, byte(0))
				}
				fmt.Println(len(out), (16 + len(in)))
			}

			iv := makeInitializationVector(out[:16])

			fmt.Println("key = ", fmt.Sprintf("%x", key))
			fmt.Println("iv  = ", fmt.Sprintf("%x", iv))
			fmt.Println("in  = ", fmt.Sprintf("%x", in))
			aes := cipher.NewCBCEncrypter(block, iv)
			aes.CryptBlocks(out[16:], in)
			fmt.Println("out = ", fmt.Sprintf("%x", out[16:16+len(in)]))
			return out[:16+len(in)], nil
		}
	},

	func(key []byte) func(in, out []byte) ([]byte, error) {
		return func(in, out []byte) ([]byte, error) {
			copy(out, in)
			return out[:len(in)], nil
		}
	},

	func(key []byte) func(in, out []byte) ([]byte, error) {
		return func(in, out []byte) ([]byte, error) {
			copy(out, in)
			return out[:len(in)], nil
		}
	},
}

var RAKPAlgorithmConfidentialityDecryptMethods = []func(key []byte) func(in, out []byte) ([]byte, error){
	func(key []byte) func(in, out []byte) ([]byte, error) {
		return func(in, out []byte) ([]byte, error) {
			copy(out, in)
			return out[:len(in)], nil
		}
	},

	func(key []byte) func(in, out []byte) ([]byte, error) {
		block, err := aes.NewCipher(key[:16])
		return func(in, out []byte) ([]byte, error) {
			if err != nil {
				return nil, err
			}

			iv := in[:16]
			fmt.Println("key = ", fmt.Sprintf("%x", key))
			fmt.Println("iv  = ", fmt.Sprintf("%x", iv))
			fmt.Println("in  = ", fmt.Sprintf("%x", in[16:]))
			aes := cipher.NewCBCDecrypter(block, iv)
			aes.CryptBlocks(out[:len(in)-16], in[16:])
			fmt.Println("out = ", fmt.Sprintf("%x", out[:len(in)-16]))
			return out[:len(in)-16], nil
		}
	},

	func(key []byte) func(in, out []byte) ([]byte, error) {
		return func(in, out []byte) ([]byte, error) {
			copy(out, in)
			return out[:len(in)], nil
		}
	},

	func(key []byte) func(in, out []byte) ([]byte, error) {
		return func(in, out []byte) ([]byte, error) {
			copy(out, in)
			return out[:len(in)], nil
		}
	},
}

func (self RAKPAlgorithmConfidentialityMethod) Encrypt(key []byte) func(in, out []byte) ([]byte, error) {
	return RAKPAlgorithmConfidentialityEncryptMethods[self](key)
}

func (self RAKPAlgorithmConfidentialityMethod) Decrypt(key []byte) func(in, out []byte) ([]byte, error) {
	return RAKPAlgorithmConfidentialityDecryptMethods[self](key)
}

type AuthenticationPayload struct {
	//PayloadType uint8
	//Reserved1   uint16
	//Length    uint8 // 8
	Algorithm RAKPAlgorithmAuthMethod
	//Reserved2   [3]uint8
}

func (self *AuthenticationPayload) WriteBytes(w *Writer) {
	w.WriteUint8(0)  // 0 is Authentication Algorithm
	w.WriteUint16(0) // Reserved1
	w.WriteUint8(8)  // Length is 8
	w.WriteUint8(uint8(self.Algorithm))
	w.WriteUint8(0) // Reserved Padding
	w.WriteUint8(0) // Reserved Padding
	w.WriteUint8(0) // Reserved Padding
}

func (self *AuthenticationPayload) ReadBytes(r *Reader) {
	if r.Len() < 8 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(8)
	if 0 != int(bs[0]) {
		r.SetError(errors.New("Authentication is excepted"))
		return
	}

	//self.PayloadType = uint8(bs[0])
	//self.Reserved1 = binary.LittleEndian.Uint16(bs[1:])
	//self.Length = uint8(bs[2])
	self.Algorithm = RAKPAlgorithmAuthMethod(bs[4])
}

type IntegrityPayload struct {
	//PayloadType uint8
	//Reserved1   uint16
	//Length      uint8 // 8
	Algorithm RAKPAlgorithmIntegrityMethod
	//Reserved2   [3]uint8
}

func (self *IntegrityPayload) WriteBytes(w *Writer) {
	w.WriteUint8(1)  // 1 is Integrity Algorithm
	w.WriteUint16(0) // Reserved1
	w.WriteUint8(8)  // Length is 8
	w.WriteUint8(uint8(self.Algorithm))
	w.WriteUint8(0) // Reserved Padding
	w.WriteUint8(0) // Reserved Padding
	w.WriteUint8(0) // Reserved Padding
}

func (self *IntegrityPayload) ReadBytes(r *Reader) {
	if r.Len() < 8 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(8)
	// if 1 != int(bs[0]) {
	//  r.SetError(errors.New("Integrity is excepted."))
	//  return
	// }
	//self.PayloadType = uint8(bs[0])
	//self.Reserved1 = binary.LittleEndian.Uint16(bs[1:])
	//self.Length = uint8(bs[2])
	self.Algorithm = RAKPAlgorithmIntegrityMethod(bs[4])
}

type ConfidentialityPayload struct {
	//PayloadType uint8
	//Reserved1   uint16
	//Length      uint8 // 8
	Algorithm RAKPAlgorithmConfidentialityMethod
	//Reserved2   [3]uint8
}

func (self *ConfidentialityPayload) WriteBytes(w *Writer) {
	w.WriteUint8(2)  // 2 is confidentiality Algorithm
	w.WriteUint16(0) // Reserved1
	w.WriteUint8(8)  // Length is 8
	w.WriteUint8(uint8(self.Algorithm))
	w.WriteUint8(0) // Reserved Padding
	w.WriteUint8(0) // Reserved Padding
	w.WriteUint8(0) // Reserved Padding
}

func (self *ConfidentialityPayload) ReadBytes(r *Reader) {
	if r.Len() < 8 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(8)
	if 2 != int(bs[0]) {
		r.SetError(errors.New("Confidentiality is excepted"))
		return
	}
	//self.PayloadType = uint8(bs[0])
	//self.Reserved1 = binary.LittleEndian.Uint16(bs[1:])
	//self.Length = uint8(bs[2])
	self.Algorithm = RAKPAlgorithmConfidentialityMethod(bs[4])
}
