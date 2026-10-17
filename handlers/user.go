package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"synthori/ediary/m/models"
	"synthori/ediary/m/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h UserHandler) GetDummyUser(w http.ResponseWriter, r *http.Request) {
	WriteSuccess(w, 200, "hellow world", models.DummyUser{
		Username: "Timur",
		Role:     "user",
	})
}

func (h UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, 400, err)
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		WriteError(w, 400, err)
		return
	}

	WriteSuccess(w, 200, "successful get user", user)
}

func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	users, err := h.service.GetUsers(limit, offset)
	if err != nil {
		WriteError(w, 400, err)
		return
	}

	WriteSuccess(w, 200, "all users", users)
}

func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserCreateForm

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteError(w, 400, err)
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			WriteError(w, 400, err)
			return
		}
	}()
	user.Role = "user"

	newUser, err := h.service.CreateUser(user)
	if err != nil {
		WriteError(w, 400, err)
		return
	}

	WriteSuccess(w, 200, "successful create user", newUser)
}

func (h UserHandler) AuthUser(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// http.Error(w, "invalid request body", http.StatusBadRequest)
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.service.LoginUser(req.Username, req.Password)
	if err != nil {
		// http.Error(
		// 	w,
		// 	"invalid username or password",
		// 	http.StatusUnauthorized,
		// )
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	response := models.LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// "error_type", fmt.Sprintf("%T", err)
		slog.Error("failed to encode login response: ", "error_login: ", fmt.Sprintf("%v", err))
	}
}
