package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	productpb "github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
)

type ProductContainer struct {
	Handler *handler.ProductHandler
}

func NewProductHandler(productClient productpb.ProductServiceClient) *ProductContainer {
	handler := handler.NewProductHandler(productClient)
	return &ProductContainer{handler}
}