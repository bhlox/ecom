package product

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
)

type StoreQueries interface {
	CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error)
	GetAllProducts(ctx context.Context) ([]db.Product, error)
	GetProduct(ctx context.Context, id int32) (db.Product, error)
	GetProductPriceAndStock(ctx context.Context, id int32) (db.GetProductPriceAndStockRow, error)
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var data types.CreateProductPayload
	if err := utils.ParseJSONReq(r, &data); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate.Struct(data); err != nil {
		// fmt.Println(err.Error())
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

	product, err := h.queries.CreateProduct(r.Context(), db.CreateProductParams{
		Name:        data.Name,
		Description: data.Description,
		Image:       data.Image,
		Price:       fmt.Sprintf("%.2f", data.Price),
		Quantity:    int32(data.Quantity),
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, product)
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "no id parameters found")
		return
	}
	numberId, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to convert type to int")
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	product, err := h.queries.GetProduct(r.Context(), int32(numberId)) //#nosec G109
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, product)
}

func (h *Handler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	products, err := h.queries.GetAllProducts(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, products)
}
