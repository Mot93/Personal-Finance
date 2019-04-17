package financedatabase

import (
	"fmt"
	"strings"
)

// createtablecategory creates the table category if it doensn't already exist
func createTableCategory() {
	categoryTable := `
	CREATE TABLE IF NOT EXISTS categories(
		category TEXT PRIMARY KEY
	);
	`
	executeCommand(categoryTable, "create table categories")
}

// Category in db
type Category string

// Equals checks that two categories are the same
func (c Category) Equals(c2 Category) bool {
	if strings.Compare(c.String(), c2.String()) == 0 {
		return true
	}
	return false
}

// Delete deletes a category from categories
func (c Category) delete() {
	sqlCategory := fmt.Sprintf(`
	DELETE FROM categories
	WHERE category = '%v'
	`, c.String())
	executeCommand(sqlCategory, "delte category")
}

// String returns the name of the category
func (c Category) String() string {
	return string(c)
}

// Categories is a multitude og Category
type Categories []Category

// GetAll return the complete list of categories
func (ca *Categories) GetAll() {
	// Emptying the slice
	*ca = nil
	sqlCategories := `
	SELECT * FROM categories
	ORDER BY category
	`
	// return the result of a select
	row, err := db.Query(sqlCategories)
	if err != nil {
		panic(fmt.Errorf("Error reading all categories: %v", err))
	}
	defer row.Close()
	// reading all element in the row
	for row.Next() {
		var c string
		err2 := row.Scan(&c)
		if err2 != nil {
			panic(fmt.Errorf("Error reading all categories: %v", err2))
		}
		*ca = append(*ca, Category(c))
	}
}

// Add add a category to the table categories
func (ca *Categories) Add(c string) {
	sqlCategory := fmt.Sprintf(`
	INSERT INTO categories(
		category
	) VALUES ('%v');
	`, c)
	storeItem(sqlCategory, "category")
	(*ca).GetAll()
}

// Delete deletes a category
func (ca *Categories) Delete(c Category) {
	c.delete()
	(*ca).GetAll()
}
