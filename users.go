package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json="id"`
	Name string `sql:"unique_index" json="name"`
}

func NewUser(name string) (*User, error) {
	u := &User{
		Name: name,
	}

	if err := db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) Amount(match *Match) (tot int) {
	var movs []Movement
	db.Where("match_id = ? and user_id = ?", match.ID, u.ID).Find(movs)
	for _, m := range movs {
		tot += m.Amount
	}
	return
}

func (u *User) SetAmount(match *Match, amount int, date time.Time) {
	m := Movement{
		User:   u,
		Match:  match,
		Amount: amount,
		Date:   date,
	}
	db.Create(&m)
}

func (u *User) String() string {
	return u.Name
}

func (u *User) Win(w int) Play {
	return Play{u, w}
}

func attachUserAPI(router *gin.Engine) {
	api := router.Group("/api/users")
	{
		api.POST("/", func(c *gin.Context) {
			var form struct {
				Name string `form:"name" binding:"required"`
			}
			if c.Bind(&form) == nil {
				u, err := NewUser(form.Name)
				if err != nil {
					c.String(http.StatusBadRequest, err.Error())
					return
				}
				c.JSON(http.StatusOK, u)
			}
		})

		api.GET("/", func(c *gin.Context) {
			var users []User
			db.Find(&users)
			c.JSON(http.StatusOK, users)
		})

		api.GET("/:id", func(c *gin.Context) {
			user_id, _ := strconv.Atoi(c.Param("id"))
			var user User
			if err := db.First(&user, user_id).Error; err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			c.JSON(http.StatusOK, user)
		})

		api.DELETE("/:id", func(c *gin.Context) {
			user_id, _ := strconv.Atoi(c.Param("id"))
			var user User
			if err := db.First(&user, user_id).Error; err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			db.Delete(user)
			c.String(http.StatusOK, "Deleted")
		})

	}
}
