package paymentwall

type ApiType uint

const (
	API_VC    ApiType = 1 // Virtual Currency API
	API_GOODS ApiType = 2 // Digital Goods API
	API_CART  ApiType = 3 // Cart API
)

type PingbackType string

const (
	PingbackTypeRegular  PingbackType = "0" // When product is purchased. Please check whether the reference ID is unique and deliver the goods.
	PingbackTypeGoodwill PingbackType = "1"
	PingbackTypeNegative PingbackType = "2" // When user issued chargeback or refund. Please take the delivered goods out from user’s account.

	PingbackTypeRiskUnderReview      PingbackType = "200" // Pending status. In case a payment is currently under risk review by Paymentwall. Please do not deliver the goods yet.
	PingbackTypeRiskReviewedAccepted PingbackType = "201" // Review is done and payment is accepted. Please check whether the reference ID is unique and deliver the goods.
	PingbackTypeRiskReviewedDeclined PingbackType = "202" // Review is done and payment is declined. Please do not deliver the goods since the user will get his money back.

	PingbackTypeRiskAuthorizationVoided PingbackType = "203" // Authorization has been voided due to no capture request received on time.

	PingbackTypeSubscriptionCancelled     PingbackType = "12" // When user cancels subscription plan. Sent immediately upon cancellation, e.g. in the middle of current premium month.
	PingbackTypeSubscriptionExpired       PingbackType = "13" // When subscription expired.
	PingbackTypeSubscriptionPaymentFailed PingbackType = "14" // When renewal subscription payment failed. Subscription stopped due to failing payments e.g. due to insufficient funds.
)

const (
	// https://docs.paymentwall.com/reference/signature-calculation
	DefaultSignVersion = "3" // sha256
	SignVersion2       = "2" // md5
	SignVersion3       = "3" // sha256
)

const (
	// https://docs.paymentwall.com/#chargeback-pingback
	PingbackChargebackReason1  = "1"  // Chargeback
	PingbackChargebackReason2  = "2"  // Credit Card fraud. Recommendation: Ban User
	PingbackChargebackReason3  = "3"  // Other fraud. Recommendation: Ban User
	PingbackChargebackReason4  = "4"  // Bad data entry
	PingbackChargebackReason5  = "5"  // Fake / Proxy user
	PingbackChargebackReason6  = "6"  // Rejected by advertiser
	PingbackChargebackReason7  = "7"  // Duplicate conversions
	PingbackChargebackReason8  = "8"  // Goodwill credit taken back
	PingbackChargebackReason9  = "9"  // Canceled order, e.g. refund
	PingbackChargebackReason10 = "10" // Partially reversed transaction
)
