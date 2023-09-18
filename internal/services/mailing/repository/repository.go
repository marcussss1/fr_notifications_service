package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"gitlab.com/fr5270937/notifications_service/internal/model"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/mailing"
)

type repository struct {
	db *sqlx.DB
}

func NewMailingRepository(db *sqlx.DB) desc.Repository {
	return repository{db: db}
}

func (r repository) CreateMailing(ctx context.Context, req model.CreateMailingRequest) (model.Mailing, error) {
	timeNow := time.Now()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return model.Mailing{}, err
	}

	mailing, err := r.createMailing(ctx, req)
	if err != nil {
		_ = tx.Rollback()
		return model.Mailing{}, err
	}
	// если рассылку нужно запустить в данный момент
	if timeNow.String() > req.CreatedAt.String() && timeNow.String() < req.FinishedAt.String() {
		err = tx.Commit()
		if err != nil {
			return model.Mailing{}, err
		}

		return mailing, err
	}

	err = r.createPendingMailing(ctx, mailing)
	if err != nil {
		_ = tx.Rollback()
		return model.Mailing{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Mailing{}, err
	}

	return mailing, nil
}

func (r repository) createMailing(ctx context.Context, req model.CreateMailingRequest) (model.Mailing, error) {
	var mailing model.Mailing

	err := r.db.QueryRowContext(ctx, `INSERT INTO mailing (created_at, text, filter, finished_at) VALUES ($1, $2, $3, $4) RETURNING *`,
		req.CreatedAt, req.Text, req.Filter, req.FinishedAt).
		Scan(&mailing.ID, &mailing.CreatedAt, &mailing.Text, &mailing.Filter, &mailing.FinishedAt)
	if err != nil {
		return model.Mailing{}, err
	}

	return mailing, nil
}

func (r repository) createPendingMailing(ctx context.Context, mailing model.Mailing) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO pending_mailing (mailing_id, status) VALUES ($1, $2)", mailing.ID, model.PENDING_MAILING_STATUS_CREATED)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) UpdateMailing(ctx context.Context, id string, req model.UpdateMailingRequest) (model.Mailing, error) {
	timeNow := time.Now()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return model.Mailing{}, err
	}

	mailing, err := r.updateMailing(ctx, id, req)
	if err != nil {
		_ = tx.Rollback()
		return model.Mailing{}, err
	}
	// если рассылку не нужно откладывать
	if timeNow.String() > req.CreatedAt.String() && timeNow.String() < req.FinishedAt.String() {
		err = tx.Commit()
		if err != nil {
			return model.Mailing{}, err
		}

		return mailing, err
	}

	err = r.createPendingMailing(ctx, mailing)
	if err != nil {
		_ = tx.Rollback()
		return model.Mailing{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Mailing{}, err
	}

	return mailing, nil
}

func (r repository) updateMailing(ctx context.Context, id string, req model.UpdateMailingRequest) (model.Mailing, error) {
	var t time.Time
	var mailing model.Mailing

	builder := squirrel.Update("mailing").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *")
	if req.CreatedAt != t {
		builder = builder.Set("created_at", req.CreatedAt)
	}
	if req.Text != "" {
		builder = builder.Set("text", req.Text)
	}
	if req.Filter != "" {
		builder = builder.Set("filter", req.Filter)
	}
	if req.FinishedAt != t {
		builder = builder.Set("finished_at", req.FinishedAt)
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.Mailing{}, err
	}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&mailing.ID, &mailing.CreatedAt, &mailing.Text, &mailing.Filter, &mailing.FinishedAt)
	if err != nil {
		return model.Mailing{}, err
	}

	return mailing, nil
}

func (r repository) UpdatePendingMailing(ctx context.Context, mailingID string, updateMailing model.UpdatePendingMailing) (model.PendingMailing, error) {
	var mailing model.PendingMailing

	builder := squirrel.Update("pending_mailing").
		Where(squirrel.Eq{"mailing_id": mailingID}).
		Suffix("RETURNING *")
	if updateMailing.Status != model.PENDING_MAILING_STATUS_UNKNOWN {
		builder = builder.Set("status", updateMailing.Status)
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.PendingMailing{}, err
	}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&mailing.ID, &mailing.MailingID, &mailing.Status)
	if err != nil {
		return model.PendingMailing{}, err
	}

	return mailing, nil
}

func (r repository) DeleteMailing(ctx context.Context, id string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = r.deleteMailing(ctx, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = r.DeletePendingMailing(ctx, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r repository) deleteMailing(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM mailing WHERE id=$1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r repository) DeletePendingMailing(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM pending_mailing WHERE mailing_id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) GetAllMailings(ctx context.Context) ([]model.Mailing, error) {
	var mailings []model.Mailing

	err := r.db.SelectContext(ctx, &mailings, "SELECT * FROM mailing")
	if err != nil {
		return nil, err
	}

	return mailings, nil
}

func (r repository) GetActiveMailings(ctx context.Context) ([]model.Mailing, error) {
	timeNow := time.Now()
	var mailings []model.Mailing

	err := r.db.SelectContext(ctx, &mailings, "SELECT * FROM mailing WHERE created_at <= $1 and $1 < finished_at", timeNow)
	if err != nil {
		return nil, err
	}

	return mailings, nil
}

func (r repository) GetMailingsToSending(ctx context.Context) ([]model.Mailing, error) {
	timeNow := time.Now()
	var mailings []model.Mailing

	err := r.db.SelectContext(ctx, &mailings, "SELECT mailing.id, created_at, text, filter, finished_at "+
		"FROM mailing JOIN pending_mailing ON mailing.id = pending_mailing.mailing_id "+
		"WHERE created_at <= $1 AND status = $2", timeNow, model.MESSAGE_STATUS_CREATED)
	if err != nil {
		return nil, err
	}

	return mailings, nil
}
