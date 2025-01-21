package domains

type UserRepository interface {
	Create(user User) User
	GetSocialAccountByProviderUserId(provider, providerUserId string) *SocialAccount
	CreateSocialAccount(account SocialAccount) SocialAccount
	GetUserById(id int64) User
	UpdateUserById(id int64, user User)
}
