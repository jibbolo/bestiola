package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const COIN = 10

type Match struct {
	ID        int
	Players   []*User `gorm:"many2many:match_players;"`
	CreatedAt time.Time
	Pool      int
}

func NewMatch(players ...*User) (*Match, error) {
	m := &Match{
		Players:   players,
		CreatedAt: time.Now(),
		Pool:      0,
	}
	if err := db.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Match) PrintStatus() {
	fmt.Println("--------------")
	fmt.Printf("Pool: %v\n", m.Pool)
	for _, p := range m.Players {
		fmt.Printf("%v:\t%v\n", p.Name, p.Amount)
	}
	fmt.Println("--------------")
}

func (m *Match) GetPlayer(user_id int) (*User, bool) {
	for _, p := range m.Players {
		if user_id == p.ID {
			return p, true
		}
	}
	return nil, false
}

func (m *Match) SetPool() {
	m.Pool = len(m.Players) * COIN
	now := time.Now()
	for _, p := range m.Players {
		p.SetAmount(m, -COIN, now)
	}
	log.Printf("Pool set %v", m.Pool)
}

type Movement struct {
	ID      int
	Match   *Match
	MatchID int
	User    *User
	UserID  int
	Amount  int
	Date    time.Time
}

type Play struct {
	Player *User
	Won    int
}

func (m *Match) NewHand(plays ...Play) {

	if m.Pool == 0 {
		m.SetPool()
	}

	handPool := m.Pool
	wonMap := make(map[*User]int)
	now := time.Now()

	var winners int

	for _, play := range plays {
		if play.Won > 0 {
			wonMap[play.Player] = play.Won
			winners += 1
		} else {
			wonMap[play.Player] = 0
		}
	}

	if winners != 3 {
		m.Pool = 0
	}
	for p, w := range wonMap {
		if w == 0 {
			p.SetAmount(m, -handPool, now)
			m.Pool += handPool
			log.Printf("%v loses: %v Bestia!!", p, handPool)
		} else {
			if winners != 3 {
				delta := int(float64(handPool) / 3 * float64(w))
				p.SetAmount(m, delta, now)
				log.Printf("%v wins (%v): %v", p, w, delta)
			} else {
				log.Printf("%v patta", p)
			}

		}
		db.Save(&p)
	}
	db.Save(&m)
}

func attachMatchAPI(router *gin.Engine) {
	api := router.Group("/api/matches")
	{
		api.POST("/", func(c *gin.Context) {

			var form struct {
				UserIDs []int `form:"users[]" binding:"required"`
			}
			if c.Bind(&form) == nil {
				var users []*User
				if err := db.Where(form.UserIDs).Find(&users).Error; err != nil {
					c.String(http.StatusBadRequest, err.Error())
					return
				}
				if len(users) == 0 {
					c.String(http.StatusBadRequest, "Missing players")
					return
				}
				m, err := NewMatch(users...)
				if err != nil {
					c.String(http.StatusBadRequest, err.Error())
					return
				}
				c.JSON(http.StatusOK, m)
			}
		})

		api.GET("/", func(c *gin.Context) {
			var matches []Match
			db.Preload("Players").Find(&matches)
			c.JSON(http.StatusOK, matches)
		})

		api.GET("/:id/movements", func(c *gin.Context) {
			m_id, _ := strconv.Atoi(c.Param("id"))
			var match Match
			if err := db.Preload("Players").First(&match, m_id).Error; err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			var movs []Movement
			db.Preload("User").Where("match_id = ?", m_id).Find(&movs)
			c.JSON(http.StatusOK, movs)
		})

		api.GET("/:id", func(c *gin.Context) {
			m_id, _ := strconv.Atoi(c.Param("id"))
			var match Match
			if err := db.Preload("Players").First(&match, m_id).Error; err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			c.JSON(http.StatusOK, match)
		})

		api.DELETE("/:id", func(c *gin.Context) {
			m_id, _ := strconv.Atoi(c.Param("id"))
			var match Match
			if err := db.First(&match, m_id).Error; err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			db.Delete(match)
			c.String(http.StatusOK, "Deleted")
		})

		api.POST("/:id", func(c *gin.Context) {

			var form struct {
				Plays []struct {
					UserID int `json:"user_id" binding:"required"`
					Won    int `json:"won" binding:"required"`
				} `json:"plays" binding:"required"`
			}

			m_id, _ := strconv.Atoi(c.Param("id"))
			var match Match
			if err := db.Preload("Players").First(&match, m_id).Error; err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}

			if c.Bind(&form) == nil {

				var plays []Play
				for _, p := range form.Plays {
					if user, ok := match.GetPlayer(p.UserID); ok == true {
						plays = append(plays, user.Win(p.Won))
					}
				}
				match.NewHand(plays...)
				c.JSON(http.StatusOK, "ok")
				return
			}

			c.JSON(http.StatusBadRequest, "not ok")
		})

	}
}
