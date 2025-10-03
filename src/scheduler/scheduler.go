package scheduler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Runner is the function executed on schedule
type Runner func(ctx context.Context) error

// Scheduler runs a Runner at fixed intervals while respecting an advisory lock
type Scheduler struct {
	DB       *sql.DB
	LockKey  int64 // advisory lock key
	Interval time.Duration
}

// Run starts the scheduler loop until ctx is cancelled
func (s *Scheduler) Run(ctx context.Context, run Runner) error {
	if s.DB == nil {
		return errors.New("db required for scheduler")
	}
	if s.Interval <= 0 {
		s.Interval = time.Minute
	}

	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	// Run immediately once
	if err := s.tryRunOnce(ctx, run); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.tryRunOnce(ctx, run); err != nil {
				// continue loop but surface error
				fmt.Printf("scheduler run error: %v\n", err)
			}
		}
	}
}

// tryRunOnce attempts to acquire advisory lock and run the job
func (s *Scheduler) tryRunOnce(ctx context.Context, run Runner) error {
	// acquire advisory lock
	var got bool
	row := s.DB.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", s.LockKey)
	if err := row.Scan(&got); err != nil {
		return fmt.Errorf("advisory lock check: %w", err)
	}
	if !got {
		// another instance is running
		return nil
	}
	defer func() {
		// release lock, best effort
		s.DB.ExecContext(context.Background(), "SELECT pg_advisory_unlock($1)", s.LockKey)
	}()

	// run the function
	return run(ctx)
}
