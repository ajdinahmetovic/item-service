package db

import (
	"database/sql"
	"time"

	"github.com/ajdinahmetovic/item-service/logger"
)

// Item struct
type Item struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"userid"`
}

//AddItem function
func AddItem(item *Item) error {
	sqlState := `
	INSERT INTO item (title, description, user_id)
	VALUES ($1, $2, $3);`
	_, err := db.Exec(sqlState, &item.Title, &item.Description, &item.UserID)
	if err != nil {
		logger.Error("Failed to save item to database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return err
	}
	return nil
}

//FindItem func
func FindItem(item *Item) ([]*Item, error) {
	var data []*Item
	var rows *sql.Rows
	var err error

	sqlState := `
		SELECT * 
		FROM item 
		WHERE
		title LIKE $1 and
		description LIKE $2`

	if item.ID != 0 {
		sqlState += ` and id = $3`
		rows, err = db.Query(sqlState, item.Title+"%", item.Description+"%", item.ID)
	} else {
		rows, err = db.Query(sqlState, item.Title+"%", item.Description+"%")
	}
	if err != nil {
		logger.Error("Failed to get items from database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return nil, err
	}

	for rows.Next() {
		item := Item{}
		err = rows.Scan(
			&item.Title,
			&item.Description,
			&item.ID,
			&item.UserID,
		)
		if err != nil {
			logger.Error("Failed to scan rows returned to database", "time", time.Now(), "err", err)
			return nil, err
		}
		data = append(data, &item)

	}
	return data, nil
}

//UpdateItem func
func UpdateItem(item *Item) error {
	sqlState := `
	UPDATE item 
	SET title = $1,
	description = $2
	WHERE id = $3`
	_, err := db.Exec(sqlState, item.Title, item.Description, item.ID)
	if err != nil {
		logger.Error("Failed to update item to database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return err
	}
	return nil

}

//DeleteItem func
func DeleteItem(id int) error {
	sqlState := `
	DELETE FROM item 
	WHERE id = $1`
	_, err := db.Exec(sqlState, id)
	if err != nil {
		logger.Error("Failed to delete item in database", "time", time.Now(), "err", err, "SQL state", sqlState)
		return err
	}
	return nil
}
