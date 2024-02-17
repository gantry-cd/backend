package controller

import (
	"fmt"
	"net/http"

	v1 "github.com/gantrycd/backend/proto/metric"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UsageController struct {
	resource v1.ResourceWatcherClient
}

func New(resource v1.ResourceWatcherClient) *UsageController {
	return &UsageController{
		resource: resource,
	}
}

func (uc *UsageController) Usage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resource, err := uc.resource.GetResource(ctx, &emptypb.Empty{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// metrics取れない場合
	if resource.IsDisable {
		// TODO:後でやる
		fmt.Printf("unsuport")
	}

}
