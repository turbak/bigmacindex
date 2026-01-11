package app

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/turbak/bigmacindex/internal/domain/price"
)

//go:embed templates
var templates embed.FS

type PriceLister interface {
	ListPrices(ctx context.Context) ([]price.PriceRecord, error)
}

type App struct {
	linksRoutes *LinksRoutes
	priceRepo   PriceLister
}

func NewApp(
	linksRoutes *LinksRoutes,
	priceRepo PriceLister,
) *App {
	return &App{
		linksRoutes: linksRoutes,
		priceRepo:   priceRepo,
	}
}

func (a *App) SetupRoutes(_ context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /links", a.linksRoutes.GetLinks())
	mux.HandleFunc("POST /links", a.linksRoutes.CreateLink())
	mux.HandleFunc("DELETE /links/{id}", a.linksRoutes.DeleteLink)
	mux.HandleFunc("GET /links/{id}/edit", a.linksRoutes.EditLink())
	mux.HandleFunc("PUT /links/{id}", a.linksRoutes.UpdateLink())
	mux.HandleFunc("GET /links/{id}", a.linksRoutes.GetLink())

	mux.HandleFunc("GET /prices", a.GetPrices)

	return http.ListenAndServe(":8080", mux)
}

func (a *App) GetPrices(rw http.ResponseWriter, req *http.Request) {
}

var errorTempl = template.Must(template.ParseFS(templates, "templates/error-toast.html"))

func renderError(rw http.ResponseWriter, err error, httpStatusCode int) {
	renderErr := errorTempl.Execute(rw, err)
	if renderErr != nil {
		log.Println(renderErr)
	}

	rw.WriteHeader(httpStatusCode)
}
