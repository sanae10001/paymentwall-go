package paymentwall

import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"net/url"
	"sort"
	"strings"
)

const (
	baseUrl = "https://api.paymentwall.com/api"

	VC_CONTROLLER    = "ps"
	GOODS_CONTROLLER = "subscription"
	CART_CONTROLLER  = "cart"
)

var (
	ErrorOnlyOneProductAllowed = errors.New("Only one product is allowed when ApiType is API_GOODS")
)

func NewWidget(
	appKey string, apiType ApiType,
	secretKey, uid, widgetCode string,
	skipSignature bool) *Widget {
	w := &Widget{
		appKey:        appKey,
		secretKey:     secretKey,
		apiType:       apiType,
		code:          widgetCode,
		uid:           uid,
		skipSignature: skipSignature,
	}
	return w
}

type Widget struct {
	appKey        string
	secretKey     string
	apiType       ApiType
	code          string // widget code
	skipSignature bool
	uid           string

	products    []Product
	extraParams map[string]string
}

func (w *Widget) AppendProduct(products ...Product) error {
	if len(w.products) == 0 {
		w.products = make([]Product, 0, len(products))
	}
	if w.apiType == API_GOODS {
		if len(products) > 1 {
			return ErrorOnlyOneProductAllowed
		}
		if len(w.products) == 1 {
			return ErrorOnlyOneProductAllowed
		}
	}
	w.products = append(w.products, products...)
	return nil
}

func (w *Widget) SetExtraParams(m map[string]string) {
	for k, v := range m {
		w.SetExtraParam(k, v)
	}
}

func (w *Widget) SetExtraParam(k, v string) {
	if w.extraParams == nil {
		w.extraParams = make(map[string]string)
	}
	w.extraParams[k] = v
}

func (w *Widget) GetHtmlCode(attributes map[string]string) string {
	defaultAttributes := map[string]string{
		"frameborder": "0",
		"width":       "750",
		"height":      "800",
	}

	for k, v := range attributes {
		defaultAttributes[k] = v
	}

	attrsSlice := make([]string, 0, len(defaultAttributes))
	for attr, value := range defaultAttributes {
		attrsSlice = append(attrsSlice,
			fmt.Sprintf(`%s="%s"`, attr, value))
	}

	return fmt.Sprintf(`<iframe src="%s" %s></iframe>`,
		w.GetUrl(), strings.Join(attrsSlice, " "))
}

func (w *Widget) GetUrl() string {
	return baseUrl + "/" + w.buildController() + "?" + w.getParams().Encode()
}

func (w *Widget) getDefaultWidgetSignature() string {
	if w.apiType != API_CART {
		return DefaultSignVersion
	} else {
		return SignVersion2
	}
}

func (w *Widget) mergeSignVersion() string {
	signVersion := w.getDefaultWidgetSignature()
	if s, ok := w.extraParams["sign_version"]; ok {
		signVersion = s
	}
	return signVersion
}

func (w *Widget) buildController() string {
	if w.apiType == API_VC {
		return VC_CONTROLLER
	} else if w.apiType == API_GOODS {
		return GOODS_CONTROLLER
	} else {
		return CART_CONTROLLER
	}
}

func (w *Widget) getParams() url.Values {
	params := url.Values{}
	params.Set("key", w.appKey)
	params.Set("widget", w.code)
	params.Set("uid", w.uid)

	if len(w.products) > 0 {
		if w.apiType == API_GOODS {
			product := w.products[0]
			params.Set("amount", product.DisplayAmount())
			params.Set("currencyCode", product.Currency)
			params.Set("ag_name", product.Name)
			params.Set("ag_external_id", product.Identity)
			params.Set("ag_type", string(product.Type))
			if product.Type == ProductTypeSubscription {
				params.Set("ag_period_length", product.DisplayPeriodLength())
				params.Set("ag_period_type", string(product.PeriodType))
				if product.Recurring {
					params.Set("ag_recurring", "1")
					// TODO: add trial product
				}
			}
		} else if w.apiType == API_CART {
			for i, product := range w.products {
				params.Set(fmt.Sprintf("external_ids[%d]", i), product.Identity)
				if product.Amount > 0 {
					params.Set(fmt.Sprintf("prices[%d]", i), product.DisplayAmount())
				}
				if product.Currency != "" {
					params.Set(fmt.Sprintf("currencies[%d]", i), product.Currency)
				}
			}
		}
	}
	if !w.skipSignature {
		signVersion := w.mergeSignVersion()
		params.Set("sign_version", string(signVersion))
		params.Set("sign", w.calculateSignature(params, signVersion))
	}
	return params
}

func (w *Widget) calculateSignature(params url.Values, signVersion string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	var h hash.Hash
	if signVersion == SignVersion3 {
		h = sha256.New()
	} else {
		h = md5.New()
	}

	sort.Strings(keys)
	for _, k := range keys {
		h.Write([]byte(fmt.Sprintf("%s=%s", k, params.Get(k))))
	}
	h.Write([]byte(w.secretKey))
	return string(h.Sum(nil))
}
