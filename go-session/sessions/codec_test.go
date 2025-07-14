package sessions

import (
	"testing"
)

type testSession struct {
	ID    string
	Value map[interface{}]interface{}
}

func newTestSession() testSession {
	return testSession{
		ID:    "secret",
		Value: make(map[interface{}]interface{}),
	}
}

func TestID(t *testing.T) {
	cases := []struct {
		Name  string
		Value interface{}
	}{
		{"1", "1"},
		{"2", "2"},
	}

	codec := CodecFromKey([]byte("secret"))
	ts := newTestSession()

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			res, _ := codec.Encode(c.Name, c.Value)
			codec.Decode(c.Name, res, &ts.ID)

			if ts.ID != c.Value.(string) {
				t.Error(c.Name, ts.ID)
			}
		})
	}
}

func TestValues(t *testing.T) {
	cases := []struct {
		Name  string
		Value interface{}
	}{
		{"3", map[interface{}]interface{}{
			"3.1": 3.1,
			"3.2": "3.2",
			3.3:   "3.3",
			3.4:   3.4,
		}},
	}

	codec := CodecFromKey([]byte("secret"))
	ts := newTestSession()

	res, _ := codec.Encode(cases[0].Name, cases[0].Value)
	codec.Decode(cases[0].Name, res, &ts.Value)

	// t.Error(ts.Value["3.1"])
	// t.Error(ts.Value["3.2"])
	// t.Error(ts.Value[3.3])
	// t.Error(ts.Value[3.4])
}
