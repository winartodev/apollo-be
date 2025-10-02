package repository

const (
	registerUserQuery = `
		INSERT INTO users (username, email, phone_number, first_name, last_name, password, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`

	updateRefreshTokenQuery = `
		UPDATE users 
			SET 
			    refresh_token = $2,
				updated_at = $3
		WHERE id = $1
	`

	getUserData = `
		SELECT 
		    usr.id,
		    usr.username,
		    usr.email,
		    usr.password
		FROM users AS usr
		WHERE usr.username = $1 OR usr.email = $2
	`

	updatePasswordQueryDB = `
		UPDATE users SET 
			password = $2, 
			updated_at = $3 
		WHERE id = $1
	`
)
