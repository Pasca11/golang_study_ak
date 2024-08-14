package main

import (
	"encoding/json"
	"fmt"
	"github.com/cinar/indicator"
	"io"
	"log"
	"net/http"
)

type Indicator interface {
	StochPrice() ([]float64, []float64)
	RSI(period int) ([]float64, []float64)
	StochRSI(rsiPeriod int) ([]float64, []float64)
	SMA(period int) []float64
	MACD() ([]float64, []float64)
	EMA() []float64
}

func UnmarshalKLines(data []byte) (KLines, error) {
	var r KLines
	err := json.Unmarshal(data, &r)
	return r, err
}
func (r *KLines) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type KLines struct {
	Pair    string   `json:"pair"`
	Candles []Candle `json:"candles"`
}
type Candle struct {
	T int64   `json:"t"`
	O float64 `json:"o"`
	C float64 `json:"c"`
	H float64 `json:"h"`
	L float64 `json:"l"`
	V float64 `json:"v"`
}
type Lines struct {
	high    []float64
	low     []float64
	closing []float64
}

func (t *Lines) StochPrice() ([]float64, []float64) {
	k, d := indicator.StochasticOscillator(t.high, t.low, t.closing)
	return k, d
}

func (t *Lines) RSI(period int) ([]float64, []float64) {
	rs, rsi := indicator.RsiPeriod(period, t.closing)
	return rs, rsi
}

func (t *Lines) SMA(period int) []float64 {
	return indicator.Sma(period, nil)
}

func (t *Lines) StochRSI(rsiPeriod int) ([]float64, []float64) {
	_, rsi := t.RSI(rsiPeriod)
	k, d := indicator.StochasticOscillator(rsi, rsi, rsi)
	return k, d
}

func (t *Lines) MACD() ([]float64, []float64) {
	return indicator.Macd(t.closing)
}
func (t *Lines) EMA() []float64 {
	return indicator.Ema(5, t.closing)
}

type LinesProxy struct {
	lines Indicator
	cache map[string][]float64
}

func (t *LinesProxy) StochPrice() ([]float64, []float64) {
	kKey := "k_stochprice"
	dKey := "d_stochprice"

	if value, ok := t.cache[kKey]; ok {
		return value, t.cache[dKey]
	}

	k, d := t.lines.StochPrice()

	t.cache[kKey] = k
	t.cache[dKey] = d

	return k, d
}

func (p *LinesProxy) RSI(period int) ([]float64, []float64) {
	rsKey := fmt.Sprintf("rs_%v", period)
	rsiKey := fmt.Sprintf("rsi_%v", period)

	if value, ok := p.cache[rsKey]; ok {
		return value, p.cache[rsiKey]
	}

	rs, rsi := p.lines.RSI(period)

	p.cache[rsKey] = rs
	p.cache[rsiKey] = rsi

	return rs, rsi
}

func (p *LinesProxy) MACD() ([]float64, []float64) {
	macdKey := "macd"
	signalKey := "signal"

	if value, ok := p.cache[macdKey]; ok {
		return value, p.cache[signalKey]
	}

	macd, signal := p.lines.MACD()
	p.cache[macdKey] = macd
	p.cache[signalKey] = signal

	return macd, signal
}

func (p *LinesProxy) SMA(period int) []float64 {
	key := fmt.Sprintf("sma_%v", period)

	if value, ok := p.cache[key]; ok {
		return value
	}

	res := p.lines.SMA(period)
	p.cache[key] = res
	return res
}

func LoadKlinesProxy(data []byte) *LinesProxy {
	lines := LoadKlines(data)
	return &LinesProxy{lines, make(map[string][]float64)}
}

func (t *LinesProxy) StochRSI(rsiPeriod int) ([]float64, []float64) {
	kKey := fmt.Sprintf("k_stochrsi_%v", rsiPeriod)
	dKey := fmt.Sprintf("d_stochrsi_%v", rsiPeriod)

	if value, ok := t.cache[kKey]; ok {
		return value, t.cache[dKey]
	}

	k, d := t.lines.StochRSI(rsiPeriod)

	t.cache[kKey] = k
	t.cache[dKey] = d

	return k, d
}

func (t *LinesProxy) EMA() []float64 {
	emaKey := "ema"

	if value, ok := t.cache[emaKey]; ok {
		return value
	}

	ema := t.lines.EMA()
	t.cache[emaKey] = ema

	return ema
}

func LoadKlines(data []byte) *Lines {
	klines, err := UnmarshalKLines(data)
	if err != nil {
		log.Fatal(err)
	}
	t := &Lines{}
	for _, v := range klines.Candles {
		t.closing = append(t.closing, v.C)
		t.low = append(t.low, v.L)
		t.high = append(t.high, v.H)
	}
	return t
}
func LoadCandles(pair string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.exmo.com/v1.1/candles_history?symbol=%s&resolution=30&from=1703056979&to=1705476839", pair), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func main() {
	pair := "BTC_USD"
	candles := LoadCandles(pair)
	lines := LoadKlinesProxy(candles)
	lines.RSI(3)
}
