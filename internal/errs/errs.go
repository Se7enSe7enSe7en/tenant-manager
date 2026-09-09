package errs

import (
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

func Http(w http.ResponseWriter, r *http.Request, err error, status int) {
	// respond to frontend with error
	http.Error(w, err.Error(), status)

	// get line trace from Program Counter (pc)
	var pc uintptr
	var pcs [1]uintptr
	if runtime.Callers(2, pcs[:]) == 1 {
		pc = pcs[0]
	}

	record := slog.NewRecord(time.Now(), slog.LevelError, err.Error(), pc)

	// // Attach request context as structured fields
	// record.AddAttrs(
	// 	slog.String("method", r.Method),
	// 	slog.String("path", r.URL.Path),
	// 	slog.Int("status", status),
	// )

	// ignore err, show error in terminal
	_ = slog.Default().Handler().Handle(r.Context(), record)
}
