package example

import (
	"fmt"

	tensorflow "github.com/ryszard/tfutils/proto/tensorflow/core/example"
)

func New(values map[string]interface{}) *tensorflow.Example {
	features := make(map[string]*tensorflow.Feature)
	for k, v := range values {
		switch t := v.(type) {
		case []byte:
			features[k] = Bytes(t)
		case [][]byte:
			features[k] = BytesList(t)
		case string:
			features[k] = Bytes([]byte(t))
		case []string:
			b := make([][]byte, len(t))
			for i, s := range t {
				b[i] = []byte(s)
			}
			features[k] = BytesList(b)
		case float32:
			features[k] = Float(t)
		case []float32:
			features[k] = FloatList(t)
		case int64:
			features[k] = Int64(t)
		case []int64:
			features[k] = Int64List(t)
		case int:
			features[k] = Int64(int64(t))
		case []int:
			ints := make([]int64, len(t))
			for i, ii := range t {
				ints[i] = int64(ii)
			}
			features[k] = Int64List(ints)
		default:
			panic(fmt.Sprintf("example: unsupported feature type %T: %q", t, t))
		}
	}
	return &tensorflow.Example{
		Features: &tensorflow.Features{features},
	}
}

func Bytes(value []byte) *tensorflow.Feature {
	values := [][]byte{value}
	return &tensorflow.Feature{&tensorflow.Feature_BytesList{&tensorflow.BytesList{values}}}
}

func BytesList(values [][]byte) *tensorflow.Feature {
	return &tensorflow.Feature{&tensorflow.Feature_BytesList{&tensorflow.BytesList{values}}}
}

func Float(value float32) *tensorflow.Feature {
	values := []float32{value}
	return &tensorflow.Feature{&tensorflow.Feature_FloatList{&tensorflow.FloatList{values}}}
}

func FloatList(values []float32) *tensorflow.Feature {
	return &tensorflow.Feature{&tensorflow.Feature_FloatList{&tensorflow.FloatList{values}}}
}

func Int64(value int64) *tensorflow.Feature {
	values := []int64{value}
	return &tensorflow.Feature{&tensorflow.Feature_Int64List{&tensorflow.Int64List{values}}}
}

func Int64List(values []int64) *tensorflow.Feature {
	return &tensorflow.Feature{&tensorflow.Feature_Int64List{&tensorflow.Int64List{values}}}
}
