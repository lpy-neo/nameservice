package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MinNamePrice is Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

type SaleType = int

const (
	SaleTypeNotSale SaleType = iota
	SaleTypeNormal
	SaleTypeAuction
)

const (
	AuctionFinishBlockInterval = 100
)

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value      string         `json:"value"`
	Owner      sdk.AccAddress `json:"owner"`
	Price      sdk.Coins      `json:"price"`
	SaleStatus SaleStatus     `json:"saleStaus"`
}

type SaleStatus struct {
	SaleType SaleType  `json:"sale_type"`
	Price    sdk.Coins `json:"price"`
	Bids     []Bid     `json:"bids,omitempty"`
}

type Bid struct {
	Buyer       sdk.AccAddress `json:"buyer,omitempty"`
	Price       sdk.Coins      `json:"price"`
	BlockHeight int64          `json:"blockHeight"`
	Timestamp   int64          `json:"timestamp,omitempty"`
}

func IsSaleTypeValid(saleType SaleType) bool {
	return saleType == SaleTypeNotSale ||
		saleType == SaleTypeAuction ||
		saleType == SaleTypeNormal
}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
		SaleStatus: SaleStatus{
			SaleType: SaleTypeNormal,
		},
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}
