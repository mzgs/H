package H

import "strings"

const (
	BTCTURK = "BTCTURK"
	BINANCE = "BINANCE"
	FTX     = "FTX"
)

var CoinIcons = map[string]string{
	"ADA":   "₳",
	"BCH":   "Ƀ",
	"BSV":   "Ɓ",
	"BTC":   "₿",
	"DAI":   "◈",
	"DOGE":  "Ð",
	"EOS":   "ε",
	"ETC":   "ξ",
	"ETH":   "Ξ",
	"LTC":   "Ł",
	"MKR":   "Μ",
	"REP":   "Ɍ",
	"STEEM": "ȿ",
	"USDT":  "₮",
	"XMR":   "ɱ",
	"XRP":   "Ʀ",
	"XTZ":   "ꜩ",
	"ZEC":   "ⓩ",
}

func GetCoinIconOrName(symbol string) string {

	name := Mstr(symbol)
	name.Remove("/USD", "USDT", "USD", "-PERP")

	if val, ok := CoinIcons[name.String()]; ok {
		name.Set(val)
	}

	return name.String()
}

func GetCurrency(symbol string) string {
	target := GetTargetCode(symbol)
	var currency string
	switch target {
	case "USDT", "USD":
		currency = "$"
	case "TRY":
		currency = "₺"
	case "BTC":
		currency = "₿"
	}
	return currency
}

func GetSymbolUrl(ex, symbol string) string {

	urlCoin := GetSourceCoin(symbol) + "_" + GetTargetCode(symbol)

	var url string

	switch ex {
	case BINANCE:
		url = "https://www.binance.com/en/trade/" + urlCoin
	case BTCTURK:
		url = "https://pro.btcturk.com/pro/al-sat/" + urlCoin

	}

	return url

}

func GetSourceCoin(symbol string) string {

	var coin string
	targets := []string{"TRY", "USDT", "BTC", "BNB", "TL"}

	for _, value := range targets {
		if strings.HasSuffix(symbol, value) {
			coin = strings.ReplaceAll(symbol, value, "")
		}
	}

	return coin

}

func GetTargetCode(symbol string) string {
	targets := []string{"TRY", "USDT", "BTC", "BNB", "TL"}

	for _, value := range targets {
		if strings.HasSuffix(symbol, value) {
			return value
		}
	}

	return ""
}
