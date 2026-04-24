package migrate

func userTable(query string) string {

	query += `
	CREATE TABLE IF NOT EXISTS "users" (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(25) NOT NULL,
		last_name VARCHAR(25) NOT NULL,
		email VARCHAR(245) NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	return query
}

func roleTable(query string) string {
	query += `
	CREATE TABLE IF NOT EXISTS roles (
	    id 	BIGSERIAL PRIMARY KEY,
	    role_name VARCHAR(25) NOT NULL,
	    permissions JSONB NOT NULL
    );`

	return query
}

func permissionTable(query string) string {
	query += `
	CREATE TABLE IF NOT EXISTS permissions (
	id 	BIGSERIAL PRIMARY KEY,
	path  TEXT NOT NULL,
	method  TEXT NOT NULL,
	endpoint_name VARCHAR(25)
);`

	return query
}

func userRoleTable(query string) string {
	query += `
	CREATE TABLE IF NOT EXISTS user_roles (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL,
	role_id BIGINT NOT NULL
);`

	return query
}
