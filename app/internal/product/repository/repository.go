package repository

import (
	"database/sql"

	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
)

type RepositoryI interface {
	CreateProduct(product *models.Product) error
	DeleteLastProduct(productID string) error
	GetLastProdcutByDate(receptionID string) (string, error)
}

type dataBase struct {
	db *sql.DB
}

func New(db *sql.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbReception *dataBase) CreateProduct(product *models.Product) error {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbReception.db)
	query := queryBuilder.Insert("products").Columns("reception_id", "type_product").
		Values(product.ReceptionID, product.TypeProduct).Suffix("RETURNING id, date_time")

	err := query.QueryRow().Scan(&product.ID, &product.DateTime)
	if err != nil {
		return errors.Wrap(err, "database error (table products, CreateProduct)")
	}
	return nil
}

func (dbReception *dataBase) DeleteLastProduct(productID string) error {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbReception.db)
	query := queryBuilder.Delete("products").Where(sq.Eq{"id": productID})

	_, err := query.Exec()
	if err != nil {
		return errors.Wrap(err, "database error (table products, DeleteLastProduct)")
	}

	return nil
}

func (dbReception *dataBase) GetLastProdcutByDate(receptionID string) (string, error) {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbReception.db)
	query := queryBuilder.Select("id").From("products").Where(sq.Eq{"reception_id": receptionID}).OrderBy("date_time DESC").Limit(1)

	lastProductID := ""
	err := query.QueryRow().Scan(&lastProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrEmptyReception
		}
		return "", errors.Wrap(err, "database error (table products, GetLastProdcutByDate)")
	}
	return lastProductID, nil
}
