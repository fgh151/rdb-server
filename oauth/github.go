package oauth

type GithubUser struct {
	Login     string `json:"login"`
	Id        int    `json:"id"`
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}
