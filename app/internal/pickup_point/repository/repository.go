package repository

import (
	"database/sql"

	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
)

type RepositoryI interface {
	CreatePickupPoint(pickupPoint *models.PickupPoint) error
	GetAllPickupPoint(startDate, endDate string, offset, limit int) ([]*models.PickupPoint, error)
	GetPickupPointByID(pickupPointID string) error
}

type dataBase struct {
	db *sql.DB
}

func New(db *sql.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbPickupPoint *dataBase) CreatePickupPoint(pickupPoint *models.PickupPoint) error {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbPickupPoint.db)
	query := queryBuilder.Insert("pickup_points").Columns("city").
		Values(pickupPoint.City).Suffix("RETURNING id, registration_date")

	err := query.QueryRow().Scan(&pickupPoint.ID, &pickupPoint.RegistrationDate)
	if err != nil {
		return errors.Wrap(err, "database error (table pcikup_points, CreatePickupPoint)")
	}

	return nil
}

func (dbPickupPoint *dataBase) GetAllPickupPoint(startDate, endDate string, offset, limit int) ([]*models.PickupPoint, error) {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbPickupPoint.db)
	query := queryBuilder.Select("*").From("pickup_points").OrderBy("registration_date ASC").Offset(uint64(offset)).Limit(uint64(limit))

	if startDate != "" {
		query = query.Where(sq.GtOrEq{"registration_date": startDate})
	}
	if endDate != "" {
		query = query.Where(sq.LtOrEq{"registration_date": endDate})
	}

	rows, err := query.Query()
	if err != nil {
		return nil, errors.Wrap(err, "database error (table pickup_points, GetAllPickupPoint)")
	}
	defer rows.Close()

	pickupPoints := make([]*models.PickupPoint, 0)

	for rows.Next() {
		var pp models.PickupPoint
		err := rows.Scan(&pp.ID, &pp.RegistrationDate, &pp.City)
		if err != nil {
			return nil, errors.Wrap(err, "database error (table pickup_points, GetAllPickupPoint)")
		}
		pickupPoints = append(pickupPoints, &pp)
	}

	return pickupPoints, nil
}

func (dbPickupPoint *dataBase) GetPickupPointByID(pickupPointID string) error {
	row := dbPickupPoint.db.QueryRow("SELECT EXISTS(SELECT 1 FROM pickup_points WHERE id = $1)", pickupPointID)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return errors.Wrap(err, "database error (table pickup_points, GetPickupPointByID)")
	}

	if !exists {
		return models.ErrPickupPointDontExist
	}

	return nil
}
