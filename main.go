// Copyright Fauna, Inc.
// SPDX-License-Identifier: MIT-0

package main

import (
	"net/http"
	"os"

	"github/fauna-labs/go-gin-fly-io-starter/internal/utils"

	"github.com/fauna/fauna-go"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	client := fauna.NewClient(os.Getenv("FAUNA_SECRET_KEY"), fauna.DefaultTimeouts())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"greeting1": "Hello go-gin-fly-io-starter"})
	})

	r.GET("/read", func(c *gin.Context) {
		q, _ := fauna.FQL(`
		order.all() {
			orderName: .name,
			customer: .customer.firstName + " " + .customer.lastName,
			orderProducts {
				product: .product.name,
				price,
				quantity
			}
		}
		`, nil)
		res, err := client.Query(q)
		if err != nil {
			s := utils.GetErrorResponseStatusCode(err)
			c.JSON(s, err)
		} else {
			c.JSON(http.StatusOK, utils.GenerateResponse(res))
		}
	})

	r.Run()
}
