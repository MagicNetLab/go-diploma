package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/order"
	"github.com/MagicNetLab/go-diploma/internal/services/user"
	"go.uber.org/zap"
)

func CreateOrderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		userID, err := user.GetAuthUserID(r)
		if err != nil {
			logger.Error("error getting auth user id from cookie", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		num, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("create order error: failed reading body", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if string(num) == "" {
			logger.Error("create order error: empty input. %v", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		number, err := strconv.Atoi(string(num))
		if err != nil {
			logger.Error("create order error: invalid input", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = order.CreateOrder(number, userID)
		if err != nil {
			if errors.Is(err, order.ErrorOrderAlreadyAddedByOtherUser) {
				logger.Error("create order error: order already added by other user")
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			if errors.Is(err, order.ErrorOrderAlreadyAddedByUser) {
				logger.Error("create order error: order already added by current user")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		logger.Info("create order success", zap.Int("userID", userID), zap.String("number", string(num)))
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Accepted"))
	}
}

func OrderListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		userID, err := user.GetAuthUserID(r)
		if err != nil {
			logger.Error("error getting auth user id from cookie", zap.Error(err))
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))

			return
		}

		userOrders, err := order.GetUserOrders(userID)
		if err != nil {
			logger.Error("error getting user orders", zap.Error(err))
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

			return
		}

		var response UserOrdersResponse
		for _, userOrder := range userOrders {
			o := UserOrder{
				Number:     strconv.Itoa(userOrder.Number),
				Status:     userOrder.Status,
				Accrual:    userOrder.Accrual,
				UploadedAt: userOrder.UploadedAt,
			}
			response = append(response, o)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if len(response) == 0 {
			_, err = w.Write([]byte("[]"))
		} else {
			err = json.NewEncoder(w).Encode(response)
		}

		if err != nil {
			logger.Error("error send user orders response", zap.Error(err))
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
	}
}

func BalanceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		userID, err := user.GetAuthUserID(r)
		if err != nil {
			logger.Error("error getting auth user id from cookie", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		result, err := order.GetBalanceByUserID(userID)
		if err != nil {
			logger.Error("error getting balance of user", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response := UserBalanceResponse{
			Current:   result.Current,
			Withdrawn: result.Withdrawn,
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("error send balance response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func WithdrawRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		userID, err := user.GetAuthUserID(r)
		if err != nil {
			logger.Error("error getting auth user id from cookie", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var withdrawRequest WithDrawRequest
		if err := json.NewDecoder(r.Body).Decode(&withdrawRequest); err != nil {
			logger.Error("failed to decode withdraw request", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !withdrawRequest.IsValid() {
			logger.Error(
				"Withdraw request: invalid request params",
				zap.String("order", withdrawRequest.Order),
				zap.Float64("sum", withdrawRequest.Sum))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		balance, err := order.GetBalanceByUserID(userID)
		if err != nil {
			logger.Error("error getting balance of user", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if balance.Current < withdrawRequest.Sum {
			logger.Error("Withdraw request error: insufficient balance.", zap.Int("user", userID))
			http.Error(w, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
			return
		}

		err = order.CreateWithdraw(withdrawRequest.Order, withdrawRequest.Sum, userID)
		if err != nil {
			if errors.Is(err, order.ErrorIncorrectWithdrawNumber) {
				logger.Error("Withdraw request error", zap.Error(err))
				http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
				return
			}

			logger.Error("error creating withdraw", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		logger.Info(
			"withdraw request success: num %s, user %v",
			zap.String("order", withdrawRequest.Order),
			zap.Int("userID", userID))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func WithdrawListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		userID, err := user.GetAuthUserID(r)
		if err != nil {
			logger.Error("error getting auth user id from cookie", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		result, err := order.GetWithdrawsByUserID(userID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var response WithdrawResponse
		for _, withdraw := range result {
			w := UserWithdraw{
				Order:       withdraw.Order,
				Sum:         withdraw.Sum,
				ProcessedAt: withdraw.ProcessedAt,
			}
			response = append(response, w)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("error send withdraws response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
