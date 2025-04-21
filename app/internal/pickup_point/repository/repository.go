package repository

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
)

type RepositoryI interface {
	CreatePickupPoint(pickupPoint *models.PickupPoint) error
	GetAllPickupPoint(startDate, endDate string, offset, limit int) ([]*models.PickupPoint, error)
	GetPickupPointByID(pickupPointID string) error
	GetListOnlyPickupPoint() ([]*models.PickupPoint, error)
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
	pickupPointQuery := queryBuilder.
		Select("id", "registration_date", "city").
		From("pickup_points").
		OrderBy("registration_date ASC").
		Offset(uint64(offset)).
		Limit(uint64(limit))

	if startDate != "" {
		pickupPointQuery = pickupPointQuery.Where(sq.GtOrEq{"registration_date": startDate})
	}
	if endDate != "" {
		pickupPointQuery = pickupPointQuery.Where(sq.LtOrEq{"registration_date": endDate})
	}

	pickupRows, err := pickupPointQuery.Query()
	if err != nil {
		return nil, errors.Wrap(err, "database error (table pcikup_points, GetAllPickupPoint, select pp)")
	}
	defer pickupRows.Close()

	pickupPoints := make([]*models.PickupPoint, 0)
	pickupPointIDs := make([]string, 0)
	pickupPointMap := make(map[string]*models.PickupPoint)

	for pickupRows.Next() {
		var pp models.PickupPoint
		if err := pickupRows.Scan(&pp.ID, &pp.RegistrationDate, &pp.City); err != nil {
			return nil, err
		}
		pickupPoints = append(pickupPoints, &pp)
		pickupPointMap[pp.ID] = &pp
		pickupPointIDs = append(pickupPointIDs, pp.ID)
	}

	receptionProductQuery := queryBuilder.
		Select(
			"r.id", "r.date_time", "r.status_reception", "r.pickup_point_id",
			"p.id", "p.date_time", "p.type_product", "p.reception_id").
		From("receptions r").
		LeftJoin("products p ON r.id = p.reception_id").
		Where(sq.Eq{"r.pickup_point_id": pickupPointIDs})

	rows, err := receptionProductQuery.Query()
	if err != nil {
		return nil, errors.Wrap(err, "database error (table pcikup_points, GetAllPickupPoint, select receptions/products)")
	}
	defer rows.Close()

	receptionMap := make(map[string]*models.Reception)

	for rows.Next() {
		var (
			recID, status, pvzID          string
			recTime                       time.Time
			prodID, prodType, receptionID sql.NullString
			prodTime                      sql.NullTime
		)

		if err := rows.Scan(&recID, &recTime, &status, &pvzID,
			&prodID, &prodTime, &prodType, &receptionID); err != nil {
			return nil, err
		}

		reception, exists := receptionMap[recID]
		if !exists {
			reception = &models.Reception{
				ID:            recID,
				DateTime:      recTime,
				Status:        status,
				PickupPointID: pvzID,
				Products:      make([]*models.Product, 0),
			}
			receptionMap[recID] = reception
			if pp, ok := pickupPointMap[pvzID]; ok {
				pp.Receptions = append(pp.Receptions, reception)
			}
		}

		if prodID.Valid {
			product := &models.Product{
				ID:          prodID.String,
				TypeProduct: prodType.String,
				ReceptionID: receptionID.String,
				DateTime:    prodTime.Time,
			}
			reception.Products = append(reception.Products, product)
		}
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

func (dbPickupPoint *dataBase) GetListOnlyPickupPoint() ([]*models.PickupPoint, error) {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbPickupPoint.db)
	pickupPointQuery := queryBuilder.
		Select("id", "registration_date", "city").
		From("pickup_points")

	pickupRows, err := pickupPointQuery.Query()
	if err != nil {
		return nil, errors.Wrap(err, "database error (table pcikup_points, GetAllPickupPoint, select pp)")
	}
	defer pickupRows.Close()

	pickupPoints := make([]*models.PickupPoint, 0)

	for pickupRows.Next() {
		var pp models.PickupPoint
		if err := pickupRows.Scan(&pp.ID, &pp.RegistrationDate, &pp.City); err != nil {
			return nil, err
		}
		pickupPoints = append(pickupPoints, &pp)
	}

	return pickupPoints, nil
}
