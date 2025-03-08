package synccredmap

import "sync"

type SyncCredentialsMap struct {
  mx sync.RWMutex
  m  map[string]string
}

func (sm *SyncCredentialsMap) Map() map[string]string {
  return sm.m
}

func NewSyncCredentialsMap() *SyncCredentialsMap {
  return &SyncCredentialsMap{
    m: make(map[string]string),
  }
}

// sync read value from map by key. returns false if value not present
func (sm *SyncCredentialsMap) Load(key string) (val string, ok bool) {
  sm.mx.RLock()
  defer sm.mx.RUnlock()

  val, ok = sm.m[key]
  return val, ok
}

// sync plase value inside map by key, does not check if value already exists
func (sm *SyncCredentialsMap) Store(key, val string) {
  sm.mx.Lock()
  defer sm.mx.Unlock()

  sm.m[key] = val
}

// sync delete value from map by key
func (sm *SyncCredentialsMap) Delete(key string) {
  sm.mx.Lock()
  defer sm.mx.Unlock()

  delete(sm.m, key)
}

// if map[key] is empty places new data in map by that key and returns true,
// otherwise returns actual value and false
func (sm *SyncCredentialsMap) LoadOrStore(key, value string) (string, bool) {
  sm.mx.RLock()

  actual, isLoaded := sm.m[key]
  sm.mx.RUnlock()

  if isLoaded {
    return actual, false
  }

  sm.mx.Lock()
  defer sm.mx.Unlock()

  sm.m[key] = value
  return value, true
}
