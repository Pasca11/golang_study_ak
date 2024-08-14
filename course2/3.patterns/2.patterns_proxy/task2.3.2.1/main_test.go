package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockIndicator struct {
	mock.Mock
}

func (m *MockIndicator) RSI(period int) ([]float64, []float64) {
	args := m.Called(period)
	return args.Get(0).([]float64), args.Get(1).([]float64)
}

func (m *MockIndicator) StochRSI(rsiPeriod int) ([]float64, []float64) {
	args := m.Called()
	return args.Get(0).([]float64), args.Get(1).([]float64)
}

func (m *MockIndicator) MACD() ([]float64, []float64) {
	args := m.Called()
	return args.Get(0).([]float64), args.Get(1).([]float64)
}

func (m *MockIndicator) EMA() []float64 {
	args := m.Called()
	return args.Get(0).([]float64)
}

func (m *MockIndicator) SMA(period int) []float64 {
	args := m.Called(period)
	return args.Get(0).([]float64)
}

func (m *MockIndicator) StochPrice() ([]float64, []float64) {
	args := m.Called()
	return args.Get(0).([]float64), args.Get(1).([]float64)
}

func TestLinesProxy_RSI(t *testing.T) {
	mock := new(MockIndicator)
	period := 3
	mockRS := []float64{1.0, 1.0, 3.0}
	mockRSI := []float64{4.0, 5.0, 6.0}
	mock.On("RSI", period).Return(mockRS, mockRSI)

	lProxy := LinesProxy{
		lines: mock,
		cache: make(map[string][]float64),
	}

	rs, rsi := lProxy.RSI(period)

	assert.Equal(t, mockRS, rs)
	assert.Equal(t, mockRSI, rsi)

	rsCached, rsiCached := lProxy.RSI(period)
	assert.Equal(t, mockRS, rsCached)
	assert.Equal(t, mockRSI, rsiCached)

	mock.AssertNumberOfCalls(t, "RSI", 1)
}

func TestLinesProxy_StochPrice(t *testing.T) {
	mockIndicator := new(MockIndicator)
	mockK := []float64{1.0, 2.0, 3.0}
	mockD := []float64{4.0, 5.0, 6.0}
	mockIndicator.On("StochPrice").Return(mockK, mockD)

	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	k, d := proxy.StochPrice()
	assert.Equal(t, mockK, k)
	assert.Equal(t, mockD, d)

	kCached, dCached := proxy.StochPrice()
	assert.Equal(t, k, kCached)
	assert.Equal(t, d, dCached)

	mockIndicator.AssertNumberOfCalls(t, "StochPrice", 1)
}

func TestLinesProxy_MACD(t *testing.T) {
	mockIndicator := new(MockIndicator)
	mockMACD := []float64{1.0, 2.0, 3.0}
	mockSignal := []float64{4.0, 5.0, 6.0}
	mockIndicator.On("MACD").Return(mockMACD, mockSignal)

	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	macd, signal := proxy.MACD()
	assert.Equal(t, mockMACD, macd)
	assert.Equal(t, mockSignal, signal)

	macdCached, signalCached := proxy.MACD()
	assert.Equal(t, macd, macdCached)
	assert.Equal(t, signal, signalCached)

	mockIndicator.AssertNumberOfCalls(t, "MACD", 1)
}

func TestLines_StochRSI(t *testing.T) {
	mockIndicator := new(MockIndicator)
	mockK := []float64{1.0, 2.0, 3.0}
	mockD := []float64{4.0, 5.0, 6.0}
	mockIndicator.On("StochRSI", mock.Anything).Return(mockK, mockD)

	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	k, d := proxy.StochRSI(14)
	assert.Equal(t, mockK, k)
	assert.Equal(t, mockD, d)

	kCached, dCached := proxy.StochRSI(14)
	assert.Equal(t, k, kCached)
	assert.Equal(t, d, dCached)

	mockIndicator.AssertNumberOfCalls(t, "StochRSI", 1)
}

func TestLines_EMA(t *testing.T) {
	mockIndicator := new(MockIndicator)
	mockEMA := []float64{10.0, 20.0, 30.0}
	mockIndicator.On("EMA").Return(mockEMA)

	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	ema := proxy.EMA()
	assert.Equal(t, mockEMA, ema)

	emaCached := proxy.EMA()
	assert.Equal(t, ema, emaCached)

	mockIndicator.AssertNumberOfCalls(t, "EMA", 1)
}

func TestLines_SMA(t *testing.T) {
	mockIndicator := new(MockIndicator)
	mockSMA := []float64{10.0, 20.0, 30.0}
	number := 5
	mockIndicator.On("SMA", number).Return(mockSMA)

	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	sma := proxy.SMA(5)
	assert.Equal(t, mockSMA, sma)

	smaCached := proxy.SMA(5)
	assert.Equal(t, sma, smaCached)

	mockIndicator.AssertNumberOfCalls(t, "SMA", 1)
}
