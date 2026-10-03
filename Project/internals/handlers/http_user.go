package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"project/helper"
	"project/internals/service"
	"project/internals/storage"
)

type HttpService struct {
	service service.UserService
}

func NewHttpService(serive service.UserService) *HttpService {
	return &HttpService{service: serive}
}

func (h *HttpService) HttpCreate(w http.ResponseWriter, r *http.Request) {

	var DTO storage.UserDTO

	if err := json.NewDecoder(r.Body).Decode(&DTO); err != nil {
		helper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(
		DTO.FirstName,
		DTO.LastName,
		DTO.Phone,
		DTO.CarDTO)
	if err != nil {
		helper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		helper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

}

func (h *HttpService) HttpGet(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	user, err := h.service.GetUserById(id)
	if err != nil {
		helper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		helper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

}

func (h *HttpService) HttpGetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		helper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		helper.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *HttpService) HttpUpdate(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	var DTO storage.UserDTO

	if err := json.NewDecoder(r.Body).Decode(&DTO); err != nil {
		helper.WriteError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUser(
		id,
		DTO.FirstName,
		DTO.LastName,
		DTO.Phone,
		DTO.CarDTO)

	if err != nil {
		if errors.Is(err, helper.UserNotFound) ||
			errors.Is(err, helper.CarNotFound) {
			helper.WriteError(w, err, http.StatusNotFound)
			return

		}

		helper.WriteError(w, err, http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		helper.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *HttpService) HttpDelete(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	err := h.service.DeleteUser(id)
	if err != nil {
		if errors.Is(err, helper.UserNotFound) {
			helper.WriteError(w, err, http.StatusNotFound)
			return
		}

		helper.WriteError(w, err, http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)

}
