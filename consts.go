package paymentwall

type APIType uint

const (
	API_VC    APIType = 1 // Virtual Currency API
	API_GOODS APIType = 2 // Digital Goods API
)

type PingbackType string

const (
	PingbackTypeGoodwill PingbackType = "1"
	PingbackTypeNegative PingbackType = "2"
)

const (
	SignVersion2 = "2" // md5
	SignVersion3 = "3" // sha256
)
