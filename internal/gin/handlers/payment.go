package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/billing"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/payment"
)

// POST /bills/pay
func PayBillHandler(p payment.Payment, b billing.Billing, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Payment.PayBill"

		log := log.With(
			slog.String("op", op),
		)

		idStr := strings.Trim(c.PostForm("id"), "/")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("cannot convert id to int64")
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "cannot convert id to int64",
			})
			return
		}

		log = log.With(
			slog.Int64("bill_id", id),
		)

		bill, err := b.GetBill(c, id)
		if err != nil {
			log.Error("cannot find bill")
			c.JSON(http.StatusNotFound, gin.H{
				"err": "cannot convert id to int64",
			})
			return
		}

		paymentStatus, err := p.ProcessPayment(c, bill.Amount)
		if err != nil {
			log.Error("cannot process payment")
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "cannot process payment",
			})
			return
		}

		if paymentStatus == payment.PAYMENT_OK {
			log.Info("payment proccessed")

			err = b.PayBill(c, id)
			if err != nil {
				log.Error("cannot pay the bill")
				c.JSON(http.StatusInternalServerError, gin.H{
					"err": "cannot pay the bill",
				})
				return
			}
		}

		log.Error("payment failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "payment failed",
		})
	}
}
