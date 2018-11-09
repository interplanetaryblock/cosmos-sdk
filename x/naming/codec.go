package naming

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// RegisterCodec concrete types on wire codec
func RegisterCodec(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSetName{}, "naming/SetName", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "naming/BuyName", nil)
}
