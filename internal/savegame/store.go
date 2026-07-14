package savegame

import "time"

// Store abstracts browser localStorage and the native in-memory adapter.
type Store interface {
	Load(key string) ([]byte, bool, error)
	Save(key string, value []byte) error
	Remove(key string) error
}

func Autosave(store Store, key string, payload any, now time.Time) error {
	model, err := New(payload, now)
	if err != nil {
		return err
	}
	raw, err := jsonMarshal(model)
	if err != nil {
		return err
	}
	return store.Save(key, raw)
}
