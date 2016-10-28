package tfrecord

import (
	"bytes"
	"testing"
)

func TestMaskedCRC(t *testing.T) {
	if unmask(mask(7)) != 7 {
		t.Error("unmask(mask(7)) != 7")
	}
}

func TestReadWrite(t *testing.T) {
	rw := new(bytes.Buffer)

	data := []byte("Ala ma kota")

	if err := Write(rw, data); err != nil {
		t.Fatalf("Write(rw, %q): %v", data, err)
	}

	read, err := Read(rw)

	if err != nil {
		t.Fatalf("Read(rw): %v", err)
	}

	if !bytes.Equal(read, data) {
		t.Errorf("%q != %q", read, data)
	}
}

func TestReadCorruptedData(t *testing.T) {
	w := new(bytes.Buffer)

	data := []byte("Ala ma kota")

	if err := Write(w, data); err != nil {
		t.Fatalf("Write(rw, %q): %v", data, err)
	}
	original := w.Bytes()
	t.Logf("original: %q", original)
	corrupted := bytes.Replace(original, data, []byte("atok am alA"), 1)

	t.Logf("corrupted: %q", corrupted)
	r := bytes.NewBuffer(corrupted)
	if _, err := Read(r); err == nil {
		t.Errorf("Read didn't return an error on corrupted data")
	}
}

func TestReadCorruptedLength(t *testing.T) {
	w := new(bytes.Buffer)

	data := []byte("Ala ma kota")

	if err := Write(w, data); err != nil {
		t.Fatalf("Write(rw, %q): %v", data, err)
	}
	serialized := w.Bytes()
	t.Logf("original: %q", serialized)

	for i, b := range uint64ToBytes(uint64(len("Ala ma"))) {
		serialized[i] = b
	}

	t.Logf("corrupted: %q", serialized)
	r := bytes.NewBuffer(serialized)
	if _, err := Read(r); err == nil {
		t.Errorf("Read didn't return an error on corrupted length %v", err)
	}
}
