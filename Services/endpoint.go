package Services

type UserRequest struct {
	UserId int `json:"user_id"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}


