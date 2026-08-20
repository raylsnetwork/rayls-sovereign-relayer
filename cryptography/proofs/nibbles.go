// This file is derived from github.com/zhangchiqing/merkle-patricia-trie,
// Copyright (c) 2020 Leo Zhang (zhangchiqing@gmail.com), licensed under the MIT
// License. Modifications Copyright 2026 Rayls Core Ltd. See the NOTICE file.
package proofs

import (
	"fmt"
)

const (
	// nibbleBitShift is the number of bits to shift for extracting the high nibble from a byte.
	nibbleBitShift = 4
	// nibbleMaxValue is the upper bound (exclusive) for valid nibble values (0x0 through 0xF).
	nibbleMaxValue = 16
)

type Nibble byte

func IsNibble(nibble byte) bool {
	n := int(nibble)
	// 0-9 && a-f
	return n >= 0 && n < nibbleMaxValue
}

func FromNibbleByte(n byte) (Nibble, error) {
	if !IsNibble(n) {
		return 0, fmt.Errorf("non-nibble byte: %v", n)
	}
	return Nibble(n), nil
}

// nibbles contain one nibble per byte
func FromNibbleBytes(nibbles []byte) ([]Nibble, error) {
	ns := make([]Nibble, 0, len(nibbles))
	for _, n := range nibbles {
		nibble, err := FromNibbleByte(n)
		if err != nil {
			return nil, fmt.Errorf("contains non-nibble byte: %w", err)
		}
		ns = append(ns, nibble)
	}
	return ns, nil
}

func FromByte(b byte) []Nibble {
	return []Nibble{
		Nibble(byte(b >> nibbleBitShift)),
		Nibble(byte(b % nibbleMaxValue)),
	}
}

func FromBytes(bs []byte) []Nibble {
	ns := make([]Nibble, 0, len(bs)*2)
	for _, b := range bs {
		ns = append(ns, FromByte(b)...)
	}
	return ns
}

func FromString(s string) []Nibble {
	return FromBytes([]byte(s))
}

// ToPrefixed add nibble prefix to a slice of nibbles to make its length even
// the prefix indicts whether a node is a leaf node.
func ToPrefixed(ns []Nibble, isLeafNode bool) []Nibble {
	// create prefix
	var prefixBytes []Nibble
	// odd number of nibbles
	if len(ns)%2 > 0 {
		prefixBytes = []Nibble{1}
	} else {
		// even number of nibbles
		prefixBytes = []Nibble{0, 0}
	}

	// append prefix to all nibble bytes
	prefixed := make([]Nibble, 0, len(prefixBytes)+len(ns))
	prefixed = append(prefixed, prefixBytes...)
	prefixed = append(prefixed, ns...)

	// update prefix if is leaf node
	if isLeafNode {
		prefixed[0] += 2
	}

	return prefixed
}

// ToBytes converts a slice of nibbles to a byte slice
// assuming the nibble slice has even number of nibbles.
func ToBytes(ns []Nibble) []byte {
	buf := make([]byte, 0, len(ns)/2)

	for i := 0; i < len(ns); i += 2 {
		b := byte(ns[i]<<nibbleBitShift) + byte(ns[i+1])
		buf = append(buf, b)
	}

	return buf
}

// [0,1,2,3], [0,1,2] => 3
// [0,1,2,3], [0,1,2,3] => 4
// [0,1,2,3], [0,1,2,3,4] => 4
func PrefixMatchedLen(node1 []Nibble, node2 []Nibble) int {
	matched := 0
	for i := 0; i < len(node1) && i < len(node2); i++ {
		n1, n2 := node1[i], node2[i]
		if n1 == n2 {
			matched++
		} else {
			break
		}
	}

	return matched
}
