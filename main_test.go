package H

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUrl(t *testing.T) {

	P(UrlString("deneme asca Ğ Iİ c "))

}

func TestArrayFromUrl(t *testing.T) {

	var m []BtcturkSymbol

	err := ParseJsonFromUrlGET("https://api.btcturk.com/api/v2/ticker", "data", &m)

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

func TestFM(t *testing.T) {

	num := 12
	s := "mustafa"
	b := true

	fm := FM("num:{num} str:{s} b:{b} {s}", num, s, b, s)

	PL(fm)

}

func TestFtoStr(t *testing.T) {

	PL(FtoStr(3243432423423.234323434))
}

func TestMstr_Lines(t *testing.T) {

	a := Mstr(ReadFile("/Users/mustafa/Documents/BitbarPlugins/data/ticker.txt"))

	for _, value := range a.Lines() {
		P("-" + value + "-")
	}

}

func TestTitleTurkish(t *testing.T) {

	s := "Ahirette insan ölümsüzdür Ilık İlk"
	P("Original: ", s)

	//PL(strings.ToUpperSpecial(unicode.TurkishCase,"ı ç ğ ö i"))

	PL(TitleTurkish(s))

}

func TestReverseAny(t *testing.T) {

	a := []int{1, 2, 3, 4}
	P(a)

	ReverseAny(a)

	P(a)

}

func TestGetRequest(t *testing.T) {

	b, err := GetRequest("https://postb.in/1604876410544-8237676401622", "X-Status", "mustafa")
	if err != nil {
		P(err)
	}

	P(string(b))

}

func TestParseJsonFromUrl(t *testing.T) {

	var posts []Post
	err := ParseJsonFromUrlGET("https://jsonplaceholder.typicode.com/posts", "", &posts)
	if err != nil {
		P(err)
	}

	var symbols []Symbol
	err = ParseJsonFromUrlGET("https://api.btcturk.com/api/v2/ticker", "data", &symbols)
	if err != nil {
		P(err)
	}

	var markets []BitexenMarket
	err = ParseJsonFromUrlGET("https://www.bitexen.com/api/v1/market_info/", "data.markets", &markets)
	if err != nil {
		P(err)
	}

	for _, value := range markets {
		P(value.MarketCode, value.PresentationDecimal)
	}

}

func TestMstr_FixMultiSpace(t *testing.T) {
	s := Mstr("Bize Mahled b. Mâlik haber verip (dedi ki)         bize Yahya b. Sa'îd rivâyet edip (dedi ki) bize İbn Cureyc rivâyet edip         dedi ki; bana Abdulmelik b. Abdillah b. Ebî Süfyân es -Sakafî, İbn         Ömer'den haber verdi ki, o şöyle demiş: Bu ilmi yazıyla kaydediniz.")

	PBL(s)

	s.FixMultiSpace()
	PBL(s)

}

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Symbol struct {
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

type BitexenMarket struct {
	MarketCode             string `json:"market_code"`
	URLSymbol              string `json:"url_symbol"`
	BaseCurrency           string `json:"base_currency"`
	CounterCurrency        string `json:"counter_currency"`
	MinimumOrderAmount     string `json:"minimum_order_amount"`
	MaximumOrderAmount     string `json:"maximum_order_amount"`
	BaseCurrencyDecimal    int    `json:"base_currency_decimal"`
	CounterCurrencyDecimal int    `json:"counter_currency_decimal"`
	PresentationDecimal    int    `json:"presentation_decimal"`
	ResellMarket           bool   `json:"resell_market"`
}

func TestUrlStringForFile(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test jpeg image file",
			args: args{"türkÇe RESİM Ğiş ŞQ I Ç.jpg"},
			want: "turkce-resim-gis-sq-i-c.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UrlStringForFile(tt.args.s); got != tt.want {
				t.Errorf("UrlStringForFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB(t *testing.T) {
	InitMongoDB("mongotest")

	type User struct {
		ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		Language string             `label:"Dil" type:"language" `
		Title    string             `label:"Başlık" type:"h_text" `
		Order    int                `label:"Sıra" type:"h_text" has:"hr"`
		Content  string             `label:"Detaylar" type:"rich" has:"hr"`
		Date     int64
	}

	var user []User
	DB.Find(&user).Sort("-language", "order").Limit(2).All()

	var userOne User
	DB.Find(&userOne).Sort("language").One()
	PL(userOne)

	PL(len(user))

	for _, v := range user {
		P(v.Language, v.Title, v.Order)
	}

	t.Errorf("UrlStringForFile() = %v, want %v", "sca", "sdvf")

}
