/*
 * Copyright (C) 2020-2021 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package hashing

import (

	// nolint:gosec
	"crypto/md5"

	// nolint:gosec
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"strings"

	"github.com/OneOfOne/xxhash"
	"github.com/spaolacci/murmur3"

	"github.com/ARM-software/golang-utils/utils/commonerrors"
)

const (
	HashMd5    = "MD5"
	HashSha256 = "SHA256"
	HashSha1   = "SHA1"
	HashMurmur = "Murmur"
	HashXXHash = "xxhash" //https://github.com/OneOfOne/xxhash
)

type IHash interface {
	Calculate(reader io.Reader) (string, error)
	GetType() string
}

type hashingAlgo struct {
	Hash hash.Hash
	Type string
}

func (h *hashingAlgo) Calculate(r io.Reader) (hashN string, err error) {
	if r == nil {
		err = commonerrors.ErrUndefined
		return
	}
	_, err = io.Copy(h.Hash, r)
	if err != nil {
		return
	}
	hashN = hex.EncodeToString(h.Hash.Sum(nil))
	h.Hash.Reset()
	return
}

func (h *hashingAlgo) GetType() string {
	return h.Type
}

func NewHashingAlgorithm(htype string) (IHash, error) {
	var hash hash.Hash
	switch htype {
	case HashMd5:
		hash = md5.New() //nolint:gosec
	case HashSha1:
		hash = sha1.New() //nolint:gosec
	case HashSha256:
		hash = sha256.New()
	case HashMurmur:
		hash = murmur3.New64()
	case HashXXHash:
		hash = xxhash.New64()
	}

	if hash == nil {
		return nil, commonerrors.ErrNotFound
	}
	return &hashingAlgo{
		Hash: hash,
		Type: htype,
	}, nil
}

func CalculateMD5Hash(text string) string {
	return CalculateHash(text, HashMd5)
}

func CalculateHash(text, htype string) string {
	hashing, err := NewHashingAlgorithm(htype)
	if err != nil {
		return ""
	}
	hash, err := hashing.Calculate(strings.NewReader(text))
	if err != nil {
		return ""
	}
	return hash
}
