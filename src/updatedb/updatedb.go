package updatedb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
)

type Message struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
	Update string  `json:"update_time"`
}

func updateStock(connInfo string, redisURL string, ticker string) {
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	redisClient, err := redis.Dial("tcp", redisURL)
	if err != nil {
		fmt.Println(err)
	}
	defer redisClient.Close()

	rand.Seed(time.Now().Unix())
	for {
		time.Sleep(time.Duration(rand.Int31n(10000)) * time.Millisecond)
		var price string
		err := db.QueryRow("select price from portfolio where ticker=$1", ticker).Scan(&price)
		switch {
		case err == sql.ErrNoRows:
			log.Fatalln("No user with that ID.")
		case err != nil:
			log.Fatal(err)
			// default:
			// 	fmt.Printf("Price is %s\n", price)
		}

		priceFloat, err := strconv.ParseFloat(price, 64)
		newPrice := priceFloat + 3*rand.NormFloat64()
		// db.Exec("update portfolio set price=$1 where ticker=$2", newPrice, ticker)

		// db.Exec("insert into update_stock_price (ticker, price) values ($1, $2)", ticker, newPrice)
		var updateTime string
		db.QueryRow("insert into update_stock_price (ticker, price) values ($1, $2) returning update_time at time zone 'EST5EDT'", ticker, newPrice).Scan(&updateTime)
		fmt.Println(updateTime)

		message := Message{
			Ticker: ticker,
			Price:  newPrice,
			Update: updateTime,
		}
		output, err := json.MarshalIndent(&message, "", "\t")
		// message := fmt.Sprintf("%s price updated to %.4f", ticker, newPrice)
		_, err = redisClient.Do("PUBLISH", "chan1", output)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func UpdateDB(connInfo string, redisURL string) {

	portfolio := []string{"AAPL", "GOOG", "AMZN", "FB"}
	for _, stock := range portfolio {
		go updateStock(connInfo, redisURL, stock)
	}
	// go subscribeRedis()
}
