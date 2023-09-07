// Package repository contains CRUD operations
package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/profile/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

// ProfileRepository represents a repository level
type ProfileRepository struct {
	pool *pgxpool.Pool
}

// NewProfileRepository creates a new ProfileRepository
func NewProfileRepository(pool *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{pool: pool}
}

func (db *ProfileRepository) GetIDByLoginPassword(ctx context.Context, login string) (ID uuid.UUID, pass string, err error) {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: "repeatable read"})
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("BeginTx: %w", err)
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				logrus.Errorf("Rollback: %v", err)
				return
			}
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				logrus.Errorf("Commit: %v", err)
				return
			}
		}
	}()

	err = tx.QueryRow(ctx, "SELECT id, password FROM profile.profile WHERE login=$1", login).Scan(&ID, &pass)
	if err != nil {
		logrus.Errorf("QueryRow: %v", err)
		return uuid.Nil, "", fmt.Errorf("QueryRow: %w", err)
	}
	return ID, pass, nil
}

// GetProfileByID function returns a profile with the given ID
func (db *ProfileRepository) GetProfileByID(ctx context.Context, id uuid.UUID) (*model.Profile, error) {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: "repeatable read"})
	if err != nil {
		return nil, fmt.Errorf("BeginTx: %w", err)
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				logrus.Errorf("Rollback: %v", err)
				return
			}
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				logrus.Errorf("Commit: %v", err)
				return
			}
		}
	}()
	profile := &model.Profile{}
	err = tx.QueryRow(ctx, "SELECT id, login, password, refresh_token FROM profile.profile WHERE id = $1", id).Scan(&profile.ID, &profile.Login, &profile.Password, &profile.RefreshToken)
	if err != nil {
		logrus.Errorf("QueryRow: %v", err)
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return profile, nil
}

// CreateProfile function creates a new profile in database
func (db *ProfileRepository) CreateProfile(ctx context.Context, profile *model.Profile) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: "repeatable read"})
	if err != nil {
		return fmt.Errorf("BeginTx: %w", err)
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				logrus.Errorf("Rollback: %v", err)
				return
			}
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				logrus.Errorf("Commit: %v", err)
				return
			}
		}
	}()
	_, err = db.pool.Exec(ctx, "INSERT INTO profile.profile (id, login, password, refresh_token) VALUES ($1, $2, $3, $4)", profile.ID, profile.Login, profile.Password, profile.RefreshToken)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

// UpdateProfile function updates the profile information in database
func (db *ProfileRepository) UpdateProfile(ctx context.Context, profile *model.Profile) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: "repeatable read"})
	if err != nil {
		return fmt.Errorf("BeginTx: %w", err)
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				logrus.Errorf("Rollback: %v", err)
				return
			}
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				logrus.Errorf("Commit: %v", err)
				return
			}
		}
	}()
	_, err = tx.Exec(
		ctx,
		"UPDATE profile.profile SET login=$1, password=$2, refresh_token=$3 WHERE id=$4",
		profile.Login, profile.Password, profile.RefreshToken, profile.ID,
	)
	if err != nil {
		logrus.Errorf("Exec: %v", err)
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (db *ProfileRepository) DeleteProfileByID(ctx context.Context, id string) error {
	tag, err := db.pool.Exec(ctx, "DELETE FROM profile.profile WHERE id=$1", id)
	if err != nil || tag.RowsAffected() == 0 {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
