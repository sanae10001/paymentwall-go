package paymentwall

import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"net/url"
	"sort"
)

var ipsWhitelist = []string{
	"174.36.92.186",
	"174.36.96.66",
	"174.36.92.187",
	"174.36.92.192",
	"174.37.14.28",
}

// https://docs.paymentwall.com/reference/pingback-home
func NewPingback(values url.Values, ip string, apiType APIType, secretKey string) *Pingback {
	p := Pingback{
		m:           make(map[string]string, len(values)),
		keys:        make([]string, 0, len(values)),
		signVersion: SignVersion2,
		isTest:      false,
		ip:          ip,
		apiType:     apiType,
		secretKey:   secretKey,
		errors:      make([]error, 0, 2),
	}
	for k := range values {
		v := values.Get(k)
		if k == "sign_version" {
			p.signVersion = v
		} else if k == "is_test" && v == "1" {
			p.isTest = true
		} else {
			p.set(k, v)
		}
	}
	return &p
}

type Pingback struct {
	keys []string
	m    map[string]string

	signVersion string
	isTest      bool

	ip        string
	apiType   APIType
	secretKey string

	errors []error
}

func (p *Pingback) set(key, value string) {
	p.m[key] = value
	p.keys = append(p.keys, key)
}

func (p *Pingback) appendToError(errMsg string) {
	p.errors = append(p.errors, errors.New(errMsg))
}

func (p *Pingback) GetError() error {
	return p.errors[0]
}

func (p *Pingback) GetErrors() []error {
	return p.errors
}

func (p *Pingback) Validate(skipIPCheck bool) bool {
	var validated bool
	if p.IsParametersValid() {
		if p.IsIPValid() || skipIPCheck {
			if p.IsSignatureValid() {
				validated = true
			} else {
				p.appendToError("Wrong signature")
			}
		} else {
			p.appendToError("IP address is not whitelisted")
		}
	}
	return validated
}

func (p *Pingback) IsParametersValid() bool {
	var requiredParams []string
	if p.apiType == API_VC {
		requiredParams = []string{"uid", "type", "ref", "sign", "sign_version", "currency"}
	} else {
		requiredParams = []string{"uid", "type", "ref", "sign", "sign_version", "goodsid"}
	}

	for _, k := range requiredParams {
		if _, ok := p.m[k]; !ok {
			p.appendToError(fmt.Sprintf("Parameter %s is missing.", k))
			return false
		}
	}
	return true
}

func (p *Pingback) IsIPValid() bool {
	for _, i := range ipsWhitelist {
		if p.ip == i {
			return true
		}
	}
	return false
}

func (p *Pingback) IsSignatureValid() bool {
	var h hash.Hash
	if p.signVersion == SignVersion2 {
		h = md5.New()
	} else {
		h = sha256.New()
	}

	sort.Strings(p.keys)
	for _, k := range p.keys {
		if k == "sign" {
			continue
		}
		h.Write([]byte(fmt.Sprintf(`%s=%s`, k, p.m[k])))
	}
	h.Write([]byte(p.secretKey))

	return string(h.Sum(nil)) == p.m["sign"]
}
