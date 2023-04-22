package routes

import (
	"github.com/Prodigy00/restaurant-mgt/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/invoices", controllers.GetInvoices())
	incomingRoutes.GET("/invoices/:invoiceId", controllers.GetInvoice())
	incomingRoutes.POST("/invoices/signup", controllers.CreateInvoice())
	incomingRoutes.PATCH("/invoices/:invoiceId", controllers.UpdateInvoice())
}
