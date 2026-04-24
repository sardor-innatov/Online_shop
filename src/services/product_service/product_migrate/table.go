package migrate

func productTable(query string) string {

	query += `
	CREATE TABLE IF NOT EXISTS "products" (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		description VARCHAR(250),
		price  NUMERIC(12, 2) NOT NULL,
		stock INTEGER NOT NULL CHECK (stock >= 0),
		category_id BIGSERIAL REFERENCES categories(id),
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	return query
}

func categoryTable(query string) string {

	query += `
	CREATE TABLE IF NOT EXISTS "categories" (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(25) NOT NULL
	);`

	return query
}
