package user

import (
	"database/sql"
	"fmt"

	"github.com/panjiasmoroart/gopher-ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	// return &Store{db: db}
	if db == nil {
		panic("database connection is nil")
	}
	return &Store{db: db}
}

// func (s *Store) CreateUser(user types.User) error {
// 	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *Store) CreateUser(user types.User) (*types.User, error) {

	// if s.db == nil {
	// 	return nil, fmt.Errorf("database connection is nil")
	// }

	query := "INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)"
	result, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return nil, fmt.Errorf("could not insert user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("could not get last insert id: %v", err)
	}

	// Set the ID of the newly created user
	user.ID = int(id)
	return &user, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT id, firstName, lastName, email, password, createdAt FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT id, firstName, lastName, email, password, createdAt FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
