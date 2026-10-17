package handlers

import (
	// "database/sql"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"synthori/ediary/m/services"
	"time"

	"github.com/go-chi/chi/v5"
)

type DiaryHandler struct {
	service services.DiaryService
}

func NewDiaryHandler(service services.DiaryService) *DiaryHandler {
	return &DiaryHandler{service: service}
}

func (h DiaryHandler) GetDayByID(w http.ResponseWriter, r *http.Request) {
	// diary ID
	idStr := chi.URLParam(r, "id")
	diaryID, err := strconv.Atoi(idStr)
	if err != nil || diaryID <= 0 {
		WriteError(w, http.StatusBadRequest, errors.New("invalid diary id"))
		return
	}

	// userID
	cookie, err := r.Cookie("user_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			WriteError(w, http.StatusUnauthorized, errors.New("user_id cookie is required"))
			return
		}
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil || userID <= 0 {
		WriteError(w, http.StatusBadRequest, errors.New("invalid user_id cookie"))
		return
	}

	fmt.Println("diaryID:", diaryID)
	fmt.Println("userID:", userID)

	day, err := h.service.GetDayByID(userID, diaryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, errors.New("diary entry not found"))
			return
		}
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteSuccess(w, http.StatusOK, "successful get day", day)
}

func (h *DiaryHandler) GetAllWeek(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		WriteError(w, http.StatusUnauthorized, errors.New("user_id cookie is required"))
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil || userID <= 0 {
		WriteError(w, http.StatusBadRequest, errors.New("invalid user_id cookie"))
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		WriteError(w, http.StatusBadRequest, errors.New("date is required"))
		return
	}

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		WriteError(
			w,
			http.StatusBadRequest,
			errors.New("date must have format YYYY-MM-DD"),
		)
		return
	}

	week, err := h.service.GetAllWeek(userID, date)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteSuccess(w, http.StatusOK, "successful get week", week)
}
