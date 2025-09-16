package config

import "errors"

const TestConfigKey = "strategy_config"

type Config struct {

	// 交易对相关设置
	PriceFilterRateThreshold float64  `mapstructure:"price_filter_rate_threshold"` // 价格波动过滤阈值（例如 0.01 表示1%）
	OnlyBaseWhiteList        bool     `mapstructure:"only_base_white_list"`        // 是否仅启用基础币种白名单
	BaseWhiteList            []string `mapstructure:"base_white_list"`             // 基础币种白名单（当 only_base_white_list=true 时生效）
	BaseBlackList            []string `mapstructure:"base_black_list"`             // 基础币种黑名单

	SymbolZeroFeeList []string `mapstructure:"symbol_zero_fee_list"` // 零手续费交易对列表（如 BTCUSDT）
	SymbolWhiteList   []string `mapstructure:"symbol_white_list"`    // 交易对白名单（优先级高于黑名单）
	// SymbolBlackList   []string        `mapstructure:"symbol_black_list"`    // 交易对黑名单（禁止交易的交易对）
	AccountList []AccountConfig `mapstructure:"account_list"` // 子账户配置列表

	// 行情数据相关设置
	MultiIpMode    bool `mapstructure:"multi_ip_mode"`     // 是否使用多IP模式
	EnableHighTick bool `mapstructure:"enable_high_tick"`  // 是否启用高速行情模式
	MarketUseMaxIP int  `mapstructure:"market_use_max_ip"` // 订阅行情最多使用的IP数量, 不配置的情况下默认使用所有IP

	// 信号计算相关设置
	S1UseHighestRate    bool    `mapstructure:"s1_use_highest_return_rate"` // 是否采用最高套利组合计算 S1
	SymbolFeeRate       float64 `mapstructure:"symbol_fee_rate"`            // symbol手续费率
	S1TriggerThreshold  float64 `mapstructure:"s1_trigger_return_rate"`     // 触发 S1 计算的最小回报率
	ExpectMinReturnRate float64 `mapstructure:"expect_min_return_rate"`     // 预期最小利润率
	S2S1MinVolRate      float64 `mapstructure:"s2s1_min_vol_rate"`          // S2/S1 盘口量最小比例值

	// 定时任务相关设置
	Interval IntervalConf `mapstructure:"interval"` // 定时任务配置
}

func (c *Config) CheckValidity() error {
	if len(c.AccountList) == 0 {
		return errors.New("account_list is empty") // 账户不能为0
	}
	return nil
}

type AccountConfig struct {
	Ed25519ApiKey     string  `mapstructure:"ed25519_api_key"`     // API 密钥
	Ed25519PrivateKey string  `mapstructure:"ed25519_private_key"` // 私钥
	HoldCoin          string  `mapstructure:"hold_coin"`           // 持仓币种（如 USDT）
	NickName          string  `mapstructure:"nick_name"`           // 账户昵称
	StopLossThreshold float64 `mapstructure:"stop_loss_threshold"` // 止损阈值（单位：百分比，例如 0.03 表示3%）
}

type IntervalConf struct {
	SyncExchangeInfo int `mapstructure:"sync_exchange_info"` //同步交易规范(s)
	//CheckAccountBalance int   `mapstructure:"check_account_balance"` //检查账户余额(s)
	CheckCommission int   `mapstructure:"check_commission"` //检查佣金费率(s)
	CheckSlowestIp  int   `mapstructure:"check_slowest_ip"` //多ip模式下检查最慢IP(s)
	Rebalance       int   `mapstructure:"rebalance"`        //账户资金平衡(s)
	BuyBNB          int64 `mapstructure:"buy_bnb"`          //unit: ms, interval of buying bnb
	S2Sleep         int64 `mapstructure:"s2_sleep"`         //s2 symbol 锁定休眠时间 unit: us
}

var _ Configer = (*Config)(nil)
