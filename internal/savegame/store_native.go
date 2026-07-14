//go:build !js

package savegame

type MemoryStore struct{ values map[string][]byte }

func NewStore() *MemoryStore { return &MemoryStore{values: map[string][]byte{}} }
func (s *MemoryStore) Load(key string) ([]byte, bool, error) {
	v, ok := s.values[key]
	return append([]byte(nil), v...), ok, nil
}
func (s *MemoryStore) Save(key string, value []byte) error {
	s.values[key] = append([]byte(nil), value...)
	return nil
}
func (s *MemoryStore) Remove(key string) error { delete(s.values, key); return nil }
