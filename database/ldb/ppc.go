// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ldb

import (
	"github.com/ppcsuite/ppcd/wire"
)

var blockMetaSuffix = []byte{'b', 'm'}

func shaMetaToKey(sha *wire.ShaHash) []byte {
	key := make([]byte, len(sha)+len(blockMetaSuffix))
	copy(key, sha[:])
	copy(key[len(sha):], blockMetaSuffix)
	return key
}

func (db *LevelDb) getBlkMeta(sha *wire.ShaHash) (rbuf []byte, err error) {
	key := shaMetaToKey(sha)
	rbuf, err = db.lDb.Get(key, db.ro)
	if err != nil {
		log.Tracef("failed to find meta for %v", sha)
	}
	return
}

func (db *LevelDb) setBlkMeta(sha *wire.ShaHash, buf []byte) {
	key := shaMetaToKey(sha)
	db.lBatch().Put(key, buf[:])
}
