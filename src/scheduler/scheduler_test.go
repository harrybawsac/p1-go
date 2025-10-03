package scheduler

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTryRunOnce_LockAcquired(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	s := &Scheduler{DB: db, LockKey: 1}

	// expect pg_try_advisory_lock to be called and return true
	mock.ExpectQuery("SELECT pg_try_advisory_lock\\(\\$1\\)").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"}).AddRow(true))
	mock.ExpectExec("SELECT pg_advisory_unlock\\(\\$1\\)").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	called := false
	err = s.tryRunOnce(context.Background(), func(ctx context.Context) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("tryRunOnce: %v", err)
	}
	if !called {
		t.Fatalf("expected runner to be called")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestTryRunOnce_LockNotAcquired(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	s := &Scheduler{DB: db, LockKey: 1}

	// expect pg_try_advisory_lock to be called and return false
	mock.ExpectQuery("SELECT pg_try_advisory_lock\\(\\$1\\)").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"}).AddRow(false))

	called := false
	err = s.tryRunOnce(context.Background(), func(ctx context.Context) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("tryRunOnce: %v", err)
	}
	if called {
		t.Fatalf("expected runner NOT to be called when lock not acquired")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
