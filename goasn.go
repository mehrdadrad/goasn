package goasn

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"errors"
	"io"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strconv"
)

// ASNInfo represents ASN description
type ASNInfo struct {
	Descr string
}

// ASNReference represents ASN source and data
type ASNReference struct {
	URL  string
	Path string
	Data map[uint64]ASNInfo
}

// NewASN create new ASN instance
func NewASN() *ASNReference {
	user, _ := user.Current()
	path := user.HomeDir

	return &ASNReference{
		URL:  "http://bgp.potaroo.net/cidr/autnums.html",
		Path: path + "/.",
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

	err = scanner.Err()

	return result, err
}

// Init loads data from origin or database
func (a *ASNReference) Init() error {

	if err := a.loadFromDB(); err == nil {
		return nil
	}

	err := a.loadFromOrigin()

	return err
}

func (a *ASNReference) loadFromOrigin() error {
	r, err := a.getDataURL()
	if err != nil {
		return err
	}

	a.Data = r

	fh, err := os.Create(a.Path + "goasn.db")
	if err != nil {
		return err
	}

	defer fh.Close()

	err = a.dump(fh)

	return err
}

func (a *ASNReference) loadFromDB() error {
	fh, err := os.Open(a.Path + "goasn.db")
	if err != nil {
		return err
	}

	defer fh.Close()

	r, err := a.load(fh)
	if err != nil {
		return err
	}

	a.Data = r

	return nil
}

// Get returns ASN description
func (a *ASNReference) Get(asn uint64) (ASNInfo, error) {
	d, ok := a.Data[asn]
	if !ok {
		return ASNInfo{}, errors.New("ASN not found")
	}

	return d, nil
}

func (a *ASNReference) load(r io.Reader) (map[uint64]ASNInfo, error) {
	buf := new(bytes.Buffer)
	zh, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	io.Copy(buf, zh)

	asn := make(map[uint64]ASNInfo)
	dec := gob.NewDecoder(buf)

	if err = dec.Decode(&asn); err != nil {
		return nil, err
	}

	return asn, nil
}

func (a *ASNReference) dump(w io.Writer) error {
	buf := new(bytes.Buffer)
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
