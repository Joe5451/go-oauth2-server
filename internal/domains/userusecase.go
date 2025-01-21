package domains

type UserUsecase interface {
	GenerateSocialProviderAuthUrl(provider, state, redirectUri string) (string, error)
	CreateUser(user User) User
	LoginSocialAccount(providerStr, authorizationCode, redirectUri string) (User, error)
	GetUserById(userId int64) User
	UpdateUserById(userId int64, user User)
}
