package db

import (
	"context"
	"fmt"
	"os"
	"time"

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

func (db *DB) InitDB(ctx context.Context, url string) {
	var err error

	db.ctx = ctx
	db.conn, err = pgxpool.New(db.ctx, url)
	if err != nil {
		fmt.Println("Error connecting to database", err)
		return
	}

	_, err = db.conn.Exec(db.ctx, "create table coins ( code text primary key, name text, price real, currency text, updated timestamptz not null )")
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

	var code, name string

	err := db.conn.QueryRow(db.ctx, "select code, name from coins where code=$1", coinCode).Scan(&code, &name)
	if err != nil {
		fmt.Println("Coin doesn't exists in database, fetching from API...", err)
		return views.Coin{}, err
	}

	return views.Coin{Code: code, Name: name}, err
}

func (db DB) GetPrice(coinCode, currency string) (views.Price, error) {

	var price float64

	err := db.conn.QueryRow(db.ctx, "select price from coins where code=$1 and currency=$2", coinCode, currency).Scan(&price)
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

func (db DB) GetExpiry(coinCode string) (time.Time, error) {
	var expiry time.Time

	err := db.conn.QueryRow(db.ctx, "select updated from coins where code=$1", coinCode).Scan(&expiry)
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
	_, err := db.conn.Exec(db.ctx, "insert into coins (code, name, price, currency, updated) values ($1, $2, $3, $4, $5)", coin.Code, coin.Name, -1.0, os.Getenv("CURRENCY"), time.Now())
	if err != nil {
		fmt.Println("Error writing to database...", err)
		return err
	}
	return err
}

func (db DB) SetPrice(price views.Price) error {
	_, err := db.conn.Exec(db.ctx, "update coins set price=$3, updated=$4 where currency=$2 and code=$1", price.Coin, price.Currency, price.Price, time.Now())
	if err != nil {
		fmt.Println("Error writing to database...", err)
		return err
	}
	return err
}

func (db *DB) Close() {
	_, err := db.conn.Exec(db.ctx, "drop table coins")
	if err != nil {
		fmt.Println("Error dropping schemas", err)
	}
	defer db.conn.Close()
}
