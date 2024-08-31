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
)

func CreateOrderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := user.GetAuthUserID(r)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error getting auth user id from cookie", args)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		num, err := io.ReadAll(r.Body)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("create order error: failed reading body", args)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if string(num) == "" {
			args := map[string]interface{}{"error": "body is empty"}
			logger.Error("create order error: empty input. %v", args)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		number, err := strconv.Atoi(string(num))
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("create order error: invalid input", args)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = order.CreateOrder(number, userID)
		if err != nil {
			if errors.Is(err, order.ErrorOrderAlreadyAddedByOtherUser) {
				logger.Error("create order error: order already added by other user", nil)
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			if errors.Is(err, order.ErrorOrderAlreadyAddedByUser) {
				logger.Error("create order error: order already added by current user", nil)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}

			if errors.Is(err, order.ErrorIncorrectOrderNumber) {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("create order error: invalid input.", args)
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte(err.Error()))
				return
			}

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		args := map[string]interface{}{"userID": userID, "number": string(num)}
		logger.Info("create order success", args)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Accepted"))
	}
}

func OrderListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := user.GetAuthUserID(r)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error getting auth user id from cookie", args)
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))

			return
		}

		userOrders, err := order.GetUserOrders(userID)
		if err != nil {
			args := map[string]interface{}{"error": err.Error(), "userID": userID}
			logger.Error("error getting user orders", args)
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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error send user orders response", args)
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
	}
}

func BalanceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := user.GetAuthUserID(r)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error getting auth user id from cookie", args)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		result, err := order.GetBalanceByUserID(userID)
		if err != nil {
			args := map[string]interface{}{"error": err.Error(), "userID": userID}
			logger.Error("error getting balance of user", args)
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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error send balance response", args)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func WithdrawRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := user.GetAuthUserID(r)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error getting auth user id from cookie", args)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var withdrawRequest WithDrawRequest
		if err := json.NewDecoder(r.Body).Decode(&withdrawRequest); err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed to decode withdraw request", args)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !withdrawRequest.IsValid() {
			args := map[string]interface{}{"order": withdrawRequest.Order, "sum": withdrawRequest.Sum}
			logger.Error("Withdraw request: invalid request params", args)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		balance, err := order.GetBalanceByUserID(userID)
		if err != nil {
			args := map[string]interface{}{"error": err.Error(), "userID": userID}
			logger.Error("error getting balance of user", args)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if balance.Current < withdrawRequest.Sum {
			args := map[string]interface{}{"user": userID, "balance": balance.Current, "sum": withdrawRequest.Sum}
			logger.Error("Withdraw request error: insufficient balance.", args)
			http.Error(w, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
			return
		}

		err = order.CreateWithdraw(withdrawRequest.Order, withdrawRequest.Sum, userID)
		if err != nil {
			if errors.Is(err, order.ErrorIncorrectWithdrawNumber) {
				args := map[string]interface{}{"error": "incorrect Withdraw number"}
				logger.Error("Withdraw request error", args)
				http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
				return
			}

			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error creating withdraw", args)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		args := map[string]interface{}{"userID": userID, "order": withdrawRequest.Order}
		logger.Info("withdraw request success: num %s, user %v", args)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func WithdrawListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := user.GetAuthUserID(r)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error getting auth user id from cookie", args)
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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error send withdraws response", args)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
