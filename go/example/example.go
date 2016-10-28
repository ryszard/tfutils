// package example gives you some syntactic sugar to build
// tensorflow.Example protobufs in a more intuitive way.
package example

import (
	"fmt"

	tensorflow "github.com/ryszard/tfutils/proto/tensorflow/core/example"
)

// New builds and returns am Example protobuf based on the contents of
// features. The values of features will be massaged to work: single
// values will be turned into one element slices. Supported types:
// byte, string (converted to []byte), float32, int64, int (converted
// to int64), slices of supported types, and *tensorflow.Feature (this
// is an escape hatch).
//
// This function will panic if you pass it an unsupported type.
func New(features map[string]interface{}) *tensorflow.Example {
	result := make(map[string]*tensorflow.Feature)
	for k, v := range features {
		switch t := v.(type) {
		case []byte:
			result[k] = toBytes(t)
		case [][]byte:
			result[k] = toBytesList(t)
		case string:
			result[k] = toBytes([]byte(t))
		case []string:
			b := make([][]byte, len(t))
			for i, s := range t {
				b[i] = []byte(s)
			}
			result[k] = toBytesList(b)
		case float32:
			result[k] = toFloat(t)
		case []float32:
			result[k] = toFloatList(t)
		case int64:
			result[k] = toInt64(t)
		case []int64:
			result[k] = toInt64List(t)
		case int:
			result[k] = toInt64(int64(t))
		case []int:
			ints := make([]int64, len(t))
			for i, ii := range t {
				ints[i] = int64(ii)
			}
			result[k] = toInt64List(ints)
		case *tensorflow.Feature:
			result[k] = t
		default:
			panic(fmt.Sprintf("example: unsupported feature type %T: %q", t, t))
		}
	}
	return &tensorflow.Example{
		Features: &tensorflow.Features{result},
	}
}

func toBytes(value []byte) *tensorflow.Feature {
	values := [][]byte{value}
	return &tensorflow.Feature{&tensorflow.Feature_BytesList{&tensorflow.BytesList{values}}}
}

func toBytesList(values [][]byte) *tensorflow.Feature {
	return &tensorflow.Feature{&tensorflow.Feature_BytesList{&tensorflow.BytesList{values}}}
}

func toFloat(value float32) *tensorflow.Feature {
	values := []float32{value}
	return &tensorflow.Feature{&tensorflow.Feature_FloatList{&tensorflow.FloatList{values}}}
}

func toFloatList(values []float32) *tensorflow.Feature {
	return &tensorflow.Feature{&tensorflow.Feature_FloatList{&tensorflow.FloatList{values}}}
}

func toInt64(value int64) *tensorflow.Feature {
	values := []int64{value}
	return &tensorflow.Feature{&tensorflow.Feature_Int64List{&tensorflow.Int64List{values}}}
}

func toInt64List(values []int64) *tensorflow.Feature {
	return &tensorflow.Feature{&tensorflow.Feature_Int64List{&tensorflow.Int64List{values}}}
}
