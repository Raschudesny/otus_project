package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Raschudesny/otus_project/v1/internal"
	"github.com/Raschudesny/otus_project/v1/storage"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ internal.Repository = (*Storage)(nil)

type Storage struct {
	db         *sqlx.DB
	driverName string
	dsn        string
}

func NewStorage(driverName, dsn string) *Storage {
	return &Storage{driverName: driverName, dsn: dsn}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	if s.db, err = sqlx.Open(s.driverName, s.dsn); err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	s.db.SetMaxOpenConns(20)
	s.db.SetMaxIdleConns(5)
	s.db.SetConnMaxLifetime(time.Minute * 3)

	if err = s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	return nil
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("error during db connection pool closing: %w", err)
	}
	return nil
}

func (s *Storage) AddSlot(ctx context.Context, description string) (string, error) {
	query := "INSERT INTO slots (slot_description) VALUES (:description) RETURNING slot_id"
	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{"description": description})
	if err != nil {
		return "", fmt.Errorf("sql execution error: %w", err)
	}
	defer func() {
		err := rows.Close()
		zap.L().Error("error closing sql rows object", zap.Error(err))
	}()

	var id string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return "", fmt.Errorf("sql AddSlot result parsing error: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("sql AddSlot result parsing error %w", err)
	}
	return id, nil
}

func (s *Storage) GetSlotByID(ctx context.Context, id string) (storage.Slot, error) {
	query := "SELECT (slot_id, slot_description) FROM slots WHERE slot_id = ?"
	row := s.db.QueryRowxContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return storage.Slot{}, fmt.Errorf("sql execution error: %w", err)
	}

	slot := new(storage.Slot)
	err := row.StructScan(slot)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return storage.Slot{}, storage.ErrSlotNotFound
	case err != nil:
		return storage.Slot{}, fmt.Errorf("sql GetSlotById result scan error: %w", err)
	default:
		return *slot, nil
	}
}

func (s *Storage) DeleteSlot(ctx context.Context, id string) error {
	query := "DELETE FROM slots WHERE slot_id = :id"
	res, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("sql delete slot delete operation query error: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error during rows affected by delete checking: %w", err)
	}
	if affected == 0 {
		return storage.ErrSlotNotFound
	}
	return nil
}

func (s *Storage) AddBanner(ctx context.Context, description string) (string, error) {
	query := "INSERT INTO banners (banner_description) VALUES (:description) RETURNING banner_id"
	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{"description": description})
	if err != nil {
		return "", fmt.Errorf("sql execution error: %w", err)
	}
	defer func() {
		err := rows.Close()
		zap.L().Error("error closing sql rows object", zap.Error(err))
	}()

	var id string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return "", fmt.Errorf("sql AddBanner result parsing error: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("sql AddBanner result parsing error %w", err)
	}
	return id, nil
}

func (s *Storage) GetBannerByID(ctx context.Context, id string) (storage.Banner, error) {
	query := "SELECT (banner_id, banner_description) FROM banners WHERE banner_id = ?"
	row := s.db.QueryRowxContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return storage.Banner{}, fmt.Errorf("sql execution error: %w", err)
	}

	banner := new(storage.Banner)
	err := row.StructScan(banner)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return storage.Banner{}, storage.ErrSlotNotFound
	case err != nil:
		return storage.Banner{}, fmt.Errorf("sql GetBannerById result scan error: %w", err)
	default:
		return *banner, nil
	}
}

func (s *Storage) FindBannersBySlot(ctx context.Context, slotID, groupID string) ([]storage.Banner, error) {
	query := `SELECT banner_id, banner_description
			  FROM banners inner join slot_banners sb using (banner_id)
			  WHERE slot_id = '6fe62f82-17fd-4cc1-bb52-efdaab79afef';`
	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"slotID":  slotID,
		"groupID": groupID,
	})
	if err != nil {
		return nil, fmt.Errorf("error during sql execution: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			zap.L().Error("error closing rows object for FindBannersBySlotAndGroup query", zap.Error(err))
		}
	}()

	var banners []storage.Banner
	var banner storage.Banner
	for rows.Next() {
		if err := rows.StructScan(&banner); err != nil {
			return nil, fmt.Errorf("sql FindBannersBySlot result parsing error: %w", err)
		}
		banners = append(banners, banner)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sql FindBannersBySlot result parsing error %w", err)
	}
	return banners, nil
}

func (s *Storage) DeleteBanner(ctx context.Context, id string) error {
	query := "DELETE FROM banners WHERE banner_id = :id"
	res, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("sql delete slot delete operation query error: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error during rows affected by delete checking: %w", err)
	}
	if affected == 0 {
		return storage.ErrBannerNotFound
	}
	return nil
}

func (s *Storage) AddBannerToSlot(ctx context.Context, bannerID, slotID string) error {
	query := "INSERT INTO slot_banners (slot_id, banner_id) VALUES (:slotId, :bannerId)"
	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"slotId":   slotID,
		"bannerId": bannerID,
	})
	if err != nil {
		return fmt.Errorf("error during sql execution: %w", err)
	}
	return nil
}

func (s *Storage) DeleteBannerFromSlot(ctx context.Context, bannerID, slotID string) error {
	query := "DELETE FROM slot_banners WHERE slot_id = :slotId AND banner_id = :bannerId"
	res, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"bannerId": bannerID,
		"slotId":   slotID,
	})
	if err != nil {
		return fmt.Errorf("sql delete slot delete operation query error: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error during rows affected by delete checking: %w", err)
	}
	if affected == 0 {
		return storage.ErrSlotToBannerRelationNotFound
	}
	return nil
}

func (s *Storage) AddGroup(ctx context.Context, description string) (string, error) {
	query := "INSERT INTO social_groups (group_description) VALUES (:description) RETURNING group_id"
	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{"description": description})
	if err != nil {
		return "", fmt.Errorf("sql execution error: %w", err)
	}
	defer func() {
		err := rows.Close()
		zap.L().Error("error closing sql rows object", zap.Error(err))
	}()

	var id string
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return "", fmt.Errorf("sql AddGroup result parsing error: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("sql AddGroup result parsing error %w", err)
	}
	return id, nil
}

func (s *Storage) DeleteGroup(ctx context.Context, id string) error {
	query := "DELETE FROM social_groups WHERE group_id = :id"
	res, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("sql delete slot delete operation query error: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error during rows affected by delete checking: %w", err)
	}
	if affected == 0 {
		return storage.ErrBannerNotFound
	}
	return nil
}

func (s *Storage) Transact(ctx context.Context, work func(*sqlx.Tx) error) error {
	txx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	if err = work(txx); err != nil {
		defer func() {
			// TODO don't have many thoughts about what to do with rollback error ...
			if err := txx.Rollback(); err != nil {
				zap.L().Error("error during transaction rollback", zap.Error(err))
			}
		}()
		return err
	}
	return txx.Commit()
}

func (s *Storage) InitStatForBanner(ctx context.Context, slotID, groupID, bannerID string) error {
	query := `INSERT INTO banner_stats (slot_id, banner_id, group_id, clicks_amount, shows_amount)
 			  VALUES (:slotId, :bannerId, :groupId, 0, 0)`
	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"slotId":   slotID,
		"bannerId": bannerID,
		"groupId":  groupID,
	})
	if err != nil {
		return fmt.Errorf("error during sql execution: %w", err)
	}
	return nil
}

func (s *Storage) PersistClick(ctx context.Context, slotID, groupID, bannerID string) error {
	query := `UPDATE banner_stats
			  SET clicks_amount = clicks_amount + 1
			  WHERE slot_id = :slotId AND group_id = :groupId AND banner_id = :bannerId`
	res, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"bannerId": bannerID,
		"slotId":   slotID,
		"groupId":  groupID,
	})
	if err != nil {
		return fmt.Errorf("error during sql execution: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error during sql rows affected by delete checking: %w", err)
	}
	if affected == 0 {
		return storage.ErrNoStatsFound
	}
	return nil
}

func (s *Storage) PersistShow(ctx context.Context, slotID, groupID, bannerID string) error {
	query := `UPDATE banner_stats 
			  SET shows_amount = shows_amount + 1 
			  WHERE slot_id = :slotId AND group_id = :groupId AND banner_id = :bannerId`
	res, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"bannerId": bannerID,
		"slotId":   slotID,
		"groupId":  groupID,
	})
	if err != nil {
		return fmt.Errorf("error during sql execution: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error during sql rows affected by delete checking: %w", err)
	}
	if affected == 0 {
		return storage.ErrNoStatsFound
	}
	return nil
}

func (s *Storage) GetShowsAmount(ctx context.Context, slotID, groupID, bannerID string) (int, error) {
	query := `SELECT shows_amount 
			  FROM banner_stats
			  WHERE slot_id = :slotId AND group_id = :groupId AND banner_id = :bannerId`
	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"slotId":   slotID,
		"groupId":  groupID,
		"bannerId": bannerID,
	})
	if err != nil {
		return -1, fmt.Errorf("error during sql execution: %w", err)
	}
	defer func() {
		err := rows.Close()
		zap.L().Error("error closing sql rows object", zap.Error(err))
	}()

	var showsAmount int
	for rows.Next() {
		if err := rows.Scan(&showsAmount); err != nil {
			return -1, fmt.Errorf("sql GetShowsCount result parsing error: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return -1, fmt.Errorf("sql AddBanner result parsing error %w", err)
	}
	return showsAmount, nil
}

func (s *Storage) GetClicksAmount(ctx context.Context, slotID, groupID, bannerID string) (int, error) {
	query := `SELECT clicks_amount
		      FROM banner_stats
 			  WHERE slot_id = :slotId AND group_id = :groupId AND banner_id = :bannerId`
	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"slotId":   slotID,
		"groupId":  groupID,
		"bannerId": bannerID,
	})
	if err != nil {
		return -1, fmt.Errorf("error during sql execution: %w", err)
	}
	defer func() {
		err := rows.Close()
		zap.L().Error("error closing sql rows object", zap.Error(err))
	}()

	var clicksAmount int
	for rows.Next() {
		if err := rows.Scan(&clicksAmount); err != nil {
			return -1, fmt.Errorf("sql GetShowsCount result parsing error: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return -1, fmt.Errorf("sql AddBanner result parsing error %w", err)
	}
	return clicksAmount, nil
}

// TODO rewrite on db.QueryRow() function.
func (s *Storage) CountTotalShowsAmount(ctx context.Context, slotID, groupID string) (uint, error) {
	query := "SELECT SUM(shows_amount) FROM banner_stats WHERE slot_id = :slotId AND group_id = :groupId"
	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"slotId":  slotID,
		"groupId": groupID,
	})
	if err != nil {
		return 0, fmt.Errorf("error during sql execution: %w", err)
	}
	var totalShows uint
	for rows.Next() {
		if err := rows.Scan(&totalShows); err != nil {
			return 0, fmt.Errorf("sql CountTotalShowsAmount result parsing error: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("sql CountTotalShowsAmount result parsing error: %w", err)
	}
	return totalShows, nil
}
