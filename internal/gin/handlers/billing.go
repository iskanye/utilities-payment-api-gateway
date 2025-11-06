package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/billing"
	"github.com/iskanye/utilities-payment-utils/pkg/logger"
)

// POST /admin/bills
func AddBillHandler(b billing.Billing, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Billing.AddBill"

		address := c.PostForm("address")
		amountStr := c.PostForm("amount")
		userIDStr := c.PostForm("user_id")

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
					"err": "amount required",
				})
				return
			}
			log.Error("cant convert amount to int", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "cant convert amount to int",
			})
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			if userIDStr == "" {
				log.Error("user_id required", logger.Err(err))
				c.JSON(http.StatusBadRequest, gin.H{
					"err": "user_id required",
				})
				return
			}
			log.Error("cant convert user_id to int64", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "cant convert user_id to int64",
			})
			return
		}

		billID, err := b.AddBill(c, address, amount, userID)
		if err != nil {
			log.Error("failed to login user", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Info("success")
		c.JSON(http.StatusOK, gin.H{
			"id": billID,
		})
	}
}

// GET /bills/:id
func GetBillHandler(b billing.Billing, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Billing.GetBill"

		billIDStr := c.Param("id")

		log := log.With(
			slog.String("op", op),
			slog.String("bill_id", billIDStr),
		)

		log.Info("attempting to get bill")

		billID, err := strconv.ParseInt(billIDStr, 10, 64)
		if err != nil {
			if billIDStr == "" {
				log.Error("bill_id required", logger.Err(err))
				c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
				return
			}
			log.Error("cant convert bill_id to int64", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		bill, err := b.GetBill(c, billID)
		if err != nil {
			log.Error("failed to get bill", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Info("success")
		c.JSON(http.StatusOK, bill)
	}
}

// GET /bills
func GetBillsHandler(b billing.Billing, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Billing.GetBills"

		userIDStr := c.Query("user_id")

		log := log.With(
			slog.String("op", op),
			slog.String("user_id", userIDStr),
		)

		log.Info("attempting to get bills")

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			if userIDStr == "" {
				log.Error("user_id required", logger.Err(err))
				c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
				return
			}
			log.Error("cant convert user_id to int64", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		bills, err := b.GetBills(c, userID)
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
