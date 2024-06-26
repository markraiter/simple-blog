package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/config"
	"github.com/markraiter/simple-blog/internal/app/api/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type AuthService interface {
	Auth
}

type PostService interface {
	PostSaver
	PostProvider
	PostProcessor
}

type CommentService interface {
	CommentSaver
}

type Handler struct {
	Healthcheck
	AuthHandler
	PostHandler
	CommentHandler
}

// The response struct is used to send a message back to the client.
type response struct {
	Message string `json:"message"`
}

func New(
	l *slog.Logger,
	v *validator.Validate,
	a AuthService,
	p PostService,
	c CommentService,
) *Handler {
	return &Handler{
		Healthcheck{log: l},
		AuthHandler{
			log:      l,
			validate: v,
			service:  a,
		},
		PostHandler{
			log:       l,
			validate:  v,
			saver:     p,
			provider:  p,
			processor: p,
		},
		CommentHandler{
			log:       l,
			validate:  v,
			saver:     c,
			provider:  nil,
			processor: nil,
		},
	}
}

func (h *Handler) Router(ctx context.Context, cfg config.Config, log *slog.Logger) http.Handler {
	m := http.NewServeMux()

	basicAuth := middleware.BasicAuth(cfg.Auth, log)

	m.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	m.Handle("GET /health", h.APIHealth())
	{
		m.Handle("POST /api/auth/register", h.RegisterUser(ctx))
		m.Handle("POST /api/auth/login", h.Login(ctx, cfg.Auth))
	}

	{
		m.Handle("POST /api/posts", basicAuth(h.CreatePost(ctx)))
		m.Handle("GET /api/posts", h.Posts(ctx))
		m.Handle("GET /api/posts/{id}", h.Post(ctx))
		m.Handle("PUT /api/posts/{id}", basicAuth(h.UpdatePost(ctx)))
		m.Handle("DELETE /api/posts/{id}", basicAuth(h.DeletePost(ctx)))
	}

	{
		m.Handle("POST /api/comments", basicAuth(h.CreateComment(ctx)))
		// m.Handle("GET /api/comments", h.Comments(ctx))
		// m.Handle("GET /api/comments/{id}", h.Comment(ctx))
		// m.Handle("PUT /api/comments/{id}", basicAuth(h.UpdateComment(ctx)))
		// m.Handle("DELETE /api/comments/{id}", basicAuth(h.DeleteComment(ctx)))
	}

	return m
}
