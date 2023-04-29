package handlers

import (
	"context"
	"encoding/json"
	"github.com/megamsquare/lemonade/repository/transaction"
	"github.com/megamsquare/lemonade/repository/user"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type TransactionHandler struct {
	TransactionRepo transaction.TransactionRepository
	UserRepo        user.UserRepository
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var newTransaction transaction.NewTransaction
	err := json.NewDecoder(r.Body).Decode(&newTransaction)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sender, err := h.UserRepo.QueryByID(r.Context(), uuid.NewString(), newTransaction.SenderID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiver, err := h.UserRepo.QueryByID(r.Context(), uuid.NewString(), newTransaction.ReceiverID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if sender.Balance < newTransaction.Amount {
		log.Println("insufficient Funds")
		http.Error(w, "insufficient Funds", http.StatusBadRequest)
		return
	}

	isSenderVerified, err := h.UserRepo.IsVerified(r.Context(), "Is User Verified", sender.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isReceiverVerified, err := h.UserRepo.IsVerified(r.Context(), "Is User Verified", receiver.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isReceiverVerified {
		log.Println(err)
		http.Error(w, "receiver is not verified", http.StatusBadRequest)
		return
	}

	if !isSenderVerified {
		log.Println(err)
		http.Error(w, "sender is not verified", http.StatusBadRequest)
		return
	}

	err = h.TransactionRepo.Create(r.Context(), "CreateTransaction", &newTransaction)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTransaction)
}

// func (h *TransactionHandler) ProcessTransactions(numWorkers int) {
// 	for i := 0; i < numWorkers; i++ {
// 		go func() {
// 			transactions, err := h.TransactionRepo.QueryAll("Get all transactions")
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 			for _, transaction := range transactions {
// 				// Deduct amount from sender and add to rece
// 				sender, _ := h.UserRepo.QueryByID(context.Background(), "", transaction.SenderID)
// 				receiver, _ := h.UserRepo.QueryByID(context.Background(), "", transaction.ReceiverID)
// 				sender.Balance -= transaction.Amount
// 				receiver.Balance += transaction.Amount
// 				err := h.UserRepo.Update(context.Background(), "", sender)
// 				if err != nil {
// 					log.Println(err)
// 				}
// 				err = h.UserRepo.Update(context.Background(), "", receiver)
// 				if err != nil {
// 					log.Println(err)
// 				}
// 			}
// 		}()
// 	}
// }

func (h *TransactionHandler) ProcessTransactions(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				go func() {
					transactions, err := h.TransactionRepo.QueryAll("Get all transactions")
					if err != nil {
						// log.Println(err)
						return
					}
					for _, transaction := range transactions {
						// Deduct amount from sender and add to rece
						sender, _ := h.UserRepo.QueryByID(context.Background(), "", transaction.SenderID)
						receiver, _ := h.UserRepo.QueryByID(context.Background(), "", transaction.ReceiverID)
						sender.Balance -= transaction.Amount
						receiver.Balance += transaction.Amount
						err := h.UserRepo.Update(context.Background(), "", sender)
						if err != nil {
							log.Println(err)
						}
						err = h.UserRepo.Update(context.Background(), "", receiver)
						if err != nil {
							log.Println(err)
						}
					}
				}()
			}
		}

	}
}
