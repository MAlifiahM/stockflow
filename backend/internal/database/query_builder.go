package database

import "gorm.io/gorm"

func BuildBaseQuery(db *gorm.DB, model interface{}, preloads []string) *gorm.DB {
	query := db.Model(model)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	return query
}

// BuildBaseQueryWithConditions creates a base query with conditions and preloads
func BuildBaseQueryWithConditions(db *gorm.DB, model interface{}, conditions map[string]interface{}, preloads []string) *gorm.DB {
	query := BuildBaseQuery(db, model, preloads)

	for field, value := range conditions {
		query = query.Where(field, value)
	}

	return query
}

// BuildQueryWithJoins creates a query with joins and preloads
func BuildQueryWithJoins(db *gorm.DB, model interface{}, joins []string, preloads []string) *gorm.DB {
	query := db.Model(model)

	for _, join := range joins {
		query = query.Joins(join)
	}

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	return query
}
