package controllers

import (
	"net/http"

	"github.com/go-chi/render"
	errResp "github.com/sertraline/messaggio/errors"
	services "github.com/sertraline/messaggio/services"
)

func GetKafkaMessageFromTopic(w http.ResponseWriter, r *http.Request) {
	message, err := services.GetKafkaMessageFromTopic(r.Context())
	if err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
	} else {
		render.JSON(w, r, message)
	}
}
