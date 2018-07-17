package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Orders main struct
type Orders struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	OrderID  int    `json:"orderid"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/orders", GetOrders).Methods("GET")     //select *
	router.HandleFunc("/orders/", CreateOrder).Methods("POST") //insert into
	router.HandleFunc("/orders/{id}", GetID).Methods("GET")    // select * where
	log.Fatal(http.ListenAndServe(":8000", router))

}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

var ord [50]Orders

//GetOrders This function handles Get Request
func GetOrders(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root123@tcp(127.0.0.1:3306)/test")
	//var id int = 9555

	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT Item, Quantity, OrderID FROM orders")
	if err != nil {
		panic(err.Error())
	}
	i := 0

	for results.Next() {
		err = results.Scan(&ord[i].Item, &ord[i].Quantity, &ord[i].OrderID)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Printf("You have ordered %d %ss and your OrderID is : %d\n", ord[i].Quantity, ord[i].Item, ord[i].OrderID)
		i = i + 1
	}
	for i := 0; i < 50; i++ {
		if ord[i].Item == "" {
			break
		} else {
			json.NewEncoder(w).Encode(ord[i])
		}
	}
}

//CreateOrder This function handles Post Request
func CreateOrder(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var data Orders
	err := decoder.Decode(&data)
	if err != nil {
		//panic(err)
		return
	}

	db, err := sql.Open("mysql", "root:root123@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	data.OrderID = random(1, 1000)
	data.OrderID = 955

	for {
		results, err := db.Query("select * from orders where OrderID = ?;", data.OrderID)
		if err != nil {
			panic(err.Error())
		}
		if results.Next() {
			data.OrderID = random(1, 1000)
		} else {
			break
		}
	}
	results, err := db.Query("insert into orders values(?, ?, ?);", data.Item, data.Quantity, data.OrderID)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(results)
	json.NewEncoder(w).Encode(data)
	return
}

//GetID This function handles Get Requests for a given Order ID
func GetID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range ord {
		if strconv.Itoa(item.OrderID) == params["id"] {
			db, err := sql.Open("mysql", "root:root123@tcp(127.0.0.1:3306)/test")
			//var id int = 9555
			if err != nil {
				log.Print(err.Error())
			}
			defer db.Close()

			results, err := db.Query("select * from orders where OrderID = ?;", params["id"])
			if err != nil {
				panic(err.Error())
			}
			var out Orders
			for results.Next() {
				err = results.Scan(&out.Item, &out.Quantity, &out.OrderID)
				if err != nil {
					panic(err.Error())
				}
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	}
}
