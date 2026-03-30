package room

import (
	"net/http"
	"ws/src/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Controller struct {
	Repo *Repository
}

func NewController(r *Repository) *Controller {
	return &Controller{Repo: r}
}

func (ctrl *Controller) Create(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}

	c.BindJSON(&input)

	userIDHex := c.MustGet(auth.UserIDKey).(string)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	room, err := ctrl.Repo.createRoom(input.Name, userID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error(),
	})
	}
	c.JSON(http.StatusOK, room)
}

func (ctrl *Controller) Get(c *gin.Context) {
	roomIDHex := c.Param("id")
	roomID, _ := bson.ObjectIDFromHex(roomIDHex)
	room, err := ctrl.Repo.getRoom(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (ctrl *Controller) List(c *gin.Context) {	
	cursor, err := ctrl.Repo.Rooms.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list rooms"})
		return
	}
	defer cursor.Close(c)

	var rooms []Room
	for cursor.Next(c) {
		var room Room	
		if err := cursor.Decode(&room); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode room"})
			return
		}
		rooms = append(rooms, room)
	}	
	c.JSON(http.StatusOK, rooms)
}