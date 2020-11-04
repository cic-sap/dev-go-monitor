module example

go 1.14

require (
	github.com/emicklei/go-restful v2.13.0+incompatible
	github.com/emicklei/go-restful-openapi v1.4.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/spec v0.19.8
	github.com/labstack/echo/v4 v4.1.16
	github.com/prometheus/client_golang v1.7.1
	github.com/cic-sap/dev-go-monitor v1.0.1
	github.com/slok/go-http-metrics v0.8.0
)

replace github.com/cic-sap/dev-go-monitor => ../
