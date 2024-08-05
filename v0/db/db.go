package db

import (
	"os"
	"fmt"
	"time"
	"strings"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"coinex-api/v0/pkg/views"
)

//Global Repository variable
var Repository DB

//Database struct with connection and context
type DB struct {
	conn *pgxpool.Pool
	ctx  context.Context
}

func (db *DB) Init(ctx context.Context, url string) {
	var err error

	db.ctx = ctx
	db.conn, err = pgxpool.New(db.ctx, url)
	if err != nil {
		fmt.Println("Error connecting to database", err)
		return
	}

	_, err = db.conn.Exec(db.ctx, "create table coins ( code text, name text, price real, currency text, updated timestamptz not null )")
	if err != nil {
		fmt.Println("Error creating schemas", err)
		return
	}
}

func (db DB) GetAllCoin() (views.CoinList, error) {
	coins := make(views.CoinList)

	rows, err := db.conn.Query(db.ctx, "select code, name from coins")
	if err != nil {
		fmt.Println("Error in fetching data", err)
		return coins, err
	}
	defer rows.Close()

	for rows.Next() {
		var code, name string
		rows.Scan(&code, &name)

		coins[code] = views.Coin{Code: code, Name: name}
	}

	return coins, err
}

func (db DB) GetCoin(coinCode string) (views.Coin, error) {

	expiry, err := db.GetExpiry(coinCode, "")
	if err != nil {
		fmt.Println("Error fetching expiry", err)
		return views.Coin{}, err
	}
	if time.Now().After(expiry) {
		return views.Coin{}, pgx.ErrNoRows
	}

	var code, name string

	err = db.conn.QueryRow(db.ctx, "select code, name from coins where code=$1", coinCode).Scan(&code, &name)
	if err != nil {
		fmt.Println("Coin doesn't exists in database, fetching from API...", err)
		return views.Coin{}, err
	}

	return views.Coin{Code: code, Name: name}, err
}

func (db DB) GetCoinInCurrency(coinCode, currency string) (views.Coin, error) {

	var code, name string

	err := db.conn.QueryRow(db.ctx, "select code, name from coins where code=$1 and currency=$2", coinCode, currency).Scan(&code, &name)
	if err != nil {
		fmt.Println("Coin doesn't exists in database, fetching from API...", err)
		return views.Coin{}, err
	}

	return views.Coin{Code: code, Name: name}, err
}

func (db DB) GetPriceIn(coinCode, currency string) (views.Price, error) {

	expiry, err := db.GetExpiry(coinCode, "")
	if err != nil {
		fmt.Println("Error fetching expiry", err)
		return views.Price{}, err
	}
	if time.Now().After(expiry) {
		return views.Price{}, pgx.ErrNoRows
	}

	var price float64

	err = db.conn.QueryRow(db.ctx, "select price from coins where code=$1 and currency=$2", coinCode, currency).Scan(&price)
	if err != nil {
		fmt.Println("Coin doesn't exists in database, fetching from API...", err)
		return views.Price{}, err
	}
	if price < 0 {
		fmt.Println("Price doesn't exists in database, fetching from API...", pgx.ErrNoRows)
		return views.Price{}, pgx.ErrNoRows
	}

	return views.Price{Coin: coinCode, Currency: currency, Price: price}, err
}

func (db DB) GetPrice(coinCode string) (views.PriceResponse, error) {
	expiry, err := db.GetExpiry(coinCode, "")
	if err != nil {
		fmt.Println("Error fetching expiry", err)
		return views.PriceResponse{}, err
	}
	if time.Now().After(expiry) {
		return views.PriceResponse{}, pgx.ErrNoRows
	}
	currencies := strings.Split(os.Getenv("CURRENCY_LIST"), "-")
	response := make(map[string]float64)
	var price float64

	for _, currency := range currencies {
		err := db.conn.QueryRow(db.ctx, "select price from coins where code=$1 and currency=$2", coinCode, currency).Scan(&price)
		if err != nil {
			fmt.Println("Coin doesn't exists in database, fetching from API...", err)
			return views.PriceResponse{}, err
		}
		if price < 0 {
			fmt.Println("Price doesn't exists in database, fetching from API...", pgx.ErrNoRows)
			return views.PriceResponse{Data: views.PriceResponseData{coinCode: {"": -1.0}}}, pgx.ErrNoRows
		}
		response[currency] = price
	}

	return views.PriceResponse{Data: views.PriceResponseData{coinCode: response}}, nil
}

func (db DB) GetExpiry(coinCode, currency string) (time.Time, error) {
	if currency == "" {
		currency = os.Getenv("CURRENCY")
	}

	var expiry time.Time

	err := db.conn.QueryRow(db.ctx, "select updated from coins where code=$1 and currency=$2", coinCode, currency).Scan(&expiry)
	if err != nil {
		fmt.Println("Coin doesn't exists in database, fetching from API...", err)
		return time.Now(), err
	}

	EXPIRY, err := time.ParseDuration(os.Getenv("EXPIRY"))
	if err != nil {
		fmt.Println("Invalid time format in environment variable 'EXPIRY'", err)
		return time.Now(), err
	}

	expiry = expiry.Add(EXPIRY)
	return expiry, err
}

func (db DB) SetCoin(coin views.Coin) error {
	currencies := strings.Split(os.Getenv("CURRENCY_LIST"), "-")
	for _, currency := range currencies {
		_, err := db.conn.Exec(db.ctx, "insert into coins (code, name, price, currency, updated) values ($1, $2, $3, $4, $5)", coin.Code, coin.Name, -1.0, currency, time.Now())
		if err != nil {
			fmt.Println("Error writing to database...", err)
			return err
		}
	}
	return nil 
}

func (db DB) SetPrice(price views.Price) error {
	_, err := db.conn.Exec(db.ctx, "update coins set price=$3, updated=$4 where currency=$2 and code=$1", price.Coin, price.Currency, price.Price, time.Now())
	if err != nil {
		fmt.Println("Error writing to database...", err)
		return err
	}
	return err
}

func (db DB) Close() {
	_, err := db.conn.Exec(db.ctx, "drop table coins")
	if err != nil {
		fmt.Println("Error dropping schemas", err)
	}
	defer db.conn.Close()
}
