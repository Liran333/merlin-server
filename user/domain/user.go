package domain

// user
type User struct {
	Id      string
	Email   Email
	Account Account

	Bio      Bio
	AvatarId AvatarId

	Version int
}

type UserInfo struct {
	Account  Account
	AvatarId AvatarId
}
