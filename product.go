package paymentwall

import "strconv"

type ProductType string

const (
	ProductTypeSubscription ProductType = "subscription"
	ProductTypeFixed        ProductType = "fixed"
)

type PeriodType string

const (
	PeriodTypeDay   PeriodType = "day"
	PeriodTypeWeek  PeriodType = "week"
	PeriodTypeMonth PeriodType = "month"
	PeriodTypeYear  PeriodType = "year"
)

func NewProduct(
	name, id string,
	amount float64, currency string,
	pType ProductType) *Product {
	p := &Product{
		Name:     name,
		Identity: id,
		Amount:   amount,
		Currency: currency,
		Type:     pType,
	}
	return p
}

type Product struct {
	Name     string
	Identity string

	Amount   float64
	Currency string

	Type ProductType

	Recurring    bool
	PeriodLength uint
	PeriodType   PeriodType
}

func (p *Product) SetSubscription(
	periodLength uint,
	periodType PeriodType,
	isRecurring bool) {
	p.Recurring = isRecurring
	p.PeriodLength = periodLength
	p.PeriodType = periodType
}

func (p *Product) DisplayAmount() string {
	return strconv.FormatFloat(p.Amount, 'f', -1, 64)
}

func (p *Product) DisplayPeriodLength() string {
	return strconv.FormatUint(uint64(p.PeriodLength), 10)
}
