package migrate 

func orderTable (query string) string{
	query += ` 
	CREATE TABLE IF NOT EXISTS "orders" (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT REFERENCES users(id),
		total_price	NUMERIC(12, 2) NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP 
	);`

	return query
}

func orderItemTable (query string) string{
	query += `
	CREATE TABLE IF NOT EXISTS "order_items" (
		order_id BIGINT REFERENCES orders(id),
		product_id BIGINT REFERENCES products(id),
		quantity INTEGER NOT NULL,
		price NUMERIC(12, 2) NOT NULL
	);`

	return query
}