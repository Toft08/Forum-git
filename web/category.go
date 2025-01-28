package web

import ( "database/sql"
"log" )

// AddCategory adds a new category to the database
func AddCategory(db *sql.DB, name string) error {
	_, err := db.Exec("INSERT INTO Category (name) VALUES (?)", name)
	if err != nil {
		log.Println("Error adding category:", err)
		return err
	}
	return nil
}

// GetCategories retrieves all categories from the database
func GetCategories(db *sql.DB) ([]Category, error) {
    rows, err := db.Query("SELECT id, name FROM Category")
    if err != nil {
        log.Println("Error retrieving categories:", err)
        return nil, err
    }
    defer rows.Close()

    var categories []Category
    for rows.Next() {
        var category Category
        if err := rows.Scan(&category.ID, &category.Name); err != nil {
            log.Println("Error scanning category:", err)
            return nil, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}
