# Currency Rate Restful API based on Golang

> Simple RESTful API to get latest, daily currency rates and analyze

## Quick Start

``` bash
go build
./go-currencyrate
```

## Endpoints

### Get Lastest Rates
``` bash
GET /rates/latest
```

#### Sample Response
```json
{
    "base": "EUR",
    "rates": {
        "AUD": 1.5339,
        "BGN": 1.9558,
        "USD": 1.2023,
        "ZAR": 14.8845
    }
}
```

### Get Rates by Date
``` bash
GET /rates/YYYY-MM-DD
```

#### Sample Response
```json
{
    "base": "EUR",
    "rates": {
        "AUD": 1.5339,
        "BGN": 1.9558,
        "USD": 1.2023,
        "ZAR": 14.8845
    }
}
```

### Analyze Rate
``` bash
GET /rates/analyze
```

#### Sample Response
```json
{
    "base": "EUR",
    "rates_analyze": {
        "AUD": {
            "min": 1.4994,
            "max": 1.5693,
            "avg": 1.5340524590163933
        },
        "BGN": {
            "min": 1.9558,
            "max": 1.9558,
            "avg": 1.9557999999999973
        },
        "USD": {
            "min": 1.1562,
            "max": 1.2065,
            "avg": 1.1783852459016388
        },
        "ZAR": {
            "min": 14.7325,
            "max": 17.0212,
            "avg": 16.06074426229508
        }
    }
}
```

## App Info

### Version

1.0.0