package example

import (
	"testing"

	tensorflow "github.com/ryszard/tfutils/proto/tensorflow/core/example"
)

func TestExample(t *testing.T) {
	_ = New(map[string]interface{}{
		"string":       "a string",
		"string list":  []string{"some strings"},
		"bytes":        []byte("bytes"),
		"bytes list":   [][]byte{[]byte("bytes")},
		"float32":      float32(10),
		"float32 list": []float32{11},
		"int64":        12,
		"int64 list":   []int64{13},
		"feature":      &tensorflow.Feature{&tensorflow.Feature_BytesList{&tensorflow.BytesList{[][]byte{[]byte{byte(1)}}}}},
	})
}
