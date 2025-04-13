package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/pkg/errors"
)

type RepositoryI interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(user *models.User) (*models.User, error)
	ChekUserByEmail(email string) error
	GetUserByEmail(user *models.User) error
}

type dataBase struct {
	db *sql.DB
}

func New(db *sql.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbUser *dataBase) CreateUser(user *models.User) (*models.User, error) {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbUser.db)
	query := queryBuilder.Insert("employees").Columns("email", "password_user", "role_user").
		Values(user.Email, user.Password, user.Role).Suffix("RETURNING id")

	err := query.QueryRow().Scan(&user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "database error (table employees, CreateUser)")
	}

	return user, nil
}

func (dbUser *dataBase) GetUserByID(user *models.User) (*models.User, error) {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbUser.db)

	query := queryBuilder.Select("email", "role_user").From("employees").Where(sq.Eq{"id": user.ID})

	err := query.QueryRow().Scan(&user.Email, &user.Role)
	if err != nil {
		return nil, errors.Wrap(err, "database error (table employees, GetUserByID)")
	}

	return user, nil
}

func (dbUser *dataBase) ChekUserByEmail(email string) error {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbUser.db)
	query := queryBuilder.Select("COUNT(*)").From("employees").Where(sq.Eq{"email": email})

	var count int
	err := query.QueryRow().Scan(&count)

	if err != nil {
		return errors.Wrap(err, "database error (table employees, ChekUserByEmail)")
	}

	if count == 1 {
		return models.ErrUserExist
	}

	return nil
}

func (dbUser *dataBase) GetUserByEmail(user *models.User) error {
	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbUser.db)

	query := queryBuilder.Select("role_user", "password_user").From("employees").Where(sq.Eq{"email": user.Email})
	err := query.QueryRow().Scan(&user.Role, &user.Password)
	if err != nil {
		return errors.Wrap(err, "database error (table employees, GetUserByEmail)")
	}

	return nil
}
