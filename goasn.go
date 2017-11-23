package goasn

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type ASNInfo struct {
	Descr string
}

type ASNReference struct {
	URL     string
	Data    map[uint64]ASNInfo
	Offline bool
}

func NewASN() *ASNReference {
	return &ASNReference{
		URL: "http://bgp.potaroo.net/cidr/autnums.html",
	}
}

func (a *ASNReference) getDataURL() (map[uint64]ASNInfo, error) {
	resp, err := http.Get(a.URL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	re := regexp.MustCompile(">AS(\\d+)\\s*</a>\\s*(.*)\\s*")

	scanner := bufio.NewScanner(resp.Body)
	result := make(map[uint64]ASNInfo)

	for scanner.Scan() {
		asn := re.FindStringSubmatch(scanner.Text())
		if len(asn) == 3 {
			if num, err := strconv.ParseUint(asn[1], 10, 64); err == nil {
				result[num] = ASNInfo{Descr: asn[2]}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *ASNReference) Init() {
	r, _ := a.getDataURL()
	a.Data = r

	fh, err := os.Create("./goasn.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := a.dump(fh); err != nil {
		log.Fatal(err)
	}

	fh.Close()
}

func (a *ASNReference) Get(asn uint64) ASNInfo {
	return a.Data[asn]
}

func (a *ASNReference) load(r io.Reader) error {
	zh, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	_ = zh

	return nil
}

func (a *ASNReference) dump(w io.Writer) error {
	var buf = new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(a.Data); err != nil {
		return err
	}

	zh := gzip.NewWriter(w)
	if _, err := zh.Write(buf.Bytes()); err != nil {
		return err
	}

	zh.Close()

	return nil
}
