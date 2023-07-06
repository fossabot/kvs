package storage

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/kvs"
)

type Value interface {
	TableName() string
	SetID(id uint32)
	Ref() interface{}
}

type Store struct {
	db  kvs.KVDB
	pks map[string]*badger.Sequence
}

func New(db kvs.KVDB) *Store {
	return &Store{db: db, pks: map[string]*badger.Sequence{}}
}

func (s Store) Save(owner kvs.UUID, value Value) error {
	rowID, err := nextRowID(s.db, value.TableName(), s.pks)
	if err != nil {
		return err
	}

	return saveValue(s.db, value.TableName(), owner, rowID, value)
}

func saveValue(db kvs.KVDB, tableName string, ownerID kvs.UUID, rowID uint32, v Value) error {
	if v == nil {
		return nil
	}
	entries := kvs.ConvertToEntries(tableName, ownerID, rowID, v)
	for _, e := range entries {
		if err := kvs.Store(db, e); err != nil {
			return err
		}
	}

	v.SetID(rowID)

	return nil
}

func nextRowID(db kvs.KVDB, tableName string, pks map[string]*badger.Sequence) (uint32, error) {
	seq, err := resolveSequence(db, tableName, pks)
	if err != nil {
		return 0, err
	}

	s, err := seq.Next()
	if err != nil {
		return 0, err
	}
	return uint32(s), nil
}

func nextSequence(seq *badger.Sequence) (uint32, error) {
	s, err := seq.Next()
	if err != nil {
		return 0, err
	}
	return uint32(s), nil
}

func resolveSequence(db kvs.KVDB, tableName string, pks map[string]*badger.Sequence) (*badger.Sequence, error) {
	seq, ok := pks[tableName]
	var err error
	if !ok {
		seq, err = db.GetSeq([]byte(tableName), 1)
		if err != nil {
			return nil, err
		}
		pks[tableName] = seq
	}

	return seq, nil
}