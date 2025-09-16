package global

const DefineBTC = "BTC"
const DefineBNB = "BNB"
const DefineUSDT = "USDT"
const DefineUSDC = "USDC"

const SideType_BUY string = "BUY"
const SideType_SELL string = "SELL"

type TradeMode int

const (
	BUY_SELL_SELL TradeMode = iota
	BUY_SELL_BUY
)
