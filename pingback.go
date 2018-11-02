package paymentwall

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"net"
	"net/url"
	"sort"
)

// The whitelisted start and end range of which Paymentwall callbacks are permissible to come from.
const (
	ipWhitelistStart = "216.127.71.0"   // Start
	ipWhitelistEnd   = "216.127.71.255" // End
)

// https://docs.paymentwall.com/reference/pingback-home
func NewPingback(
	values url.Values,
	ip string, apiType ApiType, secretKey string) *Pingback {
	p := Pingback{
		m:           make(map[string]string, len(values)),
		keys:        make([]string, 0, len(values)),
		signVersion: DefaultSignVersion,
		IsTest:      false,
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
			p.IsTest = true
		}
		p.set(k, v)
	}
	return &p
}

type Pingback struct {
	keys []string
	m    map[string]string

	signVersion string
	IsTest      bool

	ip        string
	apiType   ApiType
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
		requiredParams = []string{"uid", "type", "ref", "sig", "sign_version", "currency"}
	} else if p.apiType == API_GOODS {
		requiredParams = []string{"uid", "type", "ref", "sig", "sign_version", "goodsid"}
	}

	for _, k := range requiredParams {
		if _, ok := p.m[k]; !ok {
			p.appendToError(fmt.Sprintf("Parameter %s is missing.", k))
			return false
		}
	}
	return true
}

// IsIPValid checks to ensure the IP from which the pingback was received is within the allowed range of whitelisted
// IPs provided by Paymentwall.
// @see: https://docs.paymentwall.com/reference/pingback-new-ip-subnet
func (p *Pingback) IsIPValid() bool {
	whitelistStart := net.ParseIP(ipWhitelistStart)
	whitelistEnd := net.ParseIP(ipWhitelistEnd)
	reqIP := net.ParseIP(p.ip)

	// Ensure IP is IPv4 only.
	if reqIP.To4() == nil {
		return false
	}

	// Check if IP is within range of whitelisted IPs.
	if bytes.Compare(reqIP, whitelistStart) >= 0 && bytes.Compare(reqIP, whitelistEnd) <= 0 {
		return true
	}

	return false
}

func (p *Pingback) IsSignatureValid() bool {
	var h hash.Hash
	if p.signVersion == SignVersion3 {
		h = sha256.New()
	} else {
		h = md5.New()
	}

	sort.Strings(p.keys)
	baseString := ""
	for _, k := range p.keys {
		if k == "sig" {
			continue
		}
		baseString += fmt.Sprintf(`%s=%s`, k, p.m[k])
	}
	baseString += p.secretKey
	h.Write([]byte(baseString))
	return hex.EncodeToString(h.Sum(nil)) == p.m["sig"]
}

func (p *Pingback) Get(key string) string {
	return p.m[key]
}

func (p *Pingback) GetType() PingbackType {
	return PingbackType(p.Get("type"))
}

// Unique identifier of the user. The value of parameter “uid” from Widget call.
func (p *Pingback) GetUID() string {
	return p.Get("uid")
}

func (p *Pingback) GetVCAmount() string {
	return p.Get("currency")
}

func (p *Pingback) GetProductID() string {
	return p.Get("goodsid")
}

func (p *Pingback) GetProductPeriod() (length string, period string) {
	return p.Get("slength"), p.Get("speriod")
}

func (p *Pingback) GetReferenceID() string {
	return p.Get("ref")
}

func (p *Pingback) IsDeliverable() bool {
	type_ := p.GetType()
	return type_ == PingbackTypeRegular ||
		type_ == PingbackTypeGoodwill ||
		type_ == PingbackTypeRiskReviewedAccepted
}

func (p *Pingback) IsCancelable() bool {
	type_ := p.GetType()
	return type_ == PingbackTypeNegative ||
		type_ == PingbackTypeRiskReviewedDeclined
}

func (p *Pingback) IsUnderReview() bool {
	return p.GetType() == PingbackTypeRiskUnderReview
}
