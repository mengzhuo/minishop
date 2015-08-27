package minishop

import (
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

const (
	VERSION = "dev"
)

var (
	SESSION *mgo.Session
	DB      *mgo.Database
	SERIES  *mgo.Collection
	PRODUCT *mgo.Collection
	ACCOUNT gin.Accounts
)

func Index(c *gin.Context) {
	var Products []*Product
	err := PRODUCT.Find(bson.M{"on": true}).All(&Products)
	if err != nil {
		glog.Error(err)
	}
	c.HTML(http.StatusOK, "index.html", Products)
}

func DBInit(dbpath string) {
	var err error
	SESSION, err = mgo.Dial(dbpath)
	if err != nil {
		panic(err)
	}
	DB = SESSION.DB("minishop")
	SERIES = DB.C("series")
	PRODUCT = DB.C("product")
}

func AccountsInit() {

	admin := os.Getenv("MS_ADMIN")
	if admin == "" {
		admin = "admin"
	}
	passwd := os.Getenv("MS_PASSWD")
	if passwd == "" {
		panic("You have to set password")
	}

	ACCOUNT = gin.Accounts{admin: passwd}

}

func Serve(address, dbpath string) {

	DBInit(dbpath)
	AccountsInit()
	route := gin.Default()
	route.GET("/", Index)
	route.LoadHTMLGlob("templates/*.html")

	product := route.Group("/product", gin.BasicAuth(ACCOUNT))
	{
		product.GET(":id", getProduct)
		product.GET("", queryProduct)
		product.POST("", updateProduct)
		product.POST(":id", updateProduct)
		product.DELETE(":id", deleteProduct)
	}

	series := route.Group("/series", gin.BasicAuth(ACCOUNT))
	{
		series.POST("", addSeries)
	}

	validate := route.Group("/v")
	{
		validate.GET(":id", getSerie)
		validate.POST(":id", claimSerie)
	}

	route.Static("/static", "static")
	route.Static("/asset", "asset")
	route.Run(address)
}
