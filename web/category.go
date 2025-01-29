package web

import (
	"log"
)

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
