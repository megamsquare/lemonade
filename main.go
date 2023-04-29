package main

import (
	"github.com/megamsquare/lemonade/handlers"
	"github.com/megamsquare/lemonade/repository/user"
	"github.com/megamsquare/lemonade/repository/transaction"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	userHandler := handlers.UserHandler{UserRepo: user.NewUserMemoryStore()}
	transactionHandler := handlers.TransactionHandler{TransactionRepo: transaction.NewTransactionMemoryStore()}
	r.HandleFunc("/user", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/user", userHandler.GetAllUser).Methods("GET")
	r.HandleFunc("/transaction", transactionHandler.CreateTransaction).Methods("POST")
	go userHandler.ProcessVerify()
	go transactionHandler.ProcessTransactions(2)
	log.Println("Starting Server at port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

// type User struct {
// 	ID int
// 	Username string
// 	Balance string
// }

// type HashTable struct {
// 	table map[int]*User
// }

// func main() {}

// func (h *HashTable)AddUser(w *http.Request, r *http.Response) {
// 	var u User
// 	res := r.Body
// 	err := json.NewDecoder(res).Decode(&u)
// 	if err != nil {

// 	}
// 	h.table[u.ID] = u
// }

// func (h *HashTable) GetUserByID(w *http.Request, r *http.Response) (*User, error) {
// 	res := r.Body
// 	u, ok := h.table[res.id]
// 	if !ok {
// 		return nil, fmt.Errorf()
// 	}
// }
