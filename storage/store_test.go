package storage_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/kvs"
	"github.com/tauraamui/kvs/storage"
)

type Balloon struct {
	ID    uint32 `mdb:"ignore"`
	Color string
	Size  int
}

func (b Balloon) TableName() string { return "balloons" }
func (b *Balloon) SetID(id uint32)  { b.ID = id }
func (b *Balloon) Ref() any         { return b }

type Cake struct {
	ID       uint32 `mdb:"ignore"`
	Type     string
	Calories int
}

func (b Cake) TableName() string { return "cakes" }
func (b *Cake) SetID(id uint32)  { b.ID = id }
func (b *Cake) Ref() any         { return b }

func TestStoreNewSuccess(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemKVDB()
	is.NoErr(err)
	defer db.Close()

	store := storage.New(db)
	is.True(store != nil)
}

func TestStoreBalloonsSuccess(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemKVDB()
	is.NoErr(err)
	defer db.Close()

	store := storage.New(db)
	is.True(store != nil)
	defer func() {
		is.NoErr(store.Close())
	}()

	bigRedBalloon := Balloon{Color: "RED", Size: 695}
	smallYellowBalloon := Balloon{Color: "YELLOW", Size: 112}
	mediumWhiteBalloon := Balloon{Color: "WHITE", Size: 366}
	is.NoErr(store.Save(kvs.RootOwner{}, &bigRedBalloon))
	is.NoErr(store.Save(kvs.RootOwner{}, &smallYellowBalloon))
	is.NoErr(store.Save(kvs.RootOwner{}, &mediumWhiteBalloon))

	is.Equal(bigRedBalloon.ID, uint32(0))
	is.Equal(smallYellowBalloon.ID, uint32(1))
	is.Equal(mediumWhiteBalloon.ID, uint32(2))
}

func TestStoreBalloonsAndCakesInSuccessionRetainsCorrectRowIDs(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemKVDB()
	is.NoErr(err)
	defer db.Close()

	store := storage.New(db)
	is.True(store != nil)
	defer func() {
		is.NoErr(store.Close())
	}()

	bigRedBalloon := Balloon{Color: "RED", Size: 695}
	disguistingVeganCake := Cake{Type: "INEDIBLE", Calories: -38}
	smallYellowBalloon := Balloon{Color: "YELLOW", Size: 112}
	healthyishCarrotCake := Cake{Type: "CARROT", Calories: 280}
	mediumWhiteBalloon := Balloon{Color: "WHITE", Size: 366}
	redVelvetCake := Cake{Type: "RED_VELVET", Calories: 410}

	is.NoErr(store.Save(kvs.RootOwner{}, &bigRedBalloon))
	is.NoErr(store.Save(kvs.RootOwner{}, &smallYellowBalloon))
	is.NoErr(store.Save(kvs.RootOwner{}, &mediumWhiteBalloon))

	is.NoErr(store.Save(kvs.RootOwner{}, &disguistingVeganCake))
	is.NoErr(store.Save(kvs.RootOwner{}, &healthyishCarrotCake))
	is.NoErr(store.Save(kvs.RootOwner{}, &redVelvetCake))

	is.Equal(bigRedBalloon.ID, uint32(0))
	is.Equal(disguistingVeganCake.ID, uint32(0))
	is.Equal(smallYellowBalloon.ID, uint32(1))
	is.Equal(healthyishCarrotCake.ID, uint32(1))
	is.Equal(mediumWhiteBalloon.ID, uint32(2))
	is.Equal(redVelvetCake.ID, uint32(2))
}
