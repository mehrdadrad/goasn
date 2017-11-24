package goasn

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"net/http"
	"net/http/httptest"
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
