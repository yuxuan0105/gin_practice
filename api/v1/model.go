package v1

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuxuan0105/gin_practice/database"
	"github.com/yuxuan0105/gin_practice/setting"
)

type Model struct {
	rt  *gin.Engine
	db  *sql.DB
	srv *http.Server
}

func NewModel() (*Model, error) {
	//	md.router = setupRouter()
	//	md.db = setupDatabase()
	this := Model{}

	//setup viper
	v, err := setting.GetSetting()
	if err != nil {
		return nil, err
	}
	//setup database
	this.db, err = database.SetupDatabase(v)
	if err != nil {
		return nil, err
	}
	//setup router
	this.setupRouter()

	this.srv = &http.Server{
		Addr:    ":8080",
		Handler: this.rt,
	}

	return &this, nil
}

func (this *Model) setupRouter() {
	this.rt = gin.Default()

	user := this.rt.Group("/api/v1/user")
	{
		user.GET("", this.getUsers)
		user.GET(":uid", this.getUserById)
		user.POST("", this.addUser)
		user.PUT("", this.modifyUser)
		user.DELETE("", this.deleteUser)
	}

}

func (this *Model) RunServer() {
	go func() {
		// service connections
		if err := this.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := this.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}

func (this *Model) ServeTestRequest(method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	this.rt.ServeHTTP(w, req)
	return w
}
