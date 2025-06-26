package sqlutils

import "database/sql"

func ToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func ToInt64(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

func ToInt32(ni sql.NullInt32) int32 {
	if ni.Valid {
		return ni.Int32
	}
	return 0
}

func ToUnix(nt sql.NullTime) int64 {
	if nt.Valid {
		return nt.Time.Unix()
	}
	return 0
}
