package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	httpserver "todo/internal/http-server"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Title      string `json:"title"`
	IsFinished bool   `json:"is_finished"`
}

type TaskSaver interface {
	SaveTask(task string, is_finished bool) (int64, error)
}

func New(log *slog.Logger, tasksaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, httpserver.Response{
				Status: "Error",
				Error:  "empty request",
			})

			return
		}
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, httpserver.Response{
				Status: "Error",
				Error:  "failed to decode request",
			})

			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		if req.Title == "" {
			render.JSON(w, r, httpserver.Response{
				Status: "Error",
				Error:  "title is empty",
			})

			return
		}

		id, err := tasksaver.SaveTask(req.Title, req.IsFinished)
		if err != nil {
			log.Error("failed to add url", err)

			render.JSON(w, r, httpserver.Response{
				Status: "Error",
				Error:  "failed to decode request",
			})

			return
		}
		log.Info("task added", slog.Int64("id", id))
		render.JSON(w, r, httpserver.Response{
			Status: "OK",
		})
	}
}
