package standalone_storage

import (
	"github.com/Connor1996/badger"
	"github.com/pingcap-incubator/tinykv/kv/config"
	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/kv/util/engine_util"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// StandAloneStorage is an implementation of `Storage` for a single-node TinyKV instance. It does not
// communicate with other nodes and all data is stored locally.
type StandAloneStorage struct {
	db *badger.DB
}

func NewStandAloneStorage(conf *config.Config) *StandAloneStorage {
	badgerDB := engine_util.CreateDB(conf.DBPath, conf.Raft)
	return &StandAloneStorage{
		db: badgerDB,
	}
}

type StorageReader struct {
	Txn *badger.Txn
}

func (r StorageReader) GetCF(cf string, key []byte) ([]byte, error) {
	value, err := engine_util.GetCFFromTxn(r.Txn, cf, key)
	// check if we should handle this here ?
	if err != nil && err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return value, err
}

func (r StorageReader) IterCF(cf string) engine_util.DBIterator {
	return engine_util.NewCFIterator(cf, r.Txn)
}

func (r StorageReader) Close() {
	r.Txn.Discard()
}

func (s *StandAloneStorage) Start() error {
	return nil
}

func (s *StandAloneStorage) Stop() error {
	return s.db.Close()
}

func (s *StandAloneStorage) Reader(ctx *kvrpcpb.Context) (storage.StorageReader, error) {
	txn := s.db.NewTransaction(false)
	reader := &StorageReader{
		Txn: txn,
	}
	return reader, nil
}

func (s *StandAloneStorage) Write(ctx *kvrpcpb.Context, batch []storage.Modify) error {
	txn := s.db.NewTransaction(true)
	for _, b := range batch {
		switch b.Data.(type) {
		case storage.Put:
			txn.Set(engine_util.KeyWithCF(b.Cf(), b.Key()), b.Value())
		case storage.Delete:
			txn.Delete(engine_util.KeyWithCF(b.Cf(), b.Key()))
		}
	}
	return txn.Commit()
}
