package migrate

func cardTable(query string) string {
	query += `
	CREATE TABLE IF NOT EXISTS cards (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT REFERENCES users(id),
		number VARCHAR(16) NOT NULL,
		name VARCHAR(25),
		payment_system VARCHAR(25) NOT NULL,
		balance NUMERIC(12,2) DEFAULT 0 CHECK (balance >= 0)
	);`

	return query
}
