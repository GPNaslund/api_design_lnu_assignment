package dataaccess

import (
	webhookdto "1dv027/aad/internal/dto/user/webhook"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersWebhooksDataAccess struct {
	dbPool *pgxpool.Pool
}

func NewUsersWebhookDataAccess(dbPool *pgxpool.Pool) UsersWebhooksDataAccess {
	return UsersWebhooksDataAccess{
		dbPool: dbPool,
	}
}

func (u UsersWebhooksDataAccess) DeleteUserWebhook(ctx context.Context, userId int) error {
	query := `DELETE FROM UserWebhooks WHERE user_id = $1`
	result, err := u.dbPool.Exec(ctx, query, userId)
	if err != nil {
		return &customerrors.DatabaseError{Message: "something went wrong trying to delete user webhook"}
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return &customerrors.WebhookNotFoundError{Message: "no webhook was deleted"}
	}
	return nil
}

func (u UsersWebhooksDataAccess) GetUserWebhook(ctx context.Context, userId int) (model.Webhook, error) {
	emptyModel := model.Webhook{}
	query := `SELECT * FROM UserWebhooks WHERE user_id = $1`
	var webhookModel model.Webhook
	err := u.dbPool.QueryRow(ctx, query, userId).Scan(
		&webhookModel.Id,
		&webhookModel.EndpointUrl,
		&webhookModel.ClientSecret,
		&webhookModel.Actions,
		&webhookModel.UserId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyModel, &customerrors.WebhookNotFoundError{}
		}
		return emptyModel, err
	}
	return webhookModel, nil
}

func (u UsersWebhooksDataAccess) CreateNewWebhook(ctx context.Context, userId int, data webhookdto.NewUserWebhookDTO) error {
	query := `INSERT INTO UserWebhooks (webhook_endpoint, client_secret, webhook_actions, user_id) VALUES ($1, $2, $3, $4)`
	result, err := u.dbPool.Exec(ctx, query, data.EndpointUrl, data.ClientSecret, data.Actions, userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return &customerrors.UserNotFoundError{Message: "No registered user for the webhook"}
			}
		}

		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return &customerrors.InvalidWebhookDataError{Message: "could not create a new webhook entry"}
	}
	return nil
}

func (u UsersWebhooksDataAccess) UpdateUserWebhook(ctx context.Context, userId int, data webhookdto.UpdateUserWebhookDTO) error {
	query := u.createUpdateWebhookQuery(userId, data)
	result, err := u.dbPool.Exec(ctx, query)
	if err != nil {
		return &customerrors.DatabaseError{Message: "could not update userwebhook"}
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return &customerrors.WebhookNotFoundError{Message: "could not find the user webhook for updating"}
	}
	return nil
}

func (u UsersWebhooksDataAccess) GetAllWebhooksByAction(ctx context.Context, action model.WebhookAction) ([]model.Webhook, error) {
	query := `SELECT id, webhook_endpoint, client_secret, webhook_actions, user_id FROM UserWebhooks WHERE $1 = ANY(webhook_actions)`
	webhooks := []model.Webhook{}

	rows, err := u.dbPool.Query(ctx, query, action)
	if err != nil {
		return nil, &customerrors.DatabaseError{}
	}
	defer rows.Close()

	for rows.Next() {
		var webhook model.Webhook
		err := rows.Scan(&webhook.Id, &webhook.EndpointUrl, &webhook.ClientSecret, &webhook.Actions, &webhook.UserId)
		if err != nil {
			return nil, &customerrors.DatabaseError{}
		}
		webhooks = append(webhooks, webhook)
	}

	if err = rows.Err(); err != nil {
		return nil, &customerrors.DatabaseError{}
	}

	return webhooks, nil
}

func (u UsersWebhooksDataAccess) createUpdateWebhookQuery(userId int, webhookData webhookdto.UpdateUserWebhookDTO) string {
	query := `UPDATE UserWebhooks SET `
	updates := []string{}

	addUpdate := func(field string, value interface{}) {
		switch v := value.(type) {
		case string:
			updates = append(updates, fmt.Sprintf("%s = '%s'", field, strings.ReplaceAll(v, "'", "''")))
		case []string:
			arrayStr := "{" + strings.Join(v, ",") + "}"
			updates = append(updates, fmt.Sprintf("%s = '%s'", field, arrayStr))
		}
	}

	if webhookData.EndpointUrl != nil {
		addUpdate("webhook_endpoint", *webhookData.EndpointUrl)
	}
	if webhookData.Actions != nil {
		addUpdate("webhook_actions", *webhookData.Actions)
	}
	if webhookData.ClientSecret != nil {
		addUpdate("client_secret", *webhookData.ClientSecret)
	}

	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE user_id = %d;", userId)

	return query
}
