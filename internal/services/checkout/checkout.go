package checkout

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bhlox/ecom/internal/db"
	"github.com/bhlox/ecom/internal/response"
	"github.com/bhlox/ecom/internal/types"
	"github.com/bhlox/ecom/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CheckoutQueries interface {
	CreateOrderItems(ctx context.Context, arg db.CreateOrderItemsParams) (db.OrderItem, error)
	CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.CreateOrderRow, error)
	GetProductPriceAndStock(ctx context.Context, id int32) (db.GetProductPriceAndStockRow, error)
	UpdateProductStock(ctx context.Context, arg db.UpdateProductStockParams) error
}

func (h *Handler) checkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(types.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "user ID not found in context or is invalid ID")
		return
	}

	var cart types.OrderPayload
	if err := utils.ParseJSONReq(r, &cart); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		ValidationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			fmt.Println("err is not that type")
			response.Error(w, 500, "wtf")
			return
		}
		var errors = make(map[string]string)
		for _, validationError := range ValidationErrors {
			errors[validationError.Field()] = validationError.Error()
		}
		fmt.Println("failed on validating payload")
		response.JSON(w, http.StatusNotAcceptable, errors)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	var orderTotals float64
	for _, item := range cart.OrderItems {
		priceAndStk, err := h.queries.GetProductPriceAndStock(r.Context(), item.ProductID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		if priceAndStk.Quantity < item.Quantity {
			response.Error(w, http.StatusBadRequest, fmt.Sprintf("productId:%v has a current stock of %v", item.ProductID, priceAndStk.Quantity))
			return
		} else {
			price, err := strconv.ParseFloat(priceAndStk.Price, 64)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, "Failed to parse price")
				return
			}
			orderItemTotal := price * float64(item.Quantity)
			orderTotals += orderItemTotal
		}
	}

	// payment here

	// payment end

	createdOrder, err := h.queries.CreateOrder(r.Context(), db.CreateOrderParams{
		UserID:  userID,
		Address: cart.Address,
		Total:   fmt.Sprintf("%.2f", orderTotals),
	})

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "something went wrong inserting order")
		return
	}

	var createdOrderItems = make([]db.OrderItem, 0)
	for _, item := range cart.OrderItems {
		priceAndStk, err := h.queries.GetProductPriceAndStock(r.Context(), item.ProductID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		stockLeft := priceAndStk.Quantity - item.Quantity
		if err := h.queries.UpdateProductStock(r.Context(), db.UpdateProductStockParams{
			Quantity: stockLeft,
			ID:       item.ProductID,
		}); err != nil {
			response.Error(w, http.StatusInternalServerError, "something went wrong updating stock")
			return
		}
		price, err := strconv.ParseFloat(priceAndStk.Price, 64)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "something went wrong pasing float")
			return
		}
		totalPrice := float64(item.Quantity) * price
		order, err := h.queries.CreateOrderItems(r.Context(), db.CreateOrderItemsParams{
			OrderID:   createdOrder.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     fmt.Sprintf("%.2f", totalPrice),
		})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "something went wrong creating order items")
			return
		}
		createdOrderItems = append(createdOrderItems, order)
	}
	// resp := map[string]any{
	// 	"totals":            orderTotals,
	// 	"createdOrderItems": createdOrderItems,
	// }
	resp := struct {
		Totals            float64        `json:"totals"`
		CreatedOrderItems []db.OrderItem `json:"createdOrderItems"`
	}{Totals: orderTotals,
		CreatedOrderItems: createdOrderItems,
	}
	response.JSON(w, http.StatusOK, resp)
}
