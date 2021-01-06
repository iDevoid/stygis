package query

const (
	// UserInsert records the new user to db table users
	UserInsert = `
		INSERT INTO users (
			username,
			email,
			hashed_email,
			password,
			create_time,
			status
		) VALUES (
			:username,
			:email,
			:hashed_email,
			:password,
			:create_time,
			:status				
		) RETURNING id;
	`
)
