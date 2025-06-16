package md5

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"testing"
)

type sample struct {
	input, res string
}

var samples map[string]sample = map[string]sample{
	"sample01": {input: "haoran"},
}

func TestBuiltinMd5(t *testing.T) {
	fill()
	for name, v := range samples {
		t.Run(name, func(t *testing.T) {
			want := fmt.Sprintf("%x", md5.Sum([]byte(v.input)))
			msg := intensify([]byte(v.input))
			v.res = MainLoop(msg)
			if !reflect.DeepEqual(want, v.res) {
				t.Error("want = ", want, "res = ", v.res)
			}
		})
	}
}
