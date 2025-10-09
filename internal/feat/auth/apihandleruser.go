package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hermesgen/hm"
)

// User related API handlers

func (h *APIHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllUsers", h.Name())

	users, err := h.svc.GetUsers(r.Context())
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResources, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetAllItems, resUserNameCap)
	h.OK(w, msg, users)
}

func (h *APIHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetUser", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	user, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := fmt.Sprintf("User with ID %s not found", id)
			h.Err(w, http.StatusNotFound, msg, err)
			return
		}
		msg := fmt.Sprintf(hm.ErrCannotGetResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetItem, resUserNameCap)
	h.OK(w, msg, user)
}

func (h *APIHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateUser", h.Name())

	var form UserForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	newUser := NewUser(form.Username, form.Name, form.Email)

	err = h.svc.CreateUser(r.Context(), &newUser)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotCreateResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	createdUser, err := h.svc.GetUser(r.Context(), newUser.ID)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgCreateItem, resUserNameCap)
	h.Created(w, msg, createdUser)
}

func (h *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateUser", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var form UserForm
	err = json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	updatedUser := NewUser(form.Username, form.Name, form.Email)
	updatedUser.SetID(id)

	err = h.svc.UpdateUser(r.Context(), &updatedUser)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotUpdateResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	finalUser, err := h.svc.GetUser(r.Context(), updatedUser.ID)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgUpdateItem, resUserNameCap)
	h.OK(w, msg, finalUser)
}

func (h *APIHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteUser", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteUser(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotDeleteResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgDeleteItem, resUserNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}
