// This file is derived from github.com/zhangchiqing/merkle-patricia-trie,
// Copyright (c) 2020 Leo Zhang (zhangchiqing@gmail.com), licensed under the MIT
// License. Modifications Copyright 2026 Rayls Core Ltd. See the NOTICE file.
package proofs

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

type Node interface {
	Hash() ([]byte, error)
	Raw() ([]interface{}, error)
}

func Hash(node Node) ([]byte, error) {
	if IsEmptyNode(node) {
		return EmptyNodeHash, nil
	}
	return node.Hash()
}

func Serialize(node Node) ([]byte, error) {
	var raw interface{}

	if IsEmptyNode(node) {
		raw = EmptyNodeRaw
	} else {
		var err error
		raw, err = node.Raw()
		if err != nil {
			return nil, fmt.Errorf("getting node raw data: %w", err)
		}
	}

	rlpVar, err := rlp.EncodeToBytes(raw)
	if err != nil {
		return nil, fmt.Errorf("RLP encoding node: %w", err)
	}

	return rlpVar, nil
}
