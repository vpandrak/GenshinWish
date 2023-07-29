package Repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"rest-todo/internal/model"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, user model.User) (*model.User, error)
	All(ctx context.Context) ([]model.User, error)
	GetByName(ctx context.Context, name string) (*model.User, error)
	Update(ctx context.Context, id int, updated model.User) (*model.User, error)
	Delete(ctx context.Context, id int) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Migrate(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		money INT default 0
	    );
	`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	var id int
	err := r.db.QueryRowContext(ctx, "INSERT INTO users (name, password) values ($1, $2) returning id;", user.Name, user.Password).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	user.Id = id
	return &user, nil
}

func (r *PostgresRepository) All(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []model.User

	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.Money); err != nil {
			return nil, err
		}
		all = append(all, user)
	}
	return all, nil
}

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, "select * from users WHERE name = $1", name)
	var user model.User
	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.Money); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id int, updated model.User) (*model.User, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE users SET name = $1, password = $2, money= $3 WHERE id=$4", updated.Name, updated.Password, updated.Money, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}
	return &updated, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
