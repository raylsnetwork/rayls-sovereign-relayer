// This file is derived from github.com/zhangchiqing/merkle-patricia-trie,
// Copyright (c) 2020 Leo Zhang (zhangchiqing@gmail.com), licensed under the MIT
// License. Modifications Copyright 2026 Rayls Core Ltd. See the NOTICE file.
package proofs

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type ExtensionNode struct {
	Path []Nibble
	Next Node
}

func NewExtensionNode(nibbles []Nibble, next Node) *ExtensionNode {
	return &ExtensionNode{
		Path: nibbles,
		Next: next,
	}
}

func (e ExtensionNode) Hash() ([]byte, error) {
	s, err := e.Serialize()
	if err != nil {
		return nil, fmt.Errorf("hashing extension node: %w", err)
	}
	return crypto.Keccak256(s), nil
}

func (e ExtensionNode) Raw() ([]interface{}, error) {
	hashes := make([]interface{}, 2)
	hashes[0] = ToBytes(ToPrefixed(e.Path, false))
	s, err := Serialize(e.Next)
	if err != nil {
		return nil, fmt.Errorf("serializing extension next: %w", err)
	}
	if len(s) >= hashSize {
		h, err := e.Next.Hash()
		if err != nil {
			return nil, fmt.Errorf("hashing extension next: %w", err)
		}
		hashes[1] = h
	} else {
		r, err := e.Next.Raw()
		if err != nil {
			return nil, fmt.Errorf("getting raw extension next: %w", err)
		}
		hashes[1] = r
	}
	return hashes, nil
}

func (e ExtensionNode) Serialize() ([]byte, error) {
	return Serialize(e)
}
