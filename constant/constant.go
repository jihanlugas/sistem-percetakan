package constant

// TokenPayloadLen
// 4 int migration
// 8 int8 migration
// 24 string
const (
	TokenUserContext     = "usr"
	TokenContentLen      = 6
	FormatDateLayout     = "02-01-2006"
	FormatDatetimeLayout = "02-01-2006 15:04"
	BearerSchema         = "Bearer"
	RoleAdmin            = "ADMIN"
	RoleUseradmin        = "USERADMIN"
	RoleUser             = "USER"
	AuthHeaderKey        = "Authorization"
)
