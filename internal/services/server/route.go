package server

import (
	"github.com/MagicNetLab/go-diploma/internal/services/server/handlers"
	"github.com/go-chi/chi/v5"
)

func getRoute() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/user/register", mw(handlers.UserRegisterHandler(), mwGuestPost))
	r.Post("/api/user/login", mw(handlers.UserLoginHandler(), mwGuestPost))
	r.Post("/api/user/orders", mw(handlers.CreateOrderHandler(), mwAuthorizedPost))
	r.Get("/api/user/orders", mw(handlers.OrderListHandler(), mwAuthorizedGet))
	r.Get("/api/user/balance", mw(handlers.BalanceHandler(), mwAuthorizedGet))
	r.Post("/api/user/balance/withdraw", mw(handlers.WithdrawRequestHandler(), mwAuthorizedPost))
	r.Get("/api/user/withdrawals", mw(handlers.WithdrawListHandler(), mwAuthorizedGet))

	return r
}
