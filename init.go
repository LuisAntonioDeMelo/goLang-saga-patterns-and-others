package main

import (
	"errors"
	"fmt"
)

type OrderID int

func CreateOrder() (OrderID, error) {
	fmt.Println("Creating order...")
	return 42, nil
}

func CancelOrder(id OrderID) error {
	fmt.Println("Cancelling order...")
	return nil
}

func ReserveInventory(id OrderID) error {
	fmt.Println("Reserving inventory...")
	return errors.New("inventory not available")
}

func ProcessPayment(id OrderID) error {
	fmt.Println("Processing payment...")
	return errors.New("card declined")
}

func RefundPayment(id OrderID) error {
	fmt.Println("Refunding payment...")
	return nil
}

func ReleaseStock(orderID OrderID) error {
	fmt.Printf("[Estoque] Liberando estoque do pedido %d...\n", orderID)
	return nil
}

// funcao que orquestra a saga de checkout
func RunCheckoutSaga() error {
	var (
		orderID          OrderID
		orderCreated     bool
		stockReserved    bool
		paymentProcessed bool
	)

	id, err := CreateOrder()
	if err != nil {
		fmt.Println("Failed to create order:", err)
	}

	orderID = id
	orderCreated = true

	if err = ReserveInventory(orderID); err != nil {
		if orderCreated {
			_ = CancelOrder(orderID)
		}
		return fmt.Errorf("falha ao reservar estoque: %w", err)
	}

	stockReserved = true

	if err := ProcessPayment(orderID); err != nil {
		if stockReserved {
			_ = ReleaseStock(orderID)
		}
		if orderCreated {
			_ = CancelOrder(orderID)
		}
		return fmt.Errorf("falha ao processar pagamento: %w", err)
	}

	paymentProcessed = true

	fmt.Printf("Saga concluída com sucesso! Pedido %d pago=%v, estoque reservado=%v\n",
		orderID, paymentProcessed, stockReserved)

	return nil

}

func main() {
	if err := RunCheckoutSaga(); err != nil {
		fmt.Println("Saga failed:", err)
	} else {
		fmt.Println("Saga completed successfully")
	}
}
