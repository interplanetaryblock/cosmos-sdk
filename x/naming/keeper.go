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
func NewKeeper(cdc *wire.Codec, coinKeeper bank.Keeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:        cdc,
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
	}
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyResolve(name))
	return string(bz)
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyResolve(name), []byte(value))
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyOwner(name))
	return bz != nil
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyOwner(name))
	return bz
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyOwner(name), owner)
}

// GetPrice - gets the current price of a name.  If price doesn't exist yet, set to 1steak.
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	if !k.HasOwner(ctx, name) {
		return sdk.Coins{sdk.NewInt64Coin("mycoin", 1)}
	}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyPrice(name))
	var price sdk.Coins
	k.cdc.MustUnmarshalBinary(bz, &price)
	return price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyPrice(name), k.cdc.MustMarshalBinary(price))
}

// KeyResolve return Key for getting the resolving value of a specific name from the store
func KeyResolve(name string) []byte {
	return []byte(fmt.Sprintf("%s%s", resolvePrefix, name))
}

// KeyOwner return Key for getting the owner of a specific name from the store
func KeyOwner(name string) []byte {
	return []byte(fmt.Sprintf("%s%s", ownerPrefix, name))
}

// KeyPrice return Key for getting the price of a specific name from the store
func KeyPrice(name string) []byte {
	return []byte(fmt.Sprintf("%s%s", pricePrefix, name))
}
