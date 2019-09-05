package dao

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Rate struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	RateDate string        `bson:"rate_date" json:"rateDate"`
	Rates    []*RateItem   `bson:"rates" json:"rates"`
}

type RateItem struct {
	Currency string  `bson:"currency" json:"currency"`
	Rate     float32 `bson:"rate" json:"rate"`
}

type AnalyzeResult struct {
	Currency string  `bson:"_id" json:"Currency"`
	Max      float32 `bson:"max" json:"max"`
	Min      float32 `bson:"min" json:"min"`
	Avg      float32 `bson:"avg" json:"avg"`
}

type RatesDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	RATES_COLLECTION = "rates"
)

// Establish a connection to database
func (m *RatesDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of rates
func (m *RatesDAO) FindAll() ([]Rate, error) {
	var rates []Rate
	err := db.C(RATES_COLLECTION).Find(nil).All(&rates)
	return rates, err
}

func (m *RatesDAO) FindById(id string) (Rate, error) {
	var rate Rate
	err := db.C(RATES_COLLECTION).FindId(bson.ObjectIdHex(id)).One(&rate)
	return rate, err
}

func (m *RatesDAO) GetLatest() (Rate, error) {
	var rate Rate
	err := db.C(RATES_COLLECTION).Find(nil).Sort("-rate_date").One(&rate)
	return rate, err
}

func (m *RatesDAO) FindByDate(date string) (*Rate, error) {
	var rate Rate
	err := db.C(RATES_COLLECTION).Find(bson.M{"rate_date": date}).One(&rate)
	return &rate, err
}

func (m *RatesDAO) Analyze() ([]*AnalyzeResult, error) {
	pipe := db.C(RATES_COLLECTION).Pipe([]bson.M{
		{"$unwind": "$rates"},
		{"$project": bson.M{
			"_id":       1,
			"rate_date": 1,
			"currency":  "$rates.currency",
			"rate":      "$rates.rate",
		}},
		{"$group": bson.M{
			"_id": "$currency",
			"max": bson.M{"$max": "$rate"},
			"min": bson.M{"$min": "$rate"},
			"sum": bson.M{"$sum": "$rate"},
			"avg": bson.M{"$avg": "$rate"},
		}},
		{
			"$sort": bson.M{"_id": 1},
		},
	})
	resp := []*AnalyzeResult{}
	err := pipe.All(&resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *RatesDAO) Save(rate *Rate) error {
	oldRate, err := m.FindByDate(rate.RateDate)
	if err != nil || oldRate == nil {
		rate.ID = bson.NewObjectId()
		err = m.Insert(rate)
	} else {
		rate.ID = oldRate.ID
		err = m.Update(rate)
	}
	return err
}

// Insert a rate into database
func (m *RatesDAO) Insert(rate *Rate) error {
	err := db.C(RATES_COLLECTION).Insert(rate)
	return err
}

// Update an existing rate
func (m *RatesDAO) Update(rate *Rate) error {
	err := db.C(RATES_COLLECTION).UpdateId(rate.ID, rate)
	return err
}
