// Example packing demonstrates binary serialization with iproto Pack/Unpack.
//
// It shows how to encode and decode primitive types, use BER encoding,
// and serialize structs with iproto tags.
package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Educentr/go-iproto/iproto"
)

// Message demonstrates struct packing with mixed encoding modes.
// The Count field uses BER (variable-length) encoding, while Name uses
// the default fixed-length little-endian encoding.
type Message struct {
	Count uint32 `iproto:"ber"`
	Name  string
}

func main() {
	// --- Primitive types ---
	fmt.Println("=== Primitive Types ===")

	// Pack a uint32 in little-endian.
	data, err := iproto.Pack(uint32(258))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("uint32(258) -> %s\n", hex.EncodeToString(data))

	// Pack a string (length-prefixed).
	data, err = iproto.Pack("hello")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("string(\"hello\") -> %s\n", hex.EncodeToString(data))

	// Pack a uint16.
	data, err = iproto.Pack(uint16(1024))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("uint16(1024) -> %s\n", hex.EncodeToString(data))

	// --- BER encoding ---
	fmt.Println("\n=== BER Encoding ===")

	// BER encodes integers with variable length (7 bits per byte).
	data, err = iproto.PackBER(uint32(5))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BER uint32(5)   -> %s (1 byte)\n", hex.EncodeToString(data))

	data, err = iproto.PackBER(uint32(128))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BER uint32(128) -> %s (2 bytes)\n", hex.EncodeToString(data))

	data, err = iproto.PackBER(uint32(100000))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BER uint32(100000) -> %s (3 bytes)\n", hex.EncodeToString(data))

	// --- Struct packing with tags ---
	fmt.Println("\n=== Struct Packing ===")

	msg := Message{Count: 258, Name: "world"}

	data, err = iproto.Pack(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Message{Count:258, Name:\"world\"} -> %s\n", hex.EncodeToString(data))

	// Unpack the struct back.
	var decoded Message
	_, err = iproto.Unpack(data, &decoded)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Unpacked: Count=%d, Name=%q\n", decoded.Count, decoded.Name)

	// --- Round-trip demo ---
	fmt.Println("\n=== Round-Trip ===")

	original := uint32(42)
	packed, _ := iproto.Pack(original)

	var unpacked uint32
	_, err = iproto.Unpack(packed, &unpacked)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("original=%d -> packed=%s -> unpacked=%d\n",
		original, hex.EncodeToString(packed), unpacked)
}
