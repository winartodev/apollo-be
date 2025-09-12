package repository

const (
	getUserByID = `
		SELECT 
		    usr.id,
		    usr.username,
		    usr.email,
		    usr.first_name,
		    usr.last_name,
		    usr.phone_number
		FROM users AS usr
		WHERE usr.id = $1
	`

	getUserByEmail = `
			SELECT 
		    usr.id,
		    usr.username,
		    usr.email,
		    usr.first_name,
		    usr.last_name,
		    usr.phone_number
		FROM users AS usr
		WHERE usr.email = $1
	`

	getUserByUsername = `
			SELECT 
		    usr.id,
		    usr.username,
		    usr.email,
		    usr.first_name,
		    usr.last_name,
		    usr.phone_number
		FROM users AS usr
		WHERE usr.username = $1
	`

	getUserQuery = `
		SELECT 
		    usr.username,
		    usr.email,
		    usr.first_name,
		    usr.last_name,
		    usr.phone_number
		FROM users AS usr
	`
)
