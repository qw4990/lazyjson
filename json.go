package lazyjson

import (
	"encoding/json"
)

type JSON struct {
	a []interface{}
	m map[string]interface{}
	v interface{}
}

var empty = JSON{}

func (nj JSON) K(key string) JSON {
	if nj.empty() {
		return empty
	}

	return newJSON(nj.m[key])
}

func (nj JSON) I(index int) JSON {
	if nj.empty() || index < 0 || index >= len(nj.a) {
		return empty
	}
	return newJSON(nj.a[index])
}

func (nj JSON) Int(defaultV int) int {
	if v, ok := i2int(nj.v); ok {
		return v
	}
	if v, ok := i2float(nj.v); ok {
		return int(v)
	}
	return defaultV
}

func (nj JSON) Float(defaultV float64) float64 {
	if v, ok := i2float(nj.v); ok {
		return v
	}
	return defaultV
}

func (nj JSON) String(defaultV string) string {
	if v, ok := i2string(nj.v); ok {
		return v
	}
	return defaultV
}

func (nj JSON) Bool(defaultV bool) bool {
	if v, ok := i2bool(nj.v); ok {
		return v
	}
	return defaultV
}

func (nj JSON) MarshalJSON() ([]byte, error) {
	if nj.m != nil {
		return json.Marshal(nj.m)
	} else if nj.a != nil {
		return json.Marshal(nj.a)
	} else {
		return json.Marshal(nj.v)
	}
}

func (nj *JSON) UnmarshalJSON(data []byte) error {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch v := i.(type) {
	case map[string]interface{}:
		nj.m = v
	case []interface{}:
		nj.a = v
	default:
		nj.v = i
	}

	return nil
}

func (nj JSON) Size() int {
	if nj.empty() {
		return 0
	}
	if nj.m != nil {
		return len(nj.m)
	}
	if nj.a != nil {
		return len(nj.a)
	}
	return 1
}

func (nj JSON) empty() bool {
	return nj.a == nil && nj.m == nil && nj.v == nil
}

func newJSON(i interface{}) JSON {
	if i == nil {
		return empty
	}

	switch v := i.(type) {
	case map[string]interface{}:
		return JSON{m: v}
	case []interface{}:
		return JSON{a: v}
	default:
		return JSON{v: i}
	}
}

func i2int(i interface{}) (int, bool) {
	switch v := i.(type) {
	case int:
		return v, true
	case int8:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case uint:
		return int(v), true
	case uint8:
		return int(v), true
	case uint32:
		return int(v), true
	case uint64:
		return int(v), true
	}
	return 0, false
}

func i2bool(i interface{}) (bool, bool) {
	switch v := i.(type) {
	case bool:
		return v, true
	}
	return false, false
}

func i2float(i interface{}) (float64, bool) {
	switch v := i.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	}
	return 0, false
}

func i2string(i interface{}) (string, bool) {
	switch v := i.(type) {
	case string:
		return v, true
	}
	return "", false
}
