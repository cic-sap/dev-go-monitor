module github.com/cic-sap/dev-go-monitor

go 1.14

require (
	github.com/emicklei/go-restful v2.13.0+incompatible
	github.com/gin-gonic/gin v1.7.7
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/imroc/req v0.3.0
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/echo/v4 v4.6.1
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/prometheus/client_golang v1.11.1
	//github.com/zsais/go-gin-prometheus v0.1.0
	github.com/slok/go-http-metrics v0.10.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	golang.org/x/text v0.3.6 => golang.org/x/text v0.3.7
)
