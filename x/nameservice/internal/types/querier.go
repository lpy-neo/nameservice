package types

import (
	"fmt"
	"strings"
)

// Query endpoints supported by the nameservice querier
const (
// TODO: Describe query parameters, update <action> with your query
// Query<Action>    = "<action>"
)

// QueryResResolve Queries Result Payload for a resolve query
type QueryResResolve struct {
	Value string `json:"value"`
}

// implement fmt.Stringer
func (r QueryResResolve) String() string {
	return r.Value
}

// QueryResNames Queries Result Payload for a names query
type QueryResNames []string

// implement fmt.Stringer
func (n QueryResNames) String() string {
	return strings.Join(n[:], "\n")
}

type QuerySaleStatus struct {
	SaleType SaleType `json:"sale_type,omitempty"`
	Price    string   `json:"price,omitempty"`
	BidPrice string   `json:"bid_price,omitempty"`
}

func (n QuerySaleStatus) String() string {
	return fmt.Sprintf(`saleType: %ld,
price: %s,
lastBidPrice: %s`, n.SaleType, n.Price, n.BidPrice)
}
