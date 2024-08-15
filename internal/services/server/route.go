package server

import (
	"github.com/MagicNetLab/go-diploma/internal/services/server/handlers"
	"github.com/go-chi/chi/v5"
)

func getRoute() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/user/register", mwGuest(handlers.UserRegisterHandler()))
	r.Post("/api/user/login", mwGuest(handlers.UserLoginHandler()))
	r.Post("/api/user/orders", mwAuthorized(handlers.CreateOrderHandler()))
	r.Get("/api/user/orders", mwAuthorized(handlers.OrderListHandler()))
	r.Get("/api/user/balance", mwAuthorized(handlers.BalanceHandler()))
	r.Post("/api/user/balance/withdraw", mwAuthorized(handlers.WithdrawRequestHandler()))

	r.Get("/api/user/withdrawals", mwAuthorized(handlers.WithdrawListHandler()))

	return r
}
