// A TFRecords file contains a sequence of strings with CRC
// hashes. Each record has the format
//
//     uint64 length
//     uint32 masked_crc32_of_length
//     byte   data[length]
//     uint32 masked_crc32_of_data
//
// and the records are concatenated together to produce the file. The
// CRC32s are described here, and the mask of a CRC is
//
//     masked_crc = ((crc >> 15) | (crc << 17)) + 0xa282ead8ul
//
// For more information, please refer to
// https://www.tensorflow.org/versions/master/api_docs/python/python_io.html#tfrecords-format-details.
package tfrecord

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

// maskDelta is a magic number taken from
// https://github.com/tensorflow/tensorflow/blob/754048a0453a04a761e112ae5d99c149eb9910dd/tensorflow/core/lib/hash/crc32c.h#L33.
const maskDelta uint32 = 0xa282ead8

// mask returns a masked representation of crc.
//
// Motivation: it is problematic to compute the CRC of a string that
// contains embedded CRCs.  Therefore we recommend that CRCs stored
// somewhere (e.g., in files) should be masked before being stored.
func mask(crc uint32) uint32 {
	return ((crc >> 15) | (crc << 17)) + maskDelta
}

// unmask returns the unmasked representation of crc. See the
// docstring of mask.
func unmask(masked uint32) uint32 {
	rot := masked - maskDelta
	return ((rot >> 17) | (rot << 15))
}

// uint64ToBytes returns x as bytes.
func uint64ToBytes(x uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, x)
	return b
}

var crc32Table = crc32.MakeTable(crc32.Castagnoli)

// crc32Hash returs the crc32 has expected by the C++ TensorFlow
// libraries.
func crc32Hash(data []byte) uint32 {
	return crc32.Checksum(data, crc32Table)
}

// WriteRecord writes the provided data as a Record to w.
func Write(w io.Writer, data []byte) error {

	var (
		length    = uint64(len(data))
		lengthCRC = mask(crc32Hash(uint64ToBytes(length)))
		dataCRC   = mask(crc32Hash(data))
	)

	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, lengthCRC); err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, dataCRC); err != nil {
		return err
	}

	return nil
}

// ReadRecord reads one record from r.
func Read(r io.Reader) (data []byte, err error) {
	var (
		length         uint64
		lengthChecksum uint32
		dataChecksum   uint32
	)

	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &lengthChecksum); err != nil {
		return nil, err
	}

	if actual := mask(crc32Hash(uint64ToBytes(length))); actual != lengthChecksum {
		return nil, errors.New("data length checksum doesn't match")
	}

	data = make([]byte, length)

	if _, err := r.Read(data); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &dataChecksum); err != nil {
		return nil, err
	}

	if actual := mask(crc32Hash(data)); actual != dataChecksum {
		return nil, errors.New("data checksum doesn't match")
	}

	return data, nil
}
