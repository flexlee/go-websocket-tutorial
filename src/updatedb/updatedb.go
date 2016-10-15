package updatedb

import (
	"database/sql"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func updateStock(connInfo string, ticker string) {
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
		// updateStmt := fmt.Sprintf("update portfolio set price=%.4f where ticker='%s';", priceFloat+3*rand.NormFloat64(), ticker)
		db.Exec("update portfolio set price=$1 where ticker=$2", priceFloat+3*rand.NormFloat64(), ticker)
	}
}

func UpdateDB(connInfo string) {

	portfolio := []string{"AAPL", "GOOG", "AMZN", "FB"}
	for _, stock := range portfolio {
		go updateStock(connInfo, stock)
	}
}
