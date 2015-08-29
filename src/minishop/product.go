package minishop

import (
	"io"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var (
	PRODUCT_NOT_FOUND = gin.H{"err": "product not founded"}
)

type Product struct {
	Id          string    `json:"id" form:"id" binding:"required"`
	UpdateAt    time.Time `json:"update_at"`
	Factory     string    `json:"factory" form:"factory" binding:"required"`
	Version     string    `json:"version" form:"version" binding:"required"`
	Tags        []string  `json:"tags" form:"tags"`
	Title       string    `json:"title" form:"title"`
	Price       string    `json:"price" form:"price"`
	Images      []string  `json:"images" form:"images"`
	Link        string    `json:"link" form:"link"`
	Description string    `json:"description" form:"description"`
	ExtraInfo   string    `json:"extra_info" form:"extra_info"`
	On          bool      `json:"on" form:"on" binding:"required"`
}

func getProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	var p Product
	if err := PRODUCT.Find(bson.M{"id": id}).One(&p); err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}
	c.JSON(http.StatusOK, p)
}
func queryProduct(c *gin.Context) {

}
func deleteProduct(c *gin.Context) {

}

func updateProduct(c *gin.Context) {
	var p Product
	err := c.Bind(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if c.Param("id") != "" {
		p.Id = c.Param("id")
	}
	p.UpdateAt = time.Now()
	PRODUCT.Upsert(bson.M{"id": p.Id}, p)
	c.JSON(http.StatusOK, p)

}

func uploadAsset(c *gin.Context) {
	file, header, err := c.Request.FormFile("asset")
	filename := header.Filename
	glog.Infof("Upload file:%s", filename)
	out, err := os.Create("asset/" + filename)
	if err != nil {
		glog.Error(err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		glog.Error(err)
	}
}
