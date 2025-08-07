package main

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

//go:embed template.gohtml
var templateText string

type ListHandler struct {
	repo *EntryRepo
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		parsedTemplate     *template.Template
		page               = 1
		currentPageEntries []Entry
		totalEntries       int
	)
	if t, err := template.New("template.gohtml").Parse(templateText); err != nil {
		log.Print(err)
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		parsedTemplate = t
	}

	var pageParam = strings.TrimSpace(r.FormValue("page"))
	if pageParam != "" {
		if i, err := strconv.Atoi(pageParam); err != nil {
			log.Print(err)
			http.Error(w, "Error parsing page parameter: "+err.Error(), http.StatusBadRequest)
			return
		} else {
			page = i
		}
	}

	if page < 1 {
		http.Error(w, fmt.Sprintf("Invalid page number %d", page), http.StatusBadRequest)
		return
	}

	if entries, err := h.repo.ListEntries(page); err != nil {
		log.Print(err)
		http.Error(w, "Error loading entries: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		currentPageEntries = entries
	}

	if count, err := h.repo.CountEntries(); err != nil {
		log.Print(err)
		http.Error(w, "Error counting entries: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		totalEntries = count
	}
	var pageCount = int(math.Ceil(float64(totalEntries / MaxEntriesPerPage)))
	var pageNumbers = make([]int, 0, pageCount)
	for i := 1; i <= pageCount; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	var templateData = map[string]any{
		"entries":      currentPageEntries,
		"totalEntries": totalEntries,
		"page":         page,
		"pageNumbers":  pageNumbers,
	}

	if err := parsedTemplate.ExecuteTemplate(w, "template.gohtml", templateData); err != nil {
		log.Print(err)
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
