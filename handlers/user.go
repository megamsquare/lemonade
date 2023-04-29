package handlers

import (
	"encoding/json"
	"github.com/megamsquare/lemonade/repository/user"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type UserHandler struct {
	UserRepo user.UserRepository
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser user.NewUser
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newUser.Balance = 1000
	newUser.VerificationStatus = false

	err = h.UserRepo.Create(r.Context(), "CreateUser", &newUser)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserRepo.QueryAll(r.Context(), uuid.NewString())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}

func (h *UserHandler) ProcessVerify() {
	h.UserRepo.VerifyProcess()
}
