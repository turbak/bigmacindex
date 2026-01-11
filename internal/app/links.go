package app

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/turbak/bigmacindex/internal/domain/link"
)

type LinksCRUDer interface {
	AddLink(ctx context.Context, link link.LinkDescription) (link.LinkDescription, error)
	ListLinks(ctx context.Context) ([]link.LinkDescription, error)
	DeleteLink(ctx context.Context, linkID link.ID) error
	UpdateLink(ctx context.Context, linkDesc link.LinkDescription) (link.LinkDescription, error)
	GetLinkByID(ctx context.Context, linkID link.ID) (link.LinkDescription, error)
}

type LinksRoutes struct {
	linkRepo LinksCRUDer
}

func NewLinksRoutes(linkRepo LinksCRUDer) *LinksRoutes {
	return &LinksRoutes{
		linkRepo: linkRepo,
	}
}

func (a *LinksRoutes) GetLinks() func(rw http.ResponseWriter, req *http.Request) {
	templ := template.Must(template.ParseFS(templates, "templates/links.html"))

	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		links, err := a.linkRepo.ListLinks(ctx)
		if err != nil {
			renderError(rw, fmt.Errorf("failed to list links: %w", err), http.StatusInternalServerError)
			return
		}

		data := struct {
			Links []link.LinkDescription
		}{
			Links: links,
		}

		err = templ.Execute(rw, data)
		if err != nil {
			renderError(rw, fmt.Errorf("failed to render links: %w", err), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}

func (a *LinksRoutes) CreateLink() func(rw http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFS(templates, "templates/links.html"))

	return func(rw http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			renderError(rw, fmt.Errorf("failed to parse form: %w", err), http.StatusBadRequest)
			return
		}

		newLink := link.LinkDescription{
			ProductName:   req.FormValue("product_name"),
			URL:           req.FormValue("url"),
			LinkType:      link.LinkType(req.FormValue("link_type")),
			PriceSelector: req.FormValue("price_selector"),
			CountryCode:   req.FormValue("country_code"),
		}

		createdLink, err := a.linkRepo.AddLink(req.Context(), newLink)
		if err != nil {
			renderError(rw, fmt.Errorf("failed to save link: %w", err), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(rw, "link-row", createdLink)
		if err != nil {
			renderError(rw, fmt.Errorf("failed to render link: %w", err), http.StatusInternalServerError)
		}
	}
}

func (a *LinksRoutes) DeleteLink(rw http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	if idStr == "" {
		renderError(rw, fmt.Errorf("missing ID"), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(rw, fmt.Errorf("invalid ID format: %w", err), http.StatusBadRequest)
		return
	}

	err = a.linkRepo.DeleteLink(req.Context(), link.ID(id))
	if err != nil {
		renderError(rw, fmt.Errorf("failed to delete: %w", err.Error()), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (a *LinksRoutes) UpdateLink() func(rw http.ResponseWriter, req *http.Request) {
	templ := template.Must(template.ParseFS(templates, "templates/links.html"))

	return func(rw http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		if idStr == "" {
			renderError(rw, fmt.Errorf("missing ID"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			renderError(rw, fmt.Errorf("invalid ID format: %w", err), http.StatusBadRequest)
			return
		}

		if err := req.ParseForm(); err != nil {
			renderError(rw, fmt.Errorf("failed to parse form: %w", err), http.StatusBadRequest)
			return
		}

		updatedLink := link.LinkDescription{
			ID:            link.ID(id),
			ProductName:   req.FormValue("product_name"),
			URL:           req.FormValue("url"),
			LinkType:      link.LinkType(req.FormValue("link_type")),
			PriceSelector: req.FormValue("price_selector"),
			CountryCode:   req.FormValue("country_code"),
		}

		_, err = a.linkRepo.UpdateLink(req.Context(), updatedLink)
		if err != nil {
			renderError(rw, fmt.Errorf("failed to update link: %w", err), http.StatusInternalServerError)
			return
		}

		err = templ.ExecuteTemplate(rw, "link-row", updatedLink)
		if err != nil {
			renderError(rw, fmt.Errorf("failed to render link: %w", err), http.StatusInternalServerError)
			return
		}
	}
}

func (a *LinksRoutes) EditLink() func(rw http.ResponseWriter, req *http.Request) {
	templ := template.Must(template.ParseFS(templates, "templates/links.html"))

	return func(rw http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		if idStr == "" {
			renderError(rw, fmt.Errorf("missing ID"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			renderError(rw, fmt.Errorf("invalid ID format: %w", err), http.StatusBadRequest)
			return
		}

		links, err := a.linkRepo.GetLinkByID(req.Context(), link.ID(id))
		if err != nil {
			renderError(rw, fmt.Errorf("failed to get link: %w", err), http.StatusInternalServerError)
			return
		}

		err = templ.ExecuteTemplate(rw, "link-row-edit", links)
	}
}

func (a *LinksRoutes) GetLink() func(rw http.ResponseWriter, req *http.Request) {
	templ := template.Must(template.ParseFS(templates, "templates/links.html"))
	return func(rw http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		if idStr == "" {
			renderError(rw, fmt.Errorf("missing ID"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			renderError(rw, fmt.Errorf("invalid ID format: %w", err), http.StatusBadRequest)
			return
		}

		linkDesc, err := a.linkRepo.GetLinkByID(req.Context(), link.ID(id))
		if err != nil {
			renderError(rw, fmt.Errorf("failed to get link: %w", err), http.StatusInternalServerError)
			return
		}

		err = templ.ExecuteTemplate(rw, "link-row", linkDesc)
		if err != nil {
			renderError(rw, fmt.Errorf("failed to render link: %w", err), http.StatusInternalServerError)
			return
		}
	}
}
