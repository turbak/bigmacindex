package app

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/turbak/bigmacindex/internal/domain/link"
	"github.com/turbak/bigmacindex/internal/domain/price"
)

//go:embed templates
var templates embed.FS

type LinksAdderRemoverLister interface {
	AddLink(ctx context.Context, link link.LinkDescription) (link.LinkDescription, error)
	ListLinks(ctx context.Context) ([]link.LinkDescription, error)
	DeleteLink(ctx context.Context, linkID link.ID) error
}

type PriceLister interface {
	ListPrices(ctx context.Context) ([]price.PriceRecord, error)
}

type App struct {
	linkRepo  LinksAdderRemoverLister
	priceRepo PriceLister
}

func NewApp(
	linkRepo LinksAdderRemoverLister,
	priceRepo PriceLister,
) *App {
	return &App{
		linkRepo:  linkRepo,
		priceRepo: priceRepo,
	}
}

func (a *App) SetupRoutes(_ context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /links", a.GetLinks())
	mux.HandleFunc("POST /links", a.CreateLink())
	mux.HandleFunc("DELETE /links/{id}", a.DeleteLink)

	mux.HandleFunc("GET /prices", a.GetPrices)

	return http.ListenAndServe(":8080", mux)
}

func (a *App) GetLinks() func(rw http.ResponseWriter, req *http.Request) {
	templ := template.Must(template.ParseFS(templates, "templates/links.html"))

	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		links, err := a.linkRepo.ListLinks(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		data := struct {
			Links []link.LinkDescription
		}{
			Links: links,
		}

		err = templ.Execute(rw, data)
		if err != nil {
			log.Println(err)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}

func (a *App) CreateLink() func(rw http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFS(templates, "templates/links.html"))

	return func(rw http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(rw, "Invalid request", http.StatusBadRequest)
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
			http.Error(rw, "Failed to save link", http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(rw, "link-row", createdLink)
		if err != nil {
			log.Println("Error rendering row:", err)
		}
	}
}

func (a *App) DeleteLink(rw http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	if idStr == "" {
		http.Error(rw, "Missing ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid ID format: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = a.linkRepo.DeleteLink(req.Context(), link.ID(id))
	if err != nil {
		http.Error(rw, "Failed to delete: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (a *App) GetPrices(rw http.ResponseWriter, req *http.Request) {
}
