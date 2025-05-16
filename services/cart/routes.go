package cart

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/panjiasmoroart/gopher-ecom/services/auth"
	"github.com/panjiasmoroart/gopher-ecom/types"
	"github.com/panjiasmoroart/gopher-ecom/utils"
)

type Handler struct {
	store      types.ProductStore
	orderStore types.OrderStore
	userStore  types.UserStore
	db         *sql.DB
}

func NewHandler(
	store types.ProductStore,
	orderStore types.OrderStore,
	userStore types.UserStore,
	db *sql.DB,
) *Handler {
	return &Handler{
		store:      store,
		orderStore: orderStore,
		userStore:  userStore,
		db:         db,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	fmt.Println("userID >>> ", userID)

	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}
	// {[{1 10} {2 2}]}

	if len(cart.Items) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("cart is empty [or] null"))
		return
	}

	productIds, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// [1 2]

	// get products
	products, err := h.store.GetProductsByID(productIds)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// [
	//     {
	//         "id": 1,
	//         "product": "Sepatu",
	//         "description": "Lorem",
	//         "code": "cc",
	//         "price": 25000,
	//         "quantity": 1,
	//         "date": "2025-05-12 15:08:49 +0000 UTC"
	//     },
	//     {
	//         "id": 2,
	//         "product": "Laptop Lenovo",
	//         "description": "Lorem",
	//         "code": "dd",
	//         "price": 10000,
	//         "quantity": 5,
	//         "date": "2025-05-12 15:09:13 +0000 UTC"
	//     }
	// ]

	orderID, totalPrice, err := h.createOrder(products, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}
