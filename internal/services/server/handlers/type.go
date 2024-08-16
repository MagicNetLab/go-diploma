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

type UserOrder struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

type UserOrdersResponse []UserOrder

type UserBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type UserWithdraw struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type WithdrawResponse []UserWithdraw
