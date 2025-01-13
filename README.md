# The LTP service

As outlined in the [requirements](requirements.md), this service provides last-traded price data
with accuracy up to one minute for the following cryptocurrency pairs:

* BTC/CHF
* BTC/EUR
* BTC/USD

## Notes

This service utilizes caching for faster responses and is designed with clean architecture principles in mind.

## How to run

- `make dc` runs docker-compose with the app container on port 8090.
- `make test` runs the tests
- `make run` runs the app locally on port 8090 without docker.
- `make lint` runs the linter

## Request examples

```
curl --location 'http://127.0.0.1:8090/api/v1/ltp'
curl --location 'http://127.0.0.1:8090/api/v1/ltp?pairs=BTC/CHF,BTC/USD'
```