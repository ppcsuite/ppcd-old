// Copyright (c) 2013-2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mably/ppcd/blockchain"
	"github.com/mably/ppcd/database"
	_ "github.com/mably/ppcd/database/ldb"
	_ "github.com/mably/ppcd/database/memdb"
	"github.com/mably/btcnet"
	"github.com/mably/btcutil"
)

// testDbType is the database backend type to use for the tests.
const testDbType = "memdb"

// testDbRoot is the root directory used to create all test databases.
const testDbRoot = "testdbs"

// filesExists returns whether or not the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// isSupportedDbType returns whether or not the passed database type is
// currently supported.
func isSupportedDbType(dbType string) bool {
	supportedDBs := database.SupportedDBs()
	for _, sDbType := range supportedDBs {
		if dbType == sDbType {
			return true
		}
	}

	return false
}

// chainSetup is used to create a new db and chain instance with the genesis
// block already inserted.  In addition to the new chain instnce, it returns
// a teardown function the caller should invoke when done testing to clean up.
func chainSetup(dbName string) (*blockchain.BlockChain, func(), error) {
	if !isSupportedDbType(testDbType) {
		return nil, nil, fmt.Errorf("unsupported db type %v", testDbType)
	}

	// Handle memory database specially since it doesn't need the disk
	// specific handling.
	var db database.Db
	var teardown func()
	if testDbType == "memdb" {
		ndb, err := database.CreateDB(testDbType)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating db: %v", err)
		}
		db = ndb

		// Setup a teardown function for cleaning up.  This function is
		// returned to the caller to be invoked when it is done testing.
		teardown = func() {
			db.Close()
		}
	} else {
		// Create the root directory for test databases.
		if !fileExists(testDbRoot) {
			if err := os.MkdirAll(testDbRoot, 0700); err != nil {
				err := fmt.Errorf("unable to create test db "+
					"root: %v", err)
				return nil, nil, err
			}
		}

		// Create a new database to store the accepted blocks into.
		dbPath := filepath.Join(testDbRoot, dbName)
		_ = os.RemoveAll(dbPath)
		ndb, err := database.CreateDB(testDbType, dbPath)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating db: %v", err)
		}
		db = ndb

		// Setup a teardown function for cleaning up.  This function is
		// returned to the caller to be invoked when it is done testing.
		teardown = func() {
			dbVersionPath := filepath.Join(testDbRoot, dbName+".ver")
			db.Sync()
			db.Close()
			os.RemoveAll(dbPath)
			os.Remove(dbVersionPath)
			os.RemoveAll(testDbRoot)
		}
	}

	// Insert the main network genesis block.  This is part of the initial
	// database setup.
	genesisBlock := btcutil.NewBlock(btcnet.MainNetParams.GenesisBlock)
	_, err := db.InsertBlock(genesisBlock)
	if err != nil {
		teardown()
		err := fmt.Errorf("failed to insert genesis block: %v", err)
		return nil, nil, err
	}

	chain := blockchain.New(db, &btcnet.MainNetParams, nil)
	return chain, teardown, nil
}
