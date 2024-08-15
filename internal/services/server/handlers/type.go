package handlers

type RegisterUserRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r *RegisterUserRequest) IsValid() bool {
	return r.Login != "" && r.Password != ""
}

type UserLoginRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r *UserLoginRequest) IsValid() bool {
	return r.Login != "" && r.Password != ""
}

type WithDrawRequest struct {
	Order string  `json:"order,omitempty"`
	Sum   float64 `json:"sum,omitempty"`
}

func (r *WithDrawRequest) IsValid() bool {
	return r.Order != "" && r.Sum > 0
}
