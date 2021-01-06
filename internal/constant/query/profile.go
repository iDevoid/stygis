package query

const (
	// ProfileInsert is the query to record new data into database table profiles
	ProfileInsert = `
		INSERT INTO profiles (
			id,
			username,
			full_name,
			status,
			create_time
		) VALUES (
			:id,
			:username,
			:full_name,
			:status,
			:create_time
		)
	`
)
