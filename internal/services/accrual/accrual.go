package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/store"
	"go.uber.org/zap"
)

func CheckAccrualAmount(order int) {
	orderCh <- strconv.Itoa(order)
}

func checkOrder(orderNum string) error {
	num, err := strconv.Atoi(orderNum)
	if err != nil {
		return err
	}

	order, err := store.GetOrderByNumber(num)
	if err != nil {
		logger.Error("Failed to get order by number", zap.Error(err), zap.String("orderNum", orderNum))
		return err
	}

	req := httpc.R()
	resp, err := req.Get(fmt.Sprintf(accrualServicePath, orderNum))
	if err != nil {
		logger.Error("Failed to get order by number: request error", zap.Error(err), zap.String("orderNum", orderNum))
		return err
	}

	if resp.StatusCode() == http.StatusNoContent {
		order.Status = "INVALID"
		err = store.UpdateOrder(order)
		if err != nil {
			logger.Error(
				"failed update order status",
				zap.String("error", err.Error()),
				zap.String("orderNum", orderNum),
				zap.String("status", order.Status),
			)

			return err
		}

		return errors.New("order not found in accrual system")
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		pause := resp.Header().Get("Retry-After")
		if pause == "" {
			pauseCh <- pause
		}

		orderCh <- orderNum
		return nil
	}

	if resp.StatusCode() == http.StatusOK && resp.Header().Get("Content-Type") == "application/json" {
		var accrualResponse AccrualResponse
		if err := json.Unmarshal(resp.Body(), &accrualResponse); err != nil {
			logger.Error("failed unmarshal accrual response", zap.String("error", err.Error()))
			return err
		}

		if accrualResponse.Status == store.OrderStatusProcessed || accrualResponse.Status == store.OrderStatusInvalid {
			order.Status = "PROCESSED"
			order.Accrual = accrualResponse.Accrual
			err = store.UpdateOrder(order)
			if err != nil {
				logger.Error(
					"failed update order status",
					zap.String("error", err.Error()),
					zap.String("orderNum", orderNum))
				return err
			}

			return nil
		}

		if accrualResponse.Status == store.OrderStatusProcessing {
			order.Status = "PROCESSING"
			err = store.UpdateOrder(order)
			if err != nil {
				logger.Error(
					"failed update order status",
					zap.String("error", err.Error()),
					zap.String("orderNum", orderNum))
			}
		}

		orderCh <- orderNum
		return nil
	}

	return errors.New("failed to check accrual amount: http request error")
}
