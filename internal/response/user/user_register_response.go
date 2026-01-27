package user

type UserRegisterResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     uint   `json:"role"`
}
