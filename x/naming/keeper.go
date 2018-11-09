package naming

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper - handlers sets/gets of custom variables for your module
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.

	cdc *wire.Codec // The wire codec for binary encoding/decoding.
}

var (
	resolvePrefix = "resolve_"
	ownerPrefix   = "owner_"
	pricePrefix   = "price_"
)

// NewKeeper creates new instances of the naming Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *wire.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	store := ctx.KVStore(k.storeKey)
	key := fmt.Sprintf("%s%s", resolvePrefix, name)
	bz := store.Get([]byte(key))
	return string(bz)
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	store := ctx.KVStore(k.storeKey)
	key := fmt.Sprintf("%s%s", resolvePrefix, name)
	store.Set([]byte(key), []byte(value))
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	key := fmt.Sprintf("%s%s", ownerPrefix, name)
	bz := store.Get([]byte(key))
	return bz != nil
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	key := fmt.Sprintf("%s%s", ownerPrefix, name)
	bz := store.Get([]byte(key))
	return bz
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := fmt.Sprintf("%s%s", ownerPrefix, name)
	store.Set([]byte(key), owner)
}

// GetPrice - gets the current price of a name.  If price doesn't exist yet, set to 1steak.
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	if !k.HasOwner(ctx, name) {
		return sdk.Coins{sdk.NewInt64Coin("mycoin", 1)}
	}
	store := ctx.KVStore(k.storeKey)
	key := fmt.Sprintf("%s%s", pricePrefix, name)
	bz := store.Get([]byte(key))
	var price sdk.Coins
	k.cdc.MustUnmarshalBinary(bz, &price)
	return price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	key := fmt.Sprintf("%s%s", pricePrefix, name)
	store.Set([]byte(key), k.cdc.MustMarshalBinary(price))
}
