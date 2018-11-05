# Paymentwall Go Library
Unofficial Paymentwall Go Library

### TODO
- [ ] Pingback Processing
	* [x] Digital Goods
	* [x] Virtual Currency
	* [ ] Cart
- [x] Widget Call
- [ ] Coverage test

### Reference
- [Paymentwall](https://docs.paymentwall.com/)
- [Pingback](https://docs.paymentwall.com/reference/pingback-home)

# Code Samples

## Digital Goods API

#### Pingback Processing
```go
import (
	"fmt"
	
	"github.com/sanae10001/paymentwall-go"
)

pingback := paymentwall.NewPingback(values, ip, paymentwall.API_GOODS, PaymentwallSecretKey)

if pingback.Validate(false) {
	if pingback.IsDeliverable() { 
		// deliver the product
	} else if pingback.IsCancelable() {
		// withdraw the product
	}
} else {
	fmt.Println(pingback.GetError())
}

```
