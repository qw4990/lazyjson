package lazyjson

import (
	"encoding/json"
)

type JSON interface {
	K(key string) JSON 				// get the value pointed by the key in the JSON map
	I(index int) JSON 				// get the value pointed by the index in the JSON array
	Int(defaultV int) int 			// convert this JSON value to int, return defaultV if failed
	Float(defaultV float64) float64 // convert this JSON value to float, return defaultV if failed
	String(defaultV string) string 	// convert this JSON value to string, return defaultV if failed
	Bool(defaultV bool) bool 		// convert this JSON value to bool, return defaultV if failed

	json.Marshaler
}

var (
	None = new(noneJSON)
)

type noneJSON struct{}

func (nj *noneJSON) K(key string) JSON {
	return nj
}

func (nj *noneJSON) I(index int) JSON {
	return nj
}

func (nj *noneJSON) Int(defaultV int) int {
	return defaultV
}

func (nj *noneJSON) Float(defaultV float64) float64 {
	return defaultV
}

func (nj *noneJSON) String(defaultV string) string {
	return defaultV
}

func (nj *noneJSON) Bool(defaultV bool) bool {
	return defaultV
}

func (nj *noneJSON) MarshalJSON() ([]byte, error) {
	return []byte(""), nil
}

type lazyJSON struct {
	a []interface{}
	m map[string]interface{}
	v interface{}
}

func (nj *lazyJSON) K(key string) JSON {
	if nj.m == nil {
		return None
	}
	return newJSON(nj.m[key])
}

func (nj *lazyJSON) I(index int) JSON {
	if index < 0 || index >= len(nj.a) {
		return None
	}
	return newJSON(nj.a[index])
}

func (nj *lazyJSON) Int(defaultV int) int {
	if v, ok := i2int(nj.v); ok {
		return v
	}
	if v, ok := i2float(nj.v); ok {
		return int(v)
	}
	return defaultV
}

func (nj *lazyJSON) Float(defaultV float64) float64 {
	if v, ok := i2float(nj.v); ok {
		return v
	}
	return defaultV
}

func (nj *lazyJSON) String(defaultV string) string {
	if v, ok := i2string(nj.v); ok {
		return v
	}
	return defaultV
}

func (nj *lazyJSON) Bool(defaultV bool) bool {
	if v, ok := i2bool(nj.v); ok {
		return v
	}
	return defaultV
}

func (nj *lazyJSON) MarshalJSON() ([]byte, error) {
	if nj.m != nil {
		return json.Marshal(nj.m)
	} else if nj.a != nil {
		return json.Marshal(nj.a)
	} else {
		return json.Marshal(nj.v)
	}
}

func NewJSON(data []byte) (JSON, error) {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return nil, err
	}

	return newJSON(i), nil
}

func newJSON(i interface{}) JSON {
	switch v := i.(type) {
	case map[string]interface{}:
		return &lazyJSON{m:v}
	case []interface{}:
		return &lazyJSON{a:v}
	default:
		return &lazyJSON{v:i}
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
