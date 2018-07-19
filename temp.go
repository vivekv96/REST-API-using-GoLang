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
	"github.com/jmoiron/sqlx"
)

//Orders main struct
type Orders struct {
	Item     []string `db:"item"`
	Quantity []int    `db:"quantity"`
	OrderID  int      `db:"orderid"`
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

var ord [10]Orders

//GetOrders This function handles Get Request
func GetOrders(w http.ResponseWriter, r *http.Request) {
	db, err := sqlx.Open("mysql", "root:root123@tcp(127.0.0.1:3306)/orders")
	//var id int = 9555

	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	// rows, err := db.Query("SELECT * FROM items order by orderid")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// rows, err := db.Queryx("SELECT * FROM items order by orderid;")

	//db.Select(&ord, "SELECT * FROM items")

	//fmt.Println(ord)

	ord1 := []Orders{}
	rows, err := db.Queryx("select * from items")
	columnNames, _ := rows.Columns()
	fmt.Println(columnNames)
	for rows.Next() {
		fmt.Println(rows)
		err := rows.StructScan(&ord1)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", ord1)
	}
	// for rows.Next() {
	// 	err = rows.Scan(&ord[0].Item, &ord[0].Quantity, &ord[0].OrderID)
	// 	tempOrderID := ord[0].OrderID
	// 	for r
	// }

	// items := []Orders{{[]string{"colgate", "patanjali"}, []int{10, 1}, 955},
	// 	{[]string{"sdff", "sdf"}, []int{10, 1}, 900}}
	// b, _ := json.Marshal(items)
	// w.Write(b)
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
