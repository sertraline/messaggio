package controllers

import (
	"net/http"

	"github.com/go-chi/render"
	errResp "github.com/sertraline/messaggio/errors"
	services "github.com/sertraline/messaggio/services"
)

func GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := services.GetStats()
	if err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
	} else {
		render.JSON(w, r, stats)
	}
}