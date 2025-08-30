package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, number)

	stored, err := store.Get(number)
	require.NoError(t, err)

	expected := parcel
	expected.Number = number
	require.Equal(t, expected, stored)

	err = store.Delete(number)
	require.NoError(t, err)

	_, err = store.Get(number)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, number)

	newAddress := "new test"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	stored, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, newAddress, stored.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, number)

	newStatus := ParcelStatusSent
	err = store.SetStatus(number, newStatus)
	require.NoError(t, err)

	stored, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, newStatus, stored.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	store := NewParcelStore(db)

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])

		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	for _, stored := range storedParcels {
		expected, ok := parcelMap[stored.Number]
		require.True(t, ok)

		require.Equal(t, expected, stored)
	}
}
