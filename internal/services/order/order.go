package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/bhlox/ecom/internal/db"
	"github.com/bhlox/ecom/internal/response"
	"github.com/bhlox/ecom/internal/types"
	"github.com/google/uuid"
)

type OrderQueries interface {
	CreateOrderItems(ctx context.Context, arg db.CreateOrderItemsParams) (db.OrderItem, error)
	CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.CreateOrderRow, error)
	GetOrder(ctx context.Context, id int32) ([]db.GetOrderRow, error)
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(types.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "user ID not found in context or is invalid ID")
		return
	}
	orderIdStr := r.PathValue("orderId")

	if orderIdStr == "" {
		response.JSON(w, http.StatusBadRequest, "no orderId param found")
		return
	}
	orderId, err := strconv.Atoi(orderIdStr)

	if err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid orderId format")
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	orderData, err := h.queries.GetOrder(r.Context(), int32(orderId))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	formattedOrder := struct {
		OrderDetails db.Order
		Items        []db.OrderItem
	}{
		OrderDetails: orderData[0].Order,
	}

	if formattedOrder.OrderDetails.UserID != userID {
		response.Error(w, http.StatusUnauthorized, "userIds do not match on the order")
		return
	}
	for _, order := range orderData {
		formattedOrder.Items = append(formattedOrder.Items, order.OrderItem)
	}
	// fmt.Println(order)
	response.JSON(w, http.StatusOK, formattedOrder)
}
