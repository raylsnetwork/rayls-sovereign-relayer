// This file is derived from github.com/zhangchiqing/merkle-patricia-trie,
// Copyright (c) 2020 Leo Zhang (zhangchiqing@gmail.com), licensed under the MIT
// License. Modifications Copyright 2026 Rayls Core Ltd. See the NOTICE file.
package proofs

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

// hashSize is the byte length of a Keccak-256 hash.
const hashSize = 32

type BranchNode struct {
	Branches [16]Node
	Value    []byte
}

func NewBranchNode() *BranchNode {
	return &BranchNode{
		Branches: [16]Node{},
	}
}

func (b BranchNode) Hash() ([]byte, error) {
	s, err := b.Serialize()
	if err != nil {
		return nil, fmt.Errorf("hashing branch node: %w", err)
	}
	return crypto.Keccak256(s), nil
}

func (b *BranchNode) SetBranch(nibble Nibble, node Node) {
	b.Branches[int(nibble)] = node
}

func (b *BranchNode) RemoveBranch(nibble Nibble) {
	b.Branches[int(nibble)] = nil
}

func (b *BranchNode) SetValue(value []byte) {
	b.Value = value
}

func (b *BranchNode) RemoveValue() {
	b.Value = nil
}

func (b BranchNode) Raw() ([]interface{}, error) {
	hashes := make([]interface{}, 17)
	for i := 0; i < 16; i++ {
		if b.Branches[i] == nil {
			hashes[i] = EmptyNodeRaw
		} else {
			node := b.Branches[i]
			s, err := Serialize(node)
			if err != nil {
				return nil, fmt.Errorf("serializing branch %d: %w", i, err)
			}
			if len(s) >= hashSize {
				h, err := node.Hash()
				if err != nil {
					return nil, fmt.Errorf("hashing branch %d: %w", i, err)
				}
				hashes[i] = h
			} else {
				// if node can be serialized to less than 32 bits, then
				// use Serialized directly.
				// it has to be ">=", rather than ">",
				// so that when deserialized, the content can be distinguished
				// by length
				r, err := node.Raw()
				if err != nil {
					return nil, fmt.Errorf("getting raw branch %d: %w", i, err)
				}
				hashes[i] = r
			}
		}
	}

	hashes[16] = b.Value
	return hashes, nil
}

func (b BranchNode) Serialize() ([]byte, error) {
	return Serialize(b)
}

func (b BranchNode) HasValue() bool {
	return b.Value != nil
}
