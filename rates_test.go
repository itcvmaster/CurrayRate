package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"./app"
	"./models"
	"github.com/gorilla/mux"
)

func TestLatestRate(t *testing.T) {
	app.Init()
	req, err := http.NewRequest("GET", "/rates/latest", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.LatestRate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var dailyRateResult models.DailyRateResult
	err = json.Unmarshal(rr.Body.Bytes(), &dailyRateResult)
	if err != nil {
		t.Errorf("Result format is not correct: got %v",
			rr.Body.String())
	}
	if dailyRateResult.Rates == nil {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
	if _, ok := dailyRateResult.Rates["USD"]; !ok {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
}

func TestDailyRate(t *testing.T) {
	app.Init()
	req, err := http.NewRequest("GET", "/rates/2019-08-20", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()

	vars := map[string]string{
		"date": "2019-08-20",
	}

	req = mux.SetURLVars(req, vars)

	handler := http.HandlerFunc(app.DailyRate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var dailyRateResult models.DailyRateResult
	err = json.Unmarshal(rr.Body.Bytes(), &dailyRateResult)
	if err != nil {
		t.Errorf("Result format is not correct: got %v",
			rr.Body.String())
	}
	if dailyRateResult.Rates == nil {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
	if _, ok := dailyRateResult.Rates["USD"]; !ok {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
}

func TestAnalyze(t *testing.T) {
	app.Init()
	req, err := http.NewRequest("GET", "/rates/2019-08-20", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.AnalyzeRates)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var result models.RateAnalysisResult
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("Result format is not correct: got %v",
			rr.Body.String())
	}
	if result.Rates == nil {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}

	_, ok := result.Rates["USD"]
	if !ok {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
}
