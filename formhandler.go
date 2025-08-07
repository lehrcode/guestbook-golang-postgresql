package main

import (
	"net/http"
	"strings"
)

type FormHandler struct {
	repo *EntryRepo
}

func (h *FormHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		name    = strings.TrimSpace(r.PostFormValue("name"))
		email   = strings.TrimSpace(r.PostFormValue("email"))
		message = strings.TrimSpace(r.PostFormValue("message"))
	)

	if name == "" || email == "" || message == "" {
		http.Error(w, "name, email and message are required", http.StatusBadRequest)
		return
	}

	if err := h.repo.AddEntry(name, email, message); err != nil {
		http.Error(w, "Error creating entry: "+err.Error(), http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
