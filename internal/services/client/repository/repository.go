package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"gitlab.com/fr5270937/notifications_service/internal/model"
	desc "gitlab.com/fr5270937/notifications_service/internal/services/client"
)

type repository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) desc.Repository {
	return repository{db: db}
}

func (r repository) CreateClient(ctx context.Context, req model.CreateClientRequest) (model.Client, error) {
	var client model.Client

	err := r.db.QueryRowContext(ctx, `INSERT INTO client (phone_number, mobile_operator_code, tag) VALUES ($1, $2, $3) RETURNING *`,
		req.PhoneNumber, req.MobileOperatorCode, req.Tag).
		Scan(&client.ID, &client.PhoneNumber, &client.MobileOperatorCode, &client.Tag, &client.TimeZone)
	if err != nil {
		return model.Client{}, err
	}

	return client, nil
}

func (r repository) UpdateClient(ctx context.Context, id string, req model.UpdateClientRequest) (model.Client, error) {
	var client model.Client

	builder := squirrel.Update("client").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *")
	if req.PhoneNumber != "" {
		builder = builder.Set("phone_number", req.PhoneNumber)
	}
	if req.MobileOperatorCode != 0 {
		builder = builder.Set("mobile_operator_code", req.MobileOperatorCode)
	}
	if req.Tag != "" {
		builder = builder.Set("tag", req.Tag)
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.Client{}, err
	}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&client.ID, &client.PhoneNumber, &client.MobileOperatorCode, &client.Tag, &client.TimeZone)
	if err != nil {
		return model.Client{}, err
	}

	return client, nil
}

func (r repository) DeleteClient(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM client WHERE id=$1", id)
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

func (r repository) GetClientsFromMobileOperatorCode(ctx context.Context, mobileOperatorCode int32) ([]model.Client, error) {
	var clients []model.Client

	err := r.db.SelectContext(ctx, &clients, "SELECT * FROM client WHERE mobile_operator_code=$1", mobileOperatorCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return clients, nil
		}

		return nil, err
	}

	return clients, nil
}

func (r repository) GetClientsFromTag(ctx context.Context, tag string) ([]model.Client, error) {
	var clients []model.Client

	err := r.db.SelectContext(ctx, &clients, "SELECT * FROM client WHERE tag=$1", tag)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (r repository) GetClientFromID(ctx context.Context, id string) (model.Client, error) {
	var client model.Client

	err := r.db.GetContext(ctx, &client, "SELECT * FROM client WHERE id=$1", id)
	if err != nil {
		return client, err
	}

	return client, nil
}
