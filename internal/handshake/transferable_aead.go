package handshake

import (
	"crypto/cipher"
	"crypto/tls"
	"encoding/json"
	"github.com/lucas-clemente/quic-go/internal/qtls"
	"reflect"
	"unsafe"
)

type transferableAeadAESGCMTLS13 struct {
	NonceMask [12]byte
	// Ks is the key schedule, the length of which depends on the size of
	// the AES key.
	Ks []uint32
	// ProductTable contains pre-computed multiples of the binary-field
	// element used in GHASH.
	ProductTable [256]byte
	// NonceSize contains the expected size of the nonce, in bytes.
	NonceSize int
	// TagSize contains the size of the tag, in bytes.
	TagSize int
}
//TODO use generics when supported
func getField(structValue reflect.Value, fieldName string) reflect.Value{
	for structValue.Kind() == reflect.Ptr || structValue.Kind() == reflect.Interface {
		structValue = structValue.Elem()
	}
	fieldValue := structValue.FieldByName(fieldName)
	fieldValue = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
	return fieldValue
}

//TODO use generics when supported
func setField(structValue reflect.Value, fieldName string, fieldValue reflect.Value) {
	getField(structValue, fieldName).Set(fieldValue)
}


func marshal(suite *qtls.CipherSuiteTLS13, aead cipher.AEAD) []byte {
	if suite.ID != tls.TLS_AES_128_GCM_SHA256{
		panic("implement me")
	}
	transferable := transferableAeadAESGCMTLS13{}
	xorNonceAEAD := reflect.ValueOf(aead)
	gcmAsm := getField(xorNonceAEAD, "aead")
	transferable.NonceMask = getField(xorNonceAEAD, "nonceMask").Interface().([12]byte)
	transferable.Ks = getField(gcmAsm, "ks").Interface().([]uint32)
	transferable.ProductTable = getField(gcmAsm, "productTable").Interface().([256]uint8)
	transferable.NonceSize = getField(gcmAsm, "nonceSize").Interface().(int)
	transferable.TagSize = getField(gcmAsm, "tagSize").Interface().(int)

	b, _ := json.Marshal(&transferable)
	return b
}

func unmarshal(suite *qtls.CipherSuiteTLS13, data []byte) cipher.AEAD {
	if suite.ID != tls.TLS_AES_128_GCM_SHA256{
		panic("implement me")
	}
	aead := suite.AEAD(make([]byte, 16), make([]byte, 12))

	transferable := transferableAeadAESGCMTLS13{}
	err := json.Unmarshal(data, &transferable)
	if err != nil {
		panic("unmarshal failed")
	}
	xorNonceAEAD := reflect.ValueOf(aead)
	gcmAsm := getField(xorNonceAEAD, "aead")
	setField(xorNonceAEAD, "nonceMask", reflect.ValueOf(transferable.NonceMask))
	setField(gcmAsm, "ks", reflect.ValueOf(transferable.Ks))
	setField(gcmAsm, "productTable", reflect.ValueOf(transferable.ProductTable))
	setField(gcmAsm, "nonceSize", reflect.ValueOf(transferable.NonceSize))
	setField(gcmAsm, "tagSize", reflect.ValueOf(transferable.TagSize))
	return aead
}

func clone(suite *qtls.CipherSuiteTLS13, aead cipher.AEAD) cipher.AEAD {
	return unmarshal(suite, marshal(suite, aead))
}