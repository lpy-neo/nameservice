package keeper

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/lpy-neo/nameservice/x/nameservice/internal/types"
)

// Keeper of the nameservice store
type Keeper struct {
	CoinKeeper types.BankKeeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace types.ParamSubspace
}

// NewKeeper creates a nameservice keeper
func NewKeeper(coinKeeper bank.Keeper, cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   key,
		cdc:        cdc,
		paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) set(ctx sdk.Context, key string, value interface{}) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set([]byte(key), bz)
}

func (k Keeper) delete(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(key))
}

func (k Keeper) SetWhois(ctx sdk.Context, name string, whois types.Whois) {
	if whois.Owner.Empty() {
		return
	}
	k.set(ctx, name, whois)
}

func (k Keeper) GetWhois(ctx sdk.Context, name string) types.Whois {
	store := ctx.KVStore(k.storeKey)
	if !k.IsNamePresent(ctx, name) {
		return types.NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois types.Whois

	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &whois)
	if err != nil {
		return types.NewWhois()
	}

	return whois
}

func (k Keeper) DeleteWhois(ctx sdk.Context, name string) {
	k.delete(ctx, name)
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhois(ctx, name).Owner.Empty()
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

// GetPrice - gets the current price of a name
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhois(ctx, name).Price
}

// GetSaleStaus - gets the current sale status of a name
func (k Keeper) GetSaleStaus(ctx sdk.Context, name string) types.SaleStatus {
	whois := k.GetWhois(ctx, name)
	return whois.SaleStatus
}

func (k Keeper) AddBid(ctx sdk.Context, name string, buyer sdk.AccAddress, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.SaleStatus.Bids = append(whois.SaleStatus.Bids, types.Bid{
		Buyer:       buyer,
		Price:       price,
		BlockHeight: ctx.BlockHeight(),
		Timestamp:   time.Now().Unix(),
	})
	k.SetWhois(ctx, name, whois)
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

// Check if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// SetSale - sets the current price of a name
func (k Keeper) SetSale(ctx sdk.Context, name string, saleType types.SaleType, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	whois.SaleStatus = types.SaleStatus{
		SaleType: saleType,
		Price:    price,
	}
	k.SetWhois(ctx, name, whois)
}

func (k Keeper) FinishAuctions(ctx sdk.Context, curBlockHeight int64) int {
	finished := 0

	for it := k.GetNamesIterator(ctx); it.Valid(); it.Next() {
		name := string(it.Key())
		whois := k.GetWhois(ctx, name)

		if whois.SaleStatus.SaleType != types.SaleTypeAuction ||
			len(whois.SaleStatus.Bids) == 0 {
			continue
		}

		lastBid := whois.SaleStatus.Bids[len(whois.SaleStatus.Bids)-1]
		if curBlockHeight-lastBid.BlockHeight > types.AuctionFinishBlockInterval {
			if err := k.finishOneAuction(ctx, name, lastBid.Buyer, lastBid.Price); err != nil {

			} else {
				finished++
			}
		}
	}

	return finished
}

func (k Keeper) finishOneAuction(ctx sdk.Context, name string, buyer sdk.AccAddress, price sdk.Coins) error {
	_, err := k.CoinKeeper.SubtractCoins(ctx, buyer, price) // If so, deduct the Bid amount from the sender
	if err != nil {
		return err
	}

	whois := k.GetWhois(ctx, name)
	whois.Owner = buyer
	whois.Price = price
	whois.SaleStatus = types.SaleStatus{
		SaleType: types.SaleTypeNotSale,
	}
	k.SetWhois(ctx, name, whois)

	if k.HasOwner(ctx, name) {
		err := k.CoinKeeper.SendCoins(ctx, buyer, k.GetOwner(ctx, name), price)
		if err != nil {
			return err
		}
	}

	return nil
}
