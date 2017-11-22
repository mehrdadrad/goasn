package goasn

import (
	"bufio"
	"net/http"
	"regexp"
)

type ASNReference struct {
	URL     string
	Offline bool
}

func NewASN() *ASNReference {
	return &ASNReference{
		URL: "http://bgp.potaroo.net/cidr/autnums.html",
	}
}

func (a *ASNReference) getDataURL() (map[string]string, error) {
	resp, err := http.Get(a.URL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	re := regexp.MustCompile(">AS(\\d+)\\s*</a>\\s*(.*)\\s*")

	scanner := bufio.NewScanner(resp.Body)
	result := make(map[string]string)

	for scanner.Scan() {
		asn := re.FindStringSubmatch(scanner.Text())
		if len(asn) == 3 {
			result[asn[1]] = asn[2]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
