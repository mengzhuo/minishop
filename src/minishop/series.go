package minishop

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

var letters = []rune("abcdefghijkmnpqrstuwxyzABCDEFGHJKLMNPQRSTUWXYZ")

func init() {
	rand.Seed(time.Now().Unix())
}
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type SerieNumber struct {
	Id        string `json:"id" form:"id" binding:"required"`
	ProductID string `json:"pid" form:"pid" binding:"required" bson:"pid"`
	Claimed   bool   `json:"claimed" form:"claimed"`
}

func getSerie(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusForbidden, "lack of id")
		return
	}

	data := strings.SplitN(id, "-", 2)
	var r SerieNumber
	err := SERIES.Find(bson.M{"id": data[1], "pid": data[0]}).One(&r)
	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}
	var p Product
	err = PRODUCT.Find(bson.M{"id": data[0]}).One(&p)

	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	c.HTML(http.StatusOK, "validate.html", gin.H{"r": r, "p": p})
}

func claimSerie(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusForbidden, "lack of id")
		return
	}

	data := strings.SplitN(id, "-", 2)
	var r SerieNumber
	err := SERIES.Find(bson.M{"id": data[1], "pid": data[0]}).One(&r)

	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	if r.Claimed {
		c.JSON(http.StatusConflict, r)
		return
	}

	r.Claimed = true
	SERIES.Update(bson.M{"id": r.Id, "pid": r.ProductID}, r)

	c.JSON(http.StatusOK, r)

}

type BulkAddSeriesForm struct {
	ProductID string `json:"pid" form:"pid" binding:"required"`
	Count     int    `json:"ids" form:"ids" binding:"required"`
}

func addSeries(c *gin.Context) {
	var form BulkAddSeriesForm
	err := c.Bind(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	ids := make([]string, 0)
	for i := 0; i < form.Count; i++ {
		id := randSeq(6)
		err := SERIES.Insert(SerieNumber{id, form.ProductID, false})
		if err != nil {
			c.JSON(500, err)
		}
		ids = append(ids, id)
	}
	c.HTML(http.StatusOK, "series_output.html", gin.H{"ids": ids, "pid": form.ProductID})
}
