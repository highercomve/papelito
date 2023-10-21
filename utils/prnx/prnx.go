package prnx

import (
	"fmt"
	"strings"
)

// Prn string to define Prn
type Prn string

// PrnParseError string to define Prn parse Error
type PrnParseError string

func (s PrnParseError) Error() string {
	return string(s)
}

// PrnInfo Prn information
type PrnInfo struct {
	Domain   string
	Service  string
	Resource string
	ID       string
}

// PrnGetID make this a nice prn helper tool
func PrnGetID(prn string) string {
	idx := strings.Index(prn, "/")
	return prn[idx+1:]
}

// IDGetPrn get prn ID
func IDGetPrn(id string, serviceName string) string {
	return "prn:::" + serviceName + ":/" + id
}

// GetPrn get prn ID
func GetPrn(id string, serviceName string) Prn {
	return Prn("prn:::" + serviceName + ":/" + id)
}

// GetInfo get information
func (p *Prn) GetInfo() (*PrnInfo, error) {
	prn := string(*p)
	if !strings.HasPrefix(prn, "prn:") {
		errstr := fmt.Sprintf("ERROR: prn parse prn: prefix missing - %s", *p)
		return nil, PrnParseError(errstr)
	}

	rs := PrnInfo{}

	i := strings.Index(prn[4:], ":")
	if i > 0 {
		rs.Domain = prn[4 : 4+i]
	}

	j := strings.Index(prn[4+i+1:], ":")

	if j > 0 {
		rs.Service = prn[4+i+1 : 4+i+1+j]
	}

	if len(prn) > 4+i+1+j+1 {
		rs.Resource = prn[4+i+1+j+1:]
	}

	idx := strings.Index(prn, "/")
	rs.ID = prn[idx+1:]

	return &rs, nil
}

// Equals test if two PRN are equals
func (p *PrnInfo) Equals(c *PrnInfo) bool {
	return p.Domain == c.Domain &&
		p.Service == c.Service &&
		p.Resource == c.Resource
}
