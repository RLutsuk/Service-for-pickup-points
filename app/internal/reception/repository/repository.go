package repository

import (
	"database/sql"

	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
)

type RepositoryI interface {
	CreateReception(reception *models.Reception) error
	GetOpenReceptionByPPID(pickupPointID string) (string, string, error)
	CloseReception(id string) error
}

type dataBase struct {
	db *sql.DB
}

func New(db *sql.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbReception *dataBase) CreateReception(reception *models.Reception) error {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbReception.db)
	query := queryBuilder.Insert("receptions").Columns("pickup_point_id", "status_reception").
		Values(reception.PickupPointID, reception.Status).Suffix("RETURNING id, date_time")

	err := query.QueryRow().Scan(&reception.ID, &reception.DateTime)
	if err != nil {
		return errors.Wrap(err, "database error (table receptions, CreateReception)")
	}
	return nil
}

func (dbReception *dataBase) GetOpenReceptionByPPID(pickupPointID string) (string, string, error) {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbReception.db)
	query := queryBuilder.Select("id", "date_time").From("receptions").Where(sq.And{sq.Eq{"pickup_point_id": pickupPointID}, sq.Eq{"status_reception": "in_progress"}})
	id, dateTime := "", ""
	err := query.QueryRow().Scan(&id, &dateTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil
		}
		return "", "", errors.Wrap(err, "database error (table receptions, GetOpenReceptionByPPID)")
	}
	return id, dateTime, nil
}

func (dbReception *dataBase) CloseReception(id string) error {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbReception.db)
	query := queryBuilder.Update("receptions").Set("status_reception", "close").Where(sq.Eq{"id": id})

	_, err := query.Exec()
	if err != nil {
		return errors.Wrap(err, "database error (table receptions, CloseReception)")
	}

	return nil
}
