# dev-git-monitor-libs

golang export metrics 

go.mod
```
github.com/cic-sap/dev-go-monitor v1.0.7
```

example of gin framework with custom prometheus Registry
```go
package main

import (
	"fmt"
	"github.com/cic-sap/dev-go-monitor/monitor"
	"github.com/cic-sap/dev-go-monitor/plugin-gin"
	"github.com/gin-gonic/gin"
	req "github.com/imroc/req"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var reg = prometheus.NewRegistry()

var c1 = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "c1",
})

func main() {
	reg.MustRegister(c1)

	r := gin.Default()

	plugin_gin.Init(r, monitor.WithPath("/metrics"),
		monitor.WithGinHandle(func(c *gin.Context) string {
			if c.FullPath() == "/hi/:id" {
				return "/hi/" + strings.ToLower(c.Param("id"))
			}
			return c.FullPath()
		}),
		monitor.WithGatherer(reg),
		monitor.WithGatherer(prometheus.DefaultGatherer),
	)

	r.GET("/", func(c *gin.Context) {
		c1.Add(1.0)
		time.Sleep(time.Second * time.Duration(rand.Int31n(3)))
		c.String(http.StatusOK, "hello world")
	})
	r.GET("/500", func(c *gin.Context) {

		c.String(http.StatusOK,
			fmt.Sprintf("hello world:%d", 2/rand.Intn(2)))
	})
	r.GET("/hi/:id", func(c *gin.Context) {

		//log.Println("path:", c.FullPath())

		id := c.Param("id")
		c.String(http.StatusOK,
			fmt.Sprintf("hello world:%s", id))
	})
	r.GET("/info/*any", func(c *gin.Context) {

		//log.Println("path:", c.FullPath(), c.Param("any"))
		req.Get("https://httpbin.org/" + c.Param("any"))
		c.String(http.StatusOK,
			fmt.Sprintf("hello info:%s", c.FullPath()))
	})
	r.Run(":8094")
}


```

echo framework
```go
package main

import (
	"github.com/labstack/echo/v4"
	plugin_echo "github.com/cic-sap/dev-go-monitor/plugin-echo"
	"log"
	"net/http"
)

const (
	srvAddr = "127.0.0.1:8480"
)

func main() {
	// Create our middleware factory with the default settings.

	// Create Echo instance and global middleware.
	e := echo.New()

	// Add our handler.
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	e.GET("/500", func(c echo.Context) error {
		return c.String(http.StatusInternalServerError, "ERR 500")
	})
	plugin_echo.Init(e)

	log.Printf("server listening at %s", srvAddr)
	if err := http.ListenAndServe(srvAddr, e); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}

```
gin framework

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/cic-sap/dev-go-monitor/monitor"
	"github.com/cic-sap/dev-go-monitor/plugin-gin"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()

	plugin_gin.Init(r, monitor.WithPath("/metrics"))

	r.GET("/", func(c *gin.Context) {
		time.Sleep(time.Second * time.Duration(rand.Int31n(3)))
		c.String(http.StatusOK, "hello world")
	})
	r.GET("/500", func(c *gin.Context) {

		c.String(http.StatusOK,
			fmt.Sprint("hello world:%d", 2/rand.Intn(2)))
	})

	r.Run(":8094")
}

```

go-restful
```go
package main

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/cic-sap/dev-go-monitor/monitor"
	plugin_rest "github.com/cic-sap/dev-go-monitor/plugin-rest"
	"log"
	"net/http"
)

func main() {
	u := UserResource{map[string]User{}}
	restful.DefaultContainer.Add(u.WebService())

	plugin_rest.Init(restful.DefaultContainer, monitor.WithPath("/metrics"))

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs/?url=http://localhost:8080/apidocs.json
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("/Users/emicklei/Projects/swagger-ui/dist"))))

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// UserResource is the REST layer to the User domain
type UserResource struct {
	// normally one would use DAO (data access object)
	users map[string]User
}

// WebService creates a new service that can handle REST requests for User resources.
func (u UserResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	tags := []string{"users"}

	ws.Route(ws.GET("/").To(u.findAllUsers).
		// docs
		Doc("get all users").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]User{}).
		Returns(200, "OK", []User{}))

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		// docs
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("integer").DefaultValue("1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(User{}). // on the response
		Returns(200, "OK", User{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{user-id}").To(u.updateUser).
		// docs
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(User{})) // from the request

	ws.Route(ws.PUT("").To(u.createUser).
		// docs
		Doc("create a user").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
		// docs
		Doc("delete a user").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	return ws
}

// GET http://localhost:8080/users
//
func (u UserResource) findAllUsers(request *restful.Request, response *restful.Response) {
	list := []User{}
	for _, each := range u.users {
		list = append(list, each)
	}
	response.WriteEntity(list)
}

// GET http://localhost:8080/users/1
//
func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	usr := u.users[id]
	if len(usr.ID) == 0 {
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa Raspberry</Name></User>
//
func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = *usr
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa</Name></User>
//
func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr := User{ID: request.PathParameter("user-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = usr
		response.WriteHeaderAndEntity(http.StatusCreated, usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/users/1
//
func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "UserService",
			Description: "Resource for managing Users",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "john",
					Email: "john@doe.rp",
					URL:   "http://johndoe.org",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "users",
		Description: "Managing users"}}}
}

// User is just a sample type
type User struct {
	ID   string `json:"id" description:"identifier of the user"`
	Name string `json:"name" description:"name of the user" default:"john"`
	Age  int    `json:"age" description:"age of the user" default:"21"`
}

```

net/http
```go
package main

import (
	"fmt"
	pluginsimple "github.com/cic-sap/dev-go-monitor/plugin-simple"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})
	h := pluginsimple.Init(mux)
	http.ListenAndServe(":8290", h)
}

```
