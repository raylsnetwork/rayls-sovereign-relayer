// This file is derived from github.com/zhangchiqing/merkle-patricia-trie,
// Copyright (c) 2020 Leo Zhang (zhangchiqing@gmail.com), licensed under the MIT
// License. Modifications Copyright 2026 Rayls Core Ltd. See the NOTICE file.
package proofs

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type ProofDB struct {
	kv map[string][]byte
}

func NewProofDB() *ProofDB {
	return &ProofDB{
		kv: make(map[string][]byte),
	}
}

func (w *ProofDB) Put(key []byte, value []byte) error {
	keyS := fmt.Sprintf("%x", key)
	w.kv[keyS] = value

	return nil
}

func (w *ProofDB) Delete(key []byte) error {
	keyS := fmt.Sprintf("%x", key)
	delete(w.kv, keyS)
	return nil
}

func (w *ProofDB) Has(key []byte) (bool, error) {
	// Defense-in-depth: error rather than dereference a nil receiver.
	if w == nil {
		return false, errors.New("cannot query a nil ProofDB")
	}
	keyS := fmt.Sprintf("%x", key)
	_, ok := w.kv[keyS]
	return ok, nil
}

func (w *ProofDB) Get(key []byte) ([]byte, error) {
	// Defense-in-depth: error rather than dereference a nil receiver.
	if w == nil {
		return nil, errors.New("cannot read from a nil ProofDB")
	}
	keyS := fmt.Sprintf("%x", key)
	val, ok := w.kv[keyS]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return val, nil
}

func (w *ProofDB) Serialize() [][]byte {
	nodes := make([][]byte, 0, len(w.kv))
	for _, value := range w.kv {
		nodes = append(nodes, value)
	}
	return nodes
}

func (w *ProofDB) Export() ([]byte, error) {
	proofMap := make(map[string][]byte, len(w.kv))
	for k, v := range w.kv {
		proofMap[k] = v
	}

	proofBytes, err := json.Marshal(proofMap)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("marshaling proof: %w", err))
	}

	return proofBytes, nil
}

// Import deserializes bytes produced by Export into a new ProofDB. It is a
// package-level constructor (no receiver) since it always builds a fresh DB.
func Import(data []byte) (*ProofDB, error) {
	proofMap := make(map[string][]byte)
	if err := json.Unmarshal(data, &proofMap); err != nil {
		return nil, withstack.Wrap(fmt.Errorf("unmarshaling proof: %w", err))
	}

	proofDB := NewProofDB()
	for k, v := range proofMap {
		hash, err := hex.DecodeString(k)
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("decoding proof key %q: %w", k, err))
		}
		if err := proofDB.Put(hash, v); err != nil {
			return nil, withstack.Wrap(fmt.Errorf("storing proof node: %w", err))
		}
	}
	return proofDB, nil
}

// Prove returns the merkle proof for the given key, which is
func (t *Trie) Prove(key []byte) (*ProofDB, bool) {
	proof := NewProofDB()
	node := t.root
	nibbles := FromBytes(key)

	for {
		h, err := Hash(node)
		if err != nil {
			return nil, false
		}
		s, err := Serialize(node)
		if err != nil {
			return nil, false
		}
		if err := proof.Put(h, s); err != nil {
			return nil, false
		}

		if IsEmptyNode(node) {
			return nil, false
		}

		if leaf, ok := node.(*LeafNode); ok {
			matched := PrefixMatchedLen(leaf.Path, nibbles)
			if matched != len(leaf.Path) || matched != len(nibbles) {
				return nil, false
			}

			return proof, true
		}

		if branch, ok := node.(*BranchNode); ok {
			if len(nibbles) == 0 {
				return proof, branch.HasValue()
			}

			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			node = branch.Branches[b]
			continue
		}

		if ext, ok := node.(*ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Path, nibbles)
			// E 01020304
			//   010203
			if matched < len(ext.Path) {
				return nil, false
			}

			nibbles = nibbles[matched:]
			node = ext.Next
			continue
		}

		return nil, false
	}
}

// VerifyProof verify the proof for the given key under the given root hash using go-ethereum's VerifyProof implementation.
// It returns the value for the key if the proof is valid, otherwise error will be returned
func VerifyProof(rootHash common.Hash, key []byte, proof *ProofDB) (value []byte, err error) {
	if proof == nil {
		return nil, errors.New("cannot verify against a nil ProofDB")
	}
	return trie.VerifyProof(
		rootHash,
		key,
		proof,
	)
}

func GenerateTrie(txs []*types.Transaction) (*Trie, error) {
	newTrie := NewTrie()

	if len(txs) == 0 {
		return nil, withstack.Wrap(errors.New("no transactions in block to generate proof"))
	}
	for i, tx := range txs {
		key, err := rlp.EncodeToBytes(uint(i)) //nolint:gosec // i is a slice index, always non-negative
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("error encoding to bytes: %w", err))
		}
		transaction := fromEthTransaction(tx)
		rlpVar, err := transaction.GetRLP()
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("RLP encoding transaction %s: %w", tx.Hash().String(), err))
		}
		if err := newTrie.Put(key, rlpVar); err != nil {
			return nil, withstack.Wrap(fmt.Errorf("putting transaction %d into trie: %w", i, err))
		}
	}

	return newTrie, nil
}

func GenerateProofs(trie *Trie, txToProve uint, transactionRoot common.Hash) (*ProofDB, error) {
	key, err := rlp.EncodeToBytes(txToProve)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error encoding to bytes: %w", err))
	}

	proof, ok := trie.Prove(key)
	if !ok {
		return nil, withstack.Wrap(errors.New("error creating merkle proof"))
	}

	_, err = VerifyProof(transactionRoot, key, proof)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error verifying proof: %w", err))
	}

	return proof, nil
}
