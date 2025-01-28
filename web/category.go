package web

import (
	"log"
)

// AddCategory adds a new category to the database this is if we want to add a new category to the database
// func AddCategory(db *sql.DB, name string) error {
// 	_, err := db.Exec("INSERT INTO Category (name) VALUES (?)", name)
// 	if err != nil {
// 		log.Println("Error adding category:", err)
// 		return err
// 	}
// 	return nil
// }

// GetCategories retrieves all categories from the database
func GetCategories() ([]CategoryDetails, error) {
	rows, err := db.Query("SELECT id, name FROM Category")
	if err != nil {
		log.Println("Error retrieving categories:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryDetails
	for rows.Next() {
		var category CategoryDetails
		if err := rows.Scan(&category.CategoryID, &category.CategoryName); err != nil {
			log.Println("Error scanning category:", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
