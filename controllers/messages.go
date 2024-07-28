package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	errResp "github.com/sertraline/messaggio/errors"
	services "github.com/sertraline/messaggio/services"
	validators "github.com/sertraline/messaggio/validators"
)

func SaveMessage(w http.ResponseWriter, r *http.Request) {
	data := &validators.MessageSaveRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
		return
	}

	userData, err := services.SaveMessage(data, r.Context())
	if err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
	}

	render.JSON(w, r, userData)
}

func GetMessageByID(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "message_id")
	if idstr == "" {
		render.Render(w, r, errResp.ErrNotFound)
		return
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
	}

	message, err := services.GetMessageByID(id)
	if err != nil {
		if message.Id == 0 {
			render.Render(w, r, errResp.ErrNotFound)
		} else {
			render.Render(w, r, errResp.ErrInvalidRequest(err))
		}
	} else {
		render.JSON(w, r, message)
	}
}
