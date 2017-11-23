package goasn

import (
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
	} else if asn.descr != "ORACLE-ASNBLOCK-ASN - Oracle Corporation, US" {
		t.Errorf("Expect to have AS792 description but it has %s", asn.descr)
	}
}
