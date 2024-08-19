package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/order"
	"github.com/MagicNetLab/go-diploma/internal/services/user"
)

func CreateOrderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		userID, err := user.GetAuthUserID(r)
		if err != nil {
			logger.Error(fmt.Sprintf("error getting auth user id from cookie: %v", err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		num, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(fmt.Sprintf("create order error: failed reading body. %v", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if string(num) == "" {
			logger.Error(fmt.Sprintf("create order error: empty input. %v", err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		number, err := strconv.Atoi(string(num))
		if err != nil {
			logger.Error(fmt.Sprintf("create order error: invalid input. %v", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = order.CreateOrder(number, userID)
		if err != nil {
			if errors.Is(err, order.ErrorOrderAlreadyAddedByOtherUser) {
				logger.Error(fmt.Sprintf("create order error: order already added by other user"))
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			if errors.Is(err, order.ErrorOrderAlreadyAddedByUser) {
				logger.Error(fmt.Sprintf("create order error: order already added by current user"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		logger.Info(fmt.Sprintf("create order success: num %s, user %v", string(num), userID))
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
			logger.Error(fmt.Sprintf("error getting auth user id from cookie: %v", err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userOrders, err := order.GetUserOrders(userID)
		if err != nil {
			logger.Error(fmt.Sprintf("error getting user orders: %v", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
			logger.Error(fmt.Sprintf("error send user orders response: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			logger.Error(fmt.Sprintf("error getting auth user id from cookie: %v", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		result, err := order.GetBalanceByUserID(userID)
		if err != nil {
			logger.Error(fmt.Sprintf("error getting balance of user: %v", err))
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
			logger.Error(fmt.Sprintf("error send balance response: %v", err))
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
			logger.Error(fmt.Sprintf("error getting auth user id from cookie: %v", err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var withdrawRequest WithDrawRequest
		if err := json.NewDecoder(r.Body).Decode(&withdrawRequest); err != nil {
			logger.Error(fmt.Sprintf("failed to decode withdraw request: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !withdrawRequest.IsValid() {
			logger.Error(fmt.Sprintf("Withdraw request: invalid request params: %v", err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		balance, err := order.GetBalanceByUserID(userID)
		if err != nil {
			logger.Error(fmt.Sprintf("error getting balance of user: %v", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if balance.Current < withdrawRequest.Sum {
			logger.Error(fmt.Sprintf("Withdraw request error: insufficient balance. user %v", userID))
			http.Error(w, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
			return
		}

		err = order.CreateWithdraw(withdrawRequest.Order, withdrawRequest.Sum, userID)
		if err != nil {
			if errors.Is(err, order.ErrorIncorrectWithdrawNumber) {
				logger.Error(fmt.Sprintf("Withdraw request error: %v", err))
				http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
				return
			}

			logger.Error(fmt.Sprintf("error creating withdraw: %v", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		logger.Info(fmt.Sprintf("withdraw request success: num %s, user %v", withdrawRequest.Order, userID))
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
			logger.Error(fmt.Sprintf("error getting auth user id from cookie: %v", err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		result, err := order.GetWithdrawsByUserId(userID)
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
			logger.Error(fmt.Sprintf("error send withdraws response: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
