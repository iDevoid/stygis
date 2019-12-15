package query

const (
	// InsertNewUser is the query to make a new row inside table accont for creating new user
	InsertNewUser = `
		INSERT INTO account (
			email, 
			hash_password, 
			username, 
			created_at, 
			last_login, 
			status
		) VALUES (
			:email,
			:password,
			:username,
			NOW(),
			NOW(),
			:status
		)
		RETURNING id
	`

	// SelectUserByID query is selecting the user information only by id
	// without selecting the important data
	SelectUserByID = `
		SELECT 
			id, 
			username,  
			last_login, 
			status 
		FROM 
			account
		WHERE
			id = :user_id AND
			statuts = :status
	`

	// SelectUserByEmail is the query to get all user account data
	// with email and hashed password for making sure that user is verified to get the data
	// usually use for login
	SelectUserByEmail = `
		SELECT 
			id, 
			email, 
			hash_password, 
			username, 
			created_at, 
			last_login, 
			status 
		FROM 
			account
		WHERE
			email = :email AND
			hash_password = :password AND
			statuts = :status
	`

	// UpdateUserPassword uses old hashed password to replace it with newest hased password
	UpdateUserPassword = `
		UPDATE 
			account 
		SET  
			hash_password = :new_password
		WHERE 
			id = :user_id AND
			email = :email AND
			hash_password = :old_password AND
			statuts = :status
	`

	// DeactivateUser makes a change on accont status to change it to inactive accont
	DeactivateUser = `
		UPDATE 
			account 
		SET  
			status = :inactive_status
		WHERE 
			id = :user_id AND
			email = :email AND
			hash_password = :password AND
			statuts = :active_status
	`
)
