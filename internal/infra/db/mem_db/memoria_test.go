package memoria

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoriaDatabase(t *testing.T) {
	db := NewMemoriaDatabase[string]()
	assert.NotNil(t, db)
	assert.NotNil(t, db.data)
	assert.Empty(t, db.data)
}

func TestMemoriaDatabase_Save(t *testing.T) {
	db := NewMemoriaDatabase[string]()

	err := db.Save("key1", "value1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(db.data))
	assert.Equal(t, "value1", db.data["key1"])

	err = db.Save("key2", "value2")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(db.data))
	assert.Equal(t, "value2", db.data["key2"])

	// Overwrite existing key
	err = db.Save("key1", "newValue1")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(db.data))
	assert.Equal(t, "newValue1", db.data["key1"])
}

func TestMemoriaDatabase_Get(t *testing.T) {
	db := NewMemoriaDatabase[string]()
	db.Save("key1", "value1")
	db.Save("key2", "value2")

	val, err := db.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	val, err = db.Get("key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", val)

	// Not found case
	val, err = db.Get("nonExistentKey")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errors.New("not found")))
	assert.Empty(t, val) // Zero value for string

	// Test with custom struct
	type MyStruct struct {
		Name string
		Age  int
	}

	dbStruct := NewMemoriaDatabase[MyStruct]()
	structVal := MyStruct{Name: "test", Age: 10}
	dbStruct.Save("structKey", structVal)

	retrievedStruct, err := dbStruct.Get("structKey")
	assert.NoError(t, err)
	assert.Equal(t, structVal, retrievedStruct)
}

func TestMemoriaDatabase_Remove(t *testing.T) {
	db := NewMemoriaDatabase[string]()
	db.Save("key1", "value1")
	db.Save("key2", "value2")

	// Remove existing key
	err := db.Remove("key1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(db.data))
	_, ok := db.data["key1"]
	assert.False(t, ok)

	// Try to remove non-existent key
	err = db.Remove("nonExistentKey")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errors.New("not found")))
	assert.Equal(t, 1, len(db.data)) // Should not change size

	// Remove last key
	err = db.Remove("key2")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(db.data))
	_, ok = db.data["key2"]
	assert.False(t, ok)
}
