package db

import (
	"database/sql"
	"time"

	"github.com/ajdinahmetovic/item-service/logger"
	"golang.org/x/crypto/bcrypt"
)

//User struct
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Items    []Item `json:"items"`
}

//UserCredentials struct
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Login func
func Login(userCredentials *UserCredentials) (*int, error) {
	var id int
	var password string
	sqlState := `
	SELECT id, password FROM app_user 
	WHERE username = $1`
	err := db.QueryRow(sqlState, userCredentials.Username).Scan(&id, &password)
	if err != nil {
		logger.Error("Failed to save user to database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(userCredentials.Password))
	if err != nil {
		return nil, err
	}

	return &id, nil
}

//AddUser func
func AddUser(user *User) (*int, error) {
	logger.Info("Create user request", time.Now(), "user", user)
	var id int
	sqlState := `INSERT INTO app_user (username, full_name, password) VALUES ($1, $2, $3) RETURNING id;`

	err := db.QueryRow(sqlState, user.Username, user.FullName, user.Password).Scan(&id)
	if err != nil {
		logger.Error("Failed to add user to database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return nil, err
	}

	return &id, nil
}

//FindUser func
func FindUser(user *User) ([]*User, error) {
	var response []*User
	var rows *sql.Rows
	var itemRows *sql.Rows
	var err error

	sqlItem := `
		SELECT * 
		FROM item
		WHERE
		user_id = $1;
	`
	sqlState := `
	SELECT id, username, full_name 
	FROM app_user 
	WHERE
	username LIKE $1 and
	full_name LIKE $2`

	if user.ID != 0 {
		sqlState += ` and id = $3`
		rows, err = db.Query(sqlState, user.Username+"%", user.FullName+"%", user.ID)
	} else {
		rows, err = db.Query(sqlState, user.Username+"%", user.FullName+"%")
	}

	if err != nil {
		logger.Error("Failed to get users from database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return nil, err
	}

	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.FullName,
		)
		if err != nil {
			logger.Error("Failed to scan users from row returned by database", "time", time.Now(), "err", err)
			return nil, err
		}
		itemRows, err = db.Query(sqlItem, user.ID)
		var items []Item
		for itemRows.Next() {
			item := Item{}
			err = itemRows.Scan(
				&item.ID,
				&item.Title,
				&item.Description,
				&item.UserID,
			)
			if err != nil {
				logger.Error("Failed to scan item rows in user", "time", time.Now(), "err", err)
				return nil, err
			}
			items = append(items, item)
		}
		user.Items = items
		response = append(response, &user)
	}
	return response, nil
}

//DeleteUser func
func DeleteUser(id int) error {
	sqlState := `DELETE FROM item
	WHERE user_id = $1;`
	_, err := db.Exec(sqlState, id)
	if err != nil {
		logger.Error("Failed to delete user in database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return err
	}
	sqlState = `
	DELETE FROM app_user
	WHERE id = $1;`
	_, err = db.Exec(sqlState, id)
	if err != nil {
		logger.Error("Failed to delete users item in database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return err
	}
	return nil
}

//UpdateUser func
func UpdateUser(user *User) error {
	sqlState := `
	UPDATE app_user
	SET username = $1,
	full_name = $2
	WHERE id = $3;`
	_, err := db.Exec(sqlState, user.Username, user.FullName, user.ID)
	if err != nil {
		logger.Error("Failed to update user in database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return err
	}
	return nil
}
