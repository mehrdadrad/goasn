package goasn

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetDataURL(t *testing.T) {
	asn := NewASN()

	sampleData := ">AS792  </a> ORACLE-ASNBLOCK-ASN - Oracle Corporation, US"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, sampleData)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	asn.URL = ts.URL
	res, _ := asn.getDataURL()

	if asn, ok := res[792]; !ok {
		t.Error("Expect to have AS792 but not exit")
	} else if asn.Descr != "ORACLE-ASNBLOCK-ASN - Oracle Corporation, US" {
		t.Errorf("Expect to have AS792 description but it has %s", asn.Descr)
	}
}

func TestLoad(t *testing.T) {
	asn := NewASN()
	data := map[uint64]ASNInfo{15133: ASNInfo{Descr: "EdgeCast"}}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		t.Fatal(err)
	}

	w := new(bytes.Buffer)
	zh := gzip.NewWriter(w)
	if _, err := zh.Write(buf.Bytes()); err != nil {
		t.Fatal(err)
	}

	zh.Close()

	res, err := asn.load(w)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := res[15133]; !ok {
		t.Error("expected to have ASN 15133, but not exist!")
	}
}

func TestDump(t *testing.T) {
	asn := NewASN()
	asn.Data = map[uint64]ASNInfo{15133: ASNInfo{Descr: "EdgeCast"}}
	gzipData := []byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xe2, 0xff,
		0xdf, 0xcc, 0xc2, 0xc8, 0xf4, 0xbf, 0x85, 0x81, 0x91, 0x8d, 0xf1,
		0x7f, 0x13, 0x3, 0x83, 0xd8, 0xff, 0x46, 0x66, 0x46, 0xa6, 0xff,
		0x4d, 0xc, 0x8c, 0x8c, 0x8c, 0xac, 0x2e, 0xa9, 0xc5, 0xc9, 0x45,
		0x8c, 0x3c, 0xc, 0xc, 0xc, 0x42, 0x20, 0x15, 0xff, 0xac, 0x65,
		0x19, 0x39, 0x5c, 0x53, 0xd2, 0x53, 0x9d, 0x13, 0x8b, 0x4b, 0x18,
		0x0, 0x1, 0x0, 0x0, 0xff, 0xff, 0x95, 0xcc, 0x8e, 0x16, 0x3a, 0x0, 0x0, 0x0}

	buf := new(bytes.Buffer)
	asn.dump(buf)

	if !reflect.DeepEqual(buf.Bytes(), gzipData) {
		t.Error("dump unexpected error")
	}
}
