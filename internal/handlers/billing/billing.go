package billing

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/billing"
	"github.com/iskanye/utilities-payment/pkg/logger"
)

func AddBillHandler(b billing.Billing, log *slog.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		const op = "Billing.AddBill"

		address := c.Query("address")
		amountStr := c.Query("amount")

		log := log.With(
			slog.String("op", op),
			slog.String("address", address),
		)

		log.Info("attempting to create bill")

		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			if amountStr == "" {
				log.Error("amount required", logger.Err(err))
				c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
				return
			}
			log.Error("cant convert amount to int", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		billId, err := b.AddBill(c, address, amount)
		if err != nil {
			log.Error("failed to login user", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Info("success")
		c.JSON(http.StatusOK, gin.H{
			"id": billId,
		})
	}
}

func GetBillsHandler(b billing.Billing, log *slog.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		const op = "Billing.GetBills"

		address := c.Query("address")

		log := log.With(
			slog.String("op", op),
			slog.String("address", address),
		)

		log.Info("attempting to get bills")

		bills, err := b.GetBills(c, address)
		if err != nil {
			log.Error("failed to get bills", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Info("success")
		c.JSON(http.StatusOK, gin.H{
			"bills": bills,
		})
	}
}
