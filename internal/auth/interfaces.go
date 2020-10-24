package auth

// LogoutEventer tells, when a sessionID gets revoked.
type LogoutEventer interface {
	LogoutEvent() ([]string, error)
}