package models

type ApiResponse struct {
	CubeByDates []*CubeByDate `xml:"Cube>Cube"`
}

type CubeByDate struct {
	Time  string  `xml:"time,attr"`
	Cubes []*Cube `xml:"Cube"`
}

type Cube struct {
	Currency string  `xml:"currency,attr"`
	Rate     float32 `xml:"rate,attr"`
}

type DailyRateResult struct {
	Base  string             `json:"base"`
	Rates map[string]float32 `json:"rates"`
}

type RateAnalysisResult struct {
	Base  string                   `json:"base"`
	Rates map[string]*AnalysisData `json:"rates_analyze"`
}

type AnalysisData struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
	Avg float32 `json:"avg"`
}
