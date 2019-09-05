package app

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	. "../config"
	. "../dao"
	. "../models"
	"github.com/gorilla/mux"
)

var config = Config{}
var dao = RatesDAO{}

// Parse the configuration file 'config.toml', and establish a connection to DB
func Init() {
	config.Load()
	log.Println("Configuration loaded")

	dao.Server = config.Database.Host
	dao.Database = config.Database.DbName
	dao.Connect()
	log.Println("DB connected")
}

func ImportRates() {
	client := http.Client{}

	req, err := http.NewRequest("GET",
		config.ApiUrl, nil)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var apiResponse ApiResponse
	err = xml.Unmarshal(respBody, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	for _, cube := range apiResponse.CubeByDates {
		rateItems := []*RateItem{}
		for _, item := range cube.Cubes {
			rateItems = append(rateItems, &RateItem{
				Currency: item.Currency,
				Rate:     item.Rate,
			})
		}

		rate := &Rate{
			RateDate: cube.Time,
			Rates:    rateItems,
		}
		if err := dao.Save(rate); err != nil {
			log.Fatal(err)
		}
	}
}

func LatestRate(w http.ResponseWriter, r *http.Request) {
	Rate, err := dao.GetLatest()
	if err != nil {
		log.Println("LatestRate, error on GetLatest", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rates := map[string]float32{}
	for _, rateItem := range Rate.Rates {
		rates[rateItem.Currency] = rateItem.Rate
	}

	result := &DailyRateResult{
		Base:  "EUR",
		Rates: rates,
	}

	respondWithJson(w, http.StatusOK, result)
}

func DailyRate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	Rate, err := dao.FindByDate(params["date"])

	if err != nil {
		log.Println("DailyRate, error on FindByDate", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rates := map[string]float32{}
	for _, rateItem := range Rate.Rates {
		rates[rateItem.Currency] = rateItem.Rate
	}

	result := &DailyRateResult{
		Base:  "EUR",
		Rates: rates,
	}
	respondWithJson(w, http.StatusOK, result)
}

func AnalyzeRates(w http.ResponseWriter, r *http.Request) {
	analyzeResult, err := dao.Analyze()
	if err != nil {
		log.Println("AnalyzeRates error", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	result := &RateAnalysisResult{
		Base:  "EUR",
		Rates: map[string]*AnalysisData{},
	}

	for _, rate := range analyzeResult {
		analyzeData := &AnalysisData{
			Min: rate.Min,
			Max: rate.Max,
			Avg: rate.Avg,
		}
		result.Rates[rate.Currency] = analyzeData
	}

	respondWithJson(w, http.StatusOK, result)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Route() {
	Init()
	ImportRates()

	r := mux.NewRouter()
	r.HandleFunc("/rates/latest", LatestRate).Methods("GET")
	r.HandleFunc("/rates/analyze", AnalyzeRates).Methods("GET")
	r.HandleFunc("/rates/{date}", DailyRate).Methods("GET")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r); err != nil {
		log.Fatal(err)
	}
}
