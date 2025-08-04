package adapters

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rossi1/smart-pack/domain"
)

type SmartPackEntity struct {
	ID        int        `pg:"id,pk,auto_increment"`
	Size      int        `pg:"size,unique,notnull"`
	CreatedAt time.Time  `pg:"created_at,default:now()"`
	DeletedAt *time.Time `pg:"deleted_at"` // pointer to allow NULL
}

type SmartPackRepository struct {
	db *pgx.Conn
}

func NewSmartPackRepository(db *pgx.Conn) *SmartPackRepository {
	return &SmartPackRepository{db: db}
}

func (r *SmartPackRepository) GetPackSizes(ctx context.Context) ([]domain.SmartPack, error) {
	rows, err := r.db.Query(ctx, "SELECT size FROM smartpack WHERE deleted_at IS NULL ORDER BY size DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sizes []domain.SmartPack
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, err
		}
		sizes = append(sizes, domain.SmartPack{Size: size})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sizes, nil
}

func (r *SmartPackRepository) SetPackSizes(ctx context.Context, sizes []domain.SmartPack) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	now := time.Now()

	// Mark all existing packs as deleted
	_, err = tx.Exec(ctx, "UPDATE smartpack SET deleted_at = $1 WHERE deleted_at IS NULL", now)
	if err != nil {
		return err
	}

	// Insert new pack sizes
	for _, size := range sizes {
		_, err = tx.Exec(ctx, "INSERT INTO smartpack (size) VALUES ($1)", size.Size)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
