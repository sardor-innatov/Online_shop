package migrate

func cardTable(query string) string {
	query += `
	CREATE TABLE IF NOT EXISTS cards (
    	id BIGSERIAL PRIMARY KEY,
    	user_id BIGINT REFERENCES users(id),
    	stripe_method_id VARCHAR(255) NOT NULL, 
    	stripe_customer_id VARCHAR(255) NOT NULL, 
    	last4 VARCHAR(4) NOT NULL,
    	payment_system VARCHAR(25),
    	name VARCHAR(25)
	);`

	return query
}
