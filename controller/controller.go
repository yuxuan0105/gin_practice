package controller

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/yuxuan0105/gin_practice/api/v1"
	d "github.com/yuxuan0105/gin_practice/middleware/database"
	"github.com/yuxuan0105/gin_practice/pkg/setting"
)

type Controller struct {
	rt  *gin.Engine
	srv *http.Server
}

func NewController() *Controller {
	this := Controller{}
	//read flags
	confPath := ""
	flag.StringVar(&confPath, "c", "", "Configuration file path.")
	flag.Parse()
	//setup viper
	v, err := setting.GetSetting(confPath)
	if err != nil {
		log.Panicf("NewController: %s", err)
	}
	//setup database
	newdb, err := d.SetupDatabase(v)
	if err != nil {
		log.Panicf("NewController: %s", err)
	}

	//setup router
	this.rt = gin.Default()
	this.rt.Use(d.GetMiddlewareFunc(newdb))

	user := this.rt.Group("/api/v1/users")
	{
		user.GET("", v1.GetUsers)
		user.GET(":uid", v1.GetUserById)
		user.POST("", v1.SignUp)
		user.PATCH(":uid", v1.ModifyUserName)
		user.DELETE(":uid", v1.DeleteUser)
	}
	//setup server
	this.srv = &http.Server{
		Addr:    ":8080",
		Handler: this.rt,
	}

	return &this
}

func (this *Controller) RunServer() {
	go func() {
		// service connections
		if err := this.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server is running")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println()
	log.Println("Server is Closing ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := this.srv.Shutdown(ctx); err != nil {
		log.Fatal("Error Occur at Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}

func (this *Controller) ServeTestRequest(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	this.rt.ServeHTTP(w, req)
	return w
}
