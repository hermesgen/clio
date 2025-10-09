package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hermesgen/hm"

	"github.com/google/uuid"
)

func (h *APIHandler) GetLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetLayout", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resLayoutName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var layout Layout
	layout, err = h.svc.GetLayout(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetItem, hm.Cap(resLayoutName))
	h.OK(w, msg, layout)
}

func (h *APIHandler) CreateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateLayout", h.Name())

	var layout Layout
	var err error
	err = json.NewDecoder(r.Body).Decode(&layout)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	newLayout := Newlayout(layout.Name, layout.Description, layout.Code)
	newLayout.GenCreateValues()

	err = h.svc.CreateLayout(r.Context(), newLayout)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotCreateResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgCreateItem, hm.Cap(resLayoutName))
	h.Created(w, msg, newLayout)
}

func (h *APIHandler) UpdateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateLayout", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resLayoutName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var layout Layout
	err = json.NewDecoder(r.Body).Decode(&layout)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	updatedLayout := Newlayout(layout.Name, layout.Description, layout.Code)
	updatedLayout.SetID(id, true)
	updatedLayout.GenUpdateValues()

	err = h.svc.UpdateLayout(r.Context(), updatedLayout)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotUpdateResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgUpdateItem, hm.Cap(resLayoutName))
	h.OK(w, msg, updatedLayout)
}

func (h *APIHandler) DeleteLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteLayout", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resLayoutName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteLayout(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotDeleteResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgDeleteItem, hm.Cap(resLayoutName))
	h.OK(w, msg, json.RawMessage("null"))
}

func (h *APIHandler) GetAllLayouts(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllLayouts", h.Name())

	var layouts []Layout
	var err error
	layouts, err = h.svc.GetAllLayouts(r.Context())
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResources, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetAllItems, hm.Cap(resLayoutName))
	h.OK(w, msg, layouts)
}
