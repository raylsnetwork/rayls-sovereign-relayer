// This file is derived from github.com/zhangchiqing/merkle-patricia-trie,
// Copyright (c) 2020 Leo Zhang (zhangchiqing@gmail.com), licensed under the MIT
// License. Modifications Copyright 2026 Rayls Core Ltd. See the NOTICE file.
package proofs

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type LeafNode struct {
	Path  []Nibble
	Value []byte
}

func NewLeafNodeFromNibbleBytes(nibbles []byte, value []byte) (*LeafNode, error) {
	ns, err := FromNibbleBytes(nibbles)
	if err != nil {
		return nil, fmt.Errorf("could not leaf node from nibbles: %w", err)
	}

	return NewLeafNodeFromNibbles(ns, value), nil
}

func NewLeafNodeFromNibbles(nibbles []Nibble, value []byte) *LeafNode {
	return &LeafNode{
		Path:  nibbles,
		Value: value,
	}
}

func NewLeafNodeFromKeyValue(key, value string) *LeafNode {
	return NewLeafNodeFromBytes([]byte(key), []byte(value))
}

func NewLeafNodeFromBytes(key, value []byte) *LeafNode {
	return NewLeafNodeFromNibbles(FromBytes(key), value)
}

func (l LeafNode) Hash() ([]byte, error) {
	s, err := l.Serialize()
	if err != nil {
		return nil, fmt.Errorf("hashing leaf node: %w", err)
	}
	return crypto.Keccak256(s), nil
}

func (l LeafNode) Raw() ([]interface{}, error) {
	path := ToBytes(ToPrefixed(l.Path, true))
	raw := []interface{}{path, l.Value}
	return raw, nil
}

func (l LeafNode) Serialize() ([]byte, error) {
	return Serialize(l)
}
