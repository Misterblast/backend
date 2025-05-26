package app

func NullIfZero(id int32) interface{} {
	if id == 0 {
		return nil // PostgreSQL will auto-generate the ID if it's SERIAL
	}
	return id
}
