package main

import (
	"context"
	"fmt"
	"net/http"

	controllers "github.com/sertraline/messaggio/controllers"
	db "github.com/sertraline/messaggio/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)


func KafkaCtx(next http.Handler) http.Handler {
	// эта мидлварь будет встраивать reader кафки в контекст
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), db.CtxKey, db.Reader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // максимально допустимое значение
	}))

	db.InitDatabase()

	r.Group(func(r chi.Router) {
		r.Route("/messages", func(r chi.Router) {
			r.Use(KafkaCtx)
			r.Get("/{message_id}", controllers.GetMessageByID)
			r.Post("/", controllers.SaveMessage)
		})
		r.Route("/stats", func(r chi.Router) {
			r.Get("/", controllers.GetStats)
		})
		r.Route("/kafka", func(r chi.Router) {
			r.Use(KafkaCtx)
			r.Get("/read", controllers.GetKafkaMessageFromTopic)
		})
	})

	addr := ":3333"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, r)
}
