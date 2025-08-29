# Currency Converter

Aplikacja to API do pobierania kursów walut i wykonywania konwersji między walutami (zarówno fiat jak i krypto).  
Projekt napisany z wykorzystaniem [Gin](https://gin-gonic.com/) i [shopspring/decimal](https://github.com/shopspring/decimal).

---

## Endpointy

### `GET /rates`
Zwraca aktualne kursy wymiany dla podanych walut względem siebie.

**Parametry query:**
- `currencies` – lista kodów walut oddzielona przecinkami (np. `USD,EUR,GBP`).

**Przykład:**
- `GET /rates?currencies=USD,EUR,GBP`

**Odpowiedź:**
```json
[
  {"from": "USD", "to": "EUR", "rate": "0.92"},
  {"from": "EUR", "to": "USD", "rate": "1.09"},
  {"from": "USD", "to": "GBP", "rate": "0.78"}
]
```

### `GET /exchange`
Przelicza podaną kwotę z jednej waluty na inną.

**Parametry query:**
- `from` – waluta źródłowa (np. WBTC)
- `to` – waluta docelowa (np. USDT)
- `amount` – kwota do przeliczenia

**Przykład:**
- `GET /exchange?from=WBTC&to=USDT&amount=1.0`

**Odpowiedź:**
```json
{ "from": "WBTC", "to": "USDT", "amount": 57094.314314 }
```

---

## Uruchomienie

### Wymagania
- Go 1.22+ lub
- Docker / Docker Compose

## Uruchomienie lokalne (Musisz mieć Golang zainstalowany lokalnie)
Jeśli masz własny <api_key> dla `https://openexchangerates.org/` przypisz jego wartość do zmiennej środowiskowej OPENEXCHANGE_APP_ID. 
Jeśli nie masz, znajdziesz mój <api_key> w pliku `docker-compose.yaml`
`export OPENEXCHANGE_APP_ID=<api_key>`
`export SERVER_PORT=3001`

`go run ./cmd/app`

Serwer wystartuje na `http://localhost:3001`.

## Uruchomienie w Dockerze
`docker compose up`

Serwer również wystartuje na `http://localhost:3001`.

## Przykłady `curl`
`curl 'localhost:3001/rates?currencies=USD,GBP,EUR'`
`curl 'localhost:3001/exchange?from=USDT&to=BEER&amount=1.0'`
