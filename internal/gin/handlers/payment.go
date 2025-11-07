package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/billing"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/payment"
	"github.com/iskanye/utilities-payment-utils/pkg/logger"
)

// POST /bills/pay
func PayBillHandler(p payment.Payment, b billing.Billing, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Payment.PayBill"

		idStr := c.PostForm("id")

		log := log.With(
			slog.String("op", op),
			slog.String("bill_id", idStr),
		)

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			if idStr == "" {
				log.Error("bill_id required", logger.Err(err))
				c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
				return
			}
			log.Error("cannot convert id to int64", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "cannot convert id to int64",
			})
			return
		}

		bill, err := b.GetBill(c, id)
		if err != nil {
			log.Error("cannot find bill", logger.Err(err))
			c.JSON(http.StatusNotFound, gin.H{
				"err": "cannot find bill",
			})
			return
		}

		paymentStatus, err := p.ProcessPayment(c, bill.Amount)
		if err != nil {
			log.Error("cannot process payment", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "cannot process payment",
			})
			return
		}

		if paymentStatus == payment.PAYMENT_OK {
			log.Info("payment proccessed")

			err = b.PayBill(c, id)
			if err != nil {
				log.Error("cannot pay the bill", logger.Err(err))
				c.JSON(http.StatusInternalServerError, gin.H{
					"err": "cannot pay the bill",
				})
				return
			}

			c.Status(http.StatusNoContent)
			return
		}

		log.Error("payment failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "payment failed",
		})
	}
}
