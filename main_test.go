package H

import (
	"testing"
)

func TestUrl(t *testing.T) {

	P(UrlString("deneme asca Ğ Iİ c "))

}

func TestArrayFromUrl(t *testing.T) {

	var m []BtcturkSymbol

	err := ArrayFromUrl("https://api.btcturk.com/api/v2/ticker", "data", &m)

	if err != nil {
		t.Error(err)
	}
}

func TestPL(t *testing.T) {

	PL("mustafa", "dene")

}

type BitexenSymbol struct {
	Market struct {
		MarketCode          string `json:"market_code"`
		BaseCurrencyCode    string `json:"base_currency_code"`
		CounterCurrencyCode string `json:"counter_currency_code"`
	} `json:"market"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	LastPrice string `json:"last_price"`
	LastSize  string `json:"last_size"`
	Volume24H string `json:"volume_24h"`
	Change24H string `json:"change_24h"`
	Low24H    string `json:"low_24h"`
	High24H   string `json:"high_24h"`
	Avg24H    string `json:"avg_24h"`
	Timestamp string `json:"timestamp"`
}

type BtcturkSymbol struct {
	Pair              string  `json:"pair"`
	PairNormalized    string  `json:"pairNormalized"`
	Timestamp         int64   `json:"timestamp"`
	Last              float64 `json:"last"`
	High              float64 `json:"high"`
	Low               float64 `json:"low"`
	Bid               float64 `json:"bid"`
	Ask               float64 `json:"ask"`
	Open              float64 `json:"open"`
	Volume            float64 `json:"volume"`
	Average           float64 `json:"average"`
	Daily             float64 `json:"daily"`
	DailyPercent      float64 `json:"dailyPercent"`
	DenominatorSymbol string  `json:"denominatorSymbol"`
	NumeratorSymbol   string  `json:"numeratorSymbol"`
	Order             int     `json:"order"`
}

func TestBetween(t *testing.T) {

	text := `
	2236. [3:576, Hadîs No: 4367]
Şüpheli şeylerden hassasiyetle sakınan kişinin iki rekât namazı, şüpheli şeylere bulaşan kimsenin bin rekât namazından daha üstün­dür.[104]
 
2237. [4:2, Hadîs No: 4368.]
Ebû Hüreyre'den (r.a.) rivayetle:
Aklın başı, insanlarla hoş geçinmektir. Dünyada iyilik sahibi olan­lar, âhirette de iyilik sahibidirler.[105]
 
2238. [4:2, Hadîs No: 4369]
Said bin Müseyyeb rivayet ediyor:
Meşverete muhtaç olmayan hiç bir kimse yoktur.[106]
 
2239. [4:6, Hadîs No: 4378]
îbni Abbas'dan (r.a.) rivayetle:
Ben, melekleri Hanıza bin Abdulmuttalib'in ve Hanzale bin Rahb'in cenazesini yıkarken gördüm.[107]
 

`

	between := GetTextBetween(text, "Hadîs No: 4367]", "]")
	PL(between)

	b2 := GetTextBetween(between, "[", "")
	PL(b2)

}

func TestMstr_Remove(t *testing.T) {

	s := Mstr("deneme bu nasil sanki bu")

	s.Remove("bu")

	if s != "deneme  nasil sanki " {
		t.Error("Wrong result:", s)
	}
}

func TestMstr_Replace(t *testing.T) {

	s := Mstr("deneme bu nasil sanki bu")

	s.Replace("bu", "mzgs")

	if s != "deneme mzgs nasil sanki mzgs" {
		t.Error("Wrong result:", s)
	}

}

func TestMstr_Between(t *testing.T) {

	s := Mstr("deneme bu nasil sanki bu")

	between := s.Between("bu", "sanki")

	if between != " nasil " {
		t.Error("Wrong result:", between)
	}

}

func TestAlignText(t *testing.T) {

	//s1 := "first line "
	//s2 := "second line "
	//s3 := "third line "
	//
	//f1 := s1
	//f2 := s1 + "hello"
	//f3 := s1 + "hello world"
	//
	//
	//P(AlignText(40, f1, s2+"sad asd asd asd asd", s3+" dene", "asca sacasc"))
	//P(AlignText(40, f2, s2+"n", s3, "sdvn d"))
	//P(AlignText(40, f3, s2+" yes", s3, "acc "))

	P(AlignText(20, "LINKUSDT", "$12.5546", "-4.6%"))
	P(AlignText(20, "ATOMUSDT", "$5.292", "-3.55%"))
	P(AlignText(20, "LINKTRY", "₺9337", "-4.24%"))

}
