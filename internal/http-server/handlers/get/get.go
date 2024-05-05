package get

import (
	"log/slog"
	"net/http"
	httpserver "todo/internal/http-server"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type TaskGetter interface {
	GetTasks() ([]httpserver.Task, error)
}

func New(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		res, err := taskGetter.GetTasks()
		if err != nil {
			log.Error("failed to get url", err)

			render.JSON(w, r, httpserver.Response{
				Status: "Error",
				Error:  "internal error",
			})

			return
		}

		log.Info("got tasks", slog.Any("task", res))

		render.JSON(w, r, res)
	}
}
