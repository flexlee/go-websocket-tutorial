package updatedb

import (
	"database/sql"
	"fmt"
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
		time.Sleep(time.Duration(rand.Int31n(5000)) * time.Millisecond)
		selectStmt := fmt.Sprintf("select price from portfolio where ticker='%s';", ticker)
		var price string
		err := db.QueryRow(selectStmt).Scan(&price)
		switch {
		case err == sql.ErrNoRows:
			log.Printf("No user with that ID.")
		case err != nil:
			log.Fatal(err)
			// default:
			// 	fmt.Printf("Price is %s\n", price)
		}

		priceFloat, err := strconv.ParseFloat(price, 64)
		updateStmt := fmt.Sprintf("update portfolio set price=%.4f where ticker='%s';", priceFloat+3*rand.NormFloat64(), ticker)
		db.Exec(updateStmt)
	}
}

func UpdateDB(connInfo string) {

	portfolio := []string{"AAPL", "GOOG", "AMZN", "FB"}
	for _, stock := range portfolio {
		go updateStock(connInfo, stock)
	}
}
