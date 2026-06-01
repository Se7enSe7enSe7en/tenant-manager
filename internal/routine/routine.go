package routine

import (
	"context"
	"time"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
)

func DeleteExpiredSessions(ctx context.Context, queries *repo.Queries) {
	run := func() {
		// delete expired sessions
		if err := queries.DeleteExpiredSessions(ctx); err != nil {
			logger.Error("DeleteExpiredSessions(): ", err)
		}
	}
	run() // run once on startup

	// run every 24hrs.
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			run()
		case <-ctx.Done():
			return
		}
	}
}
