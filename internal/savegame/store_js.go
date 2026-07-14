//go:build js && wasm

package savegame

import "syscall/js"

type LocalStore struct{}

func NewStore() *LocalStore { return &LocalStore{} }
func (LocalStore) Load(key string) ([]byte, bool, error) {
	v := js.Global().Get("localStorage").Call("getItem", key)
	if v.Type() == js.TypeNull {
		return nil, false, nil
	}
	return []byte(v.String()), true, nil
}
func (LocalStore) Save(key string, value []byte) error {
	js.Global().Get("localStorage").Call("setItem", key, string(value))
	return nil
}
func (LocalStore) Remove(key string) error {
	js.Global().Get("localStorage").Call("removeItem", key)
	return nil
}
