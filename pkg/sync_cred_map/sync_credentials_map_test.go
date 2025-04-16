package synccredmap_test

import (
	"testing"

	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
	"github.com/stretchr/testify/assert"
)

func TestStoreAndLoad(t *testing.T) {
	scm := synccredmap.NewSyncCredentialsMap()
	key := "key"
	value := "value"

	val, ok := scm.Load(key)
	assert.False(t, ok)
	assert.Empty(t, val)

	scm.Store(key, value)

	loadedVal, ok := scm.Load(key)
	assert.True(t, ok)
	assert.Equal(t, value, loadedVal)
}

func TestDelete(t *testing.T) {
	scm := synccredmap.NewSyncCredentialsMap()
	key := "key"
	value := "value"

	scm.Store(key, value)

	_, ok := scm.Load(key)
	assert.True(t, ok)

	scm.Delete(key)

	_, ok = scm.Load(key)
	assert.False(t, ok)
}

func TestLoadOrStore(t *testing.T) {
	scm := synccredmap.NewSyncCredentialsMap()
	key := "key"
	value1 := "value1"
	value2 := "value2"

	storedVal, loaded := scm.LoadOrStore(key, value1)
	assert.Equal(t, value1, storedVal)
	assert.True(t, loaded)

	storedVal, loaded = scm.LoadOrStore(key, value2)
	assert.Equal(t, value1, storedVal)
	assert.False(t, loaded)
}

func TestMap(t *testing.T) {
	scm := synccredmap.NewSyncCredentialsMap()

	entries := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range entries {
		scm.Store(k, v)
	}

	m := scm.Map()

	assert.Equal(t, entries, m)
}
