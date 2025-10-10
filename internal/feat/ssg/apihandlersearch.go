package ssg

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hermesgen/hm"
)

func (h *APIHandler) SearchContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling SearchContent", h.Name())

	query := r.URL.Query()
	searchQuery := query.Get("search")
	pageStr := query.Get("page")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	const itemsPerPage = 25
	offset := (page - 1) * itemsPerPage

	contents, totalCount, err := h.svc.GetContentWithPaginationAndSearch(r.Context(), offset, itemsPerPage, searchQuery)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResources, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	totalPages := (totalCount + itemsPerPage - 1) / itemsPerPage

	response := struct {
		Contents   []Content `json:"contents"`
		Page       int       `json:"page"`
		TotalPages int       `json:"total_pages"`
		TotalCount int       `json:"total_count"`
		Search     string    `json:"search"`
	}{
		Contents:   contents,
		Page:       page,
		TotalPages: totalPages,
		TotalCount: totalCount,
		Search:     searchQuery,
	}

	msg := fmt.Sprintf("Search results for '%s'", searchQuery)
	if searchQuery == "" {
		msg = fmt.Sprintf(hm.MsgGetAllItems, hm.Cap(resContentName))
	}

	h.OK(w, msg, response)
}
