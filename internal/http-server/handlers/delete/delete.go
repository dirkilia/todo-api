package delete

import (
	"log/slog"
	"net/http"
	"strconv"
	httpserver "todo/internal/http-server"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type TaskDelete interface {
	DeleteTaskById(id int64) (int64, error)
}

func New(log *slog.Logger, taskDelete TaskDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		task_id_parameter := chi.URLParam(r, "id")
		task_id, err := strconv.Atoi(task_id_parameter)
		if err != nil {
			log.Error("failed to convert task id", err)

			render.JSON(w, r, httpserver.Response{
				Status: "Error",
				Error:  "internal error",
			})

			return
		}

		res, err := taskDelete.DeleteTaskById(int64(task_id))
		if err != nil {
			log.Error("failed to get url", err)

			render.JSON(w, r, httpserver.Response{
				Status: "Error",
				Error:  "internal error",
			})

			return
		}

		render.JSON(w, r, res)
	}
}
