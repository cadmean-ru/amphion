package a

import "github.com/cadmean-ru/amphion/common/require"

// SiMap provides an abstraction over Go's map[string]interface{}.
type SiMap map[string]interface{}

func (m SiMap) GetString(key string) string {
	return require.String(m[key])
}

func (m SiMap) GetInt(key string) int {
	return require.Int(m[key])
}

func (m SiMap) GetInt32(key string) int32 {
	return require.Int32(m[key])
}

func (m SiMap) GetInt64(key string) int64 {
	return require.Int64(m[key])
}

func (m SiMap) GetFloat32(key string) float32 {
	return require.Float32(m[key])
}

func (m SiMap) GetFloat64(key string) float64 {
	return require.Float64(m[key])
}

func (m SiMap) GetBool(key string) bool {
	return require.Bool(m[key])
}
func (m SiMap) ContainsKey(key string) bool {
	_, ok := m[key]
	return ok
}