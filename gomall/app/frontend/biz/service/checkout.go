// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"strconv"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/checkout"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	rpccart "github.com/cloudwego/biz-demo/gomall/app/frontend/kitex_gen/cart"
	rpcproduct "github.com/cloudwego/biz-demo/gomall/app/frontend/kitex_gen/product"
	frontendutils "github.com/cloudwego/biz-demo/gomall/app/frontend/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type CheckoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCheckoutService(Context context.Context, RequestContext *app.RequestContext) *CheckoutService {
	return &CheckoutService{RequestContext: RequestContext, Context: Context}
}

func (h *CheckoutService) Run(req *checkout.CheckoutReq) (resp map[string]any, err error) {
	var items []map[string]string
	userId := frontendutils.GetUserIdFromCtx(h.Context)

	carts, err := rpc.CartClient.GetCart(h.Context, &rpccart.GetCartRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	var total float32
	for _, v := range carts.Items {
		p, err := rpc.ProductClient.GetProduct(h.Context, &rpcproduct.GetProductRequest{Id: v.ProductId})
		if err != nil {
			continue
		}
		items = append(items, map[string]string{
			"Name":    p.Name,
			"Price":   strconv.FormatFloat(float64(p.Price), 'f', 2, 64),
			"Picture": p.Picture,
			"Qty":     strconv.Itoa(int(v.Quantity)),
		})
		total += float32(v.Quantity) * p.Price
	}

	return utils.H{
		"title":    "Checkout",
		"items":    items,
		"cart_num": len(items),
		"total":    strconv.FormatFloat(float64(total), 'f', 2, 64),
	}, nil
}
