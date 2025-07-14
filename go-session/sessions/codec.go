package sessions

import (
	"encoding/json"
	"errors"
	"strings"
)

// Codec defines an struct{} to encode and decode cookie values.
type Codec struct {
	Key []byte
}

func CodecFromKey(key []byte) Codec {
	return Codec{Key: key}
}

/*
Encode(name, session.ID)
Decode(name, value, &session.ID)
name + "-" + session.ID

Encode(name, session.Values)
Decode(name, value, &session.Values)
name + "-" + json.Marshal(session.Value)
*/

func (c *Codec) Encode(name string, value interface{}) (string, error) {
	prefix := name + "-"

	if len(c.Key) == 0 {
		return prefix + "", errors.New("no key provided")
	}

	switch value.(type) {
	case string:
		return prefix + value.(string), nil
	case map[string]interface{}:
		v, _ := json.Marshal(value)
		return prefix + string(v), nil
	default:
		return prefix + "", nil
	}
}

func (c *Codec) Decode(name, value string, dst interface{}) error {
	if len(c.Key) == 0 {
		return errors.New("no key provided")
	}

	s := strings.Split(value, "-")
	if len(s) != 2 {
		return errors.New("wrong value")
	}
	if s[0] != name {
		return errors.New("decode failed")
	}

	switch dst.(type) {
	case *string:
		p, _ := dst.(*string)
		*p = s[1]
		return nil
	case *map[string]interface{}:
		return json.Unmarshal([]byte(s[1]), dst)
	default:
		return errors.New("wrong type")
	}
}
