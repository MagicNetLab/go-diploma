package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/store"
)

func CheckAccrualAmount(order int) {
	orderCh <- strconv.Itoa(order)
}

func processOrderAccrual(orderNum string) error {
	num, err := strconv.Atoi(orderNum)
	if err != nil {
		return err
	}

	order, err := store.GetOrderByNumber(num)
	if err != nil {
		args := map[string]interface{}{"error": err.Error(), "orderNum": orderNum}
		logger.Error("Failed to get order by number", args)
		return err
	}

	req := httpc.R()
	resp, err := req.Get(fmt.Sprintf(accrualServicePath, orderNum))
	if err != nil {
		args := map[string]interface{}{"error": err.Error(), "orderNum": orderNum}
		logger.Error("Failed to get order by number: request error", args)
		return err
	}

	switch resp.StatusCode() {
	case http.StatusNoContent:
		order.Status = "INVALID"
		if err = store.UpdateOrder(order); err != nil {
			args := map[string]interface{}{"error": err.Error(), "orderNum": orderNum, "status": order.Status}
			logger.Error("failed update order status", args)
			return err
		}
		return errors.New("order not found in accrual system")
	case http.StatusTooManyRequests:
		if pause := resp.Header().Get("Retry-After"); pause != "" {
			pauseCh <- pause
		}
		orderCh <- orderNum
		return nil
	case http.StatusOK:
		if resp.Header().Get("Content-Type") == "application/json" {
			var accrualResponse AccrualResponse
			if err := json.Unmarshal(resp.Body(), &accrualResponse); err != nil {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("failed unmarshal accrual response", args)
				return err
			}

			if accrualResponse.Status == store.OrderStatusProcessed || accrualResponse.Status == store.OrderStatusInvalid {
				order.Status = "PROCESSED"
				order.Accrual = accrualResponse.Accrual
				err = store.UpdateOrder(order)
				if err != nil {
					args := map[string]interface{}{"error": err.Error(), "orderNum": orderNum}
					logger.Error("failed update order status", args)
					return err
				}
				return nil
			}

			if accrualResponse.Status == store.OrderStatusProcessing {
				order.Status = "PROCESSING"
				err = store.UpdateOrder(order)
				if err != nil {
					args := map[string]interface{}{"error": err.Error(), "orderNum": orderNum}
					logger.Error("failed update order status", args)
				}
			}

			orderCh <- orderNum
			return nil
		}
	}

	return errors.New("failed to check accrual amount: http request error")
}
