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

	getUserQuery = `
		SELECT 
		    usr.id,
		    usr.username,
		    usr.email,
		    usr.first_name,
		    usr.last_name,
		    usr.phone_number
		FROM users AS usr
	`
)
