package controller

import (
	"capstone/config"
	"capstone/middleware"
	"capstone/model"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/novalagung/gubrak/v2"

	"github.com/gorilla/websocket"

	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsconnuser   = map[int]*websocket.Conn{}
	wsconndoctor = map[int]*websocket.Conn{}
)

type Message struct {
	From    int
	To      int    `json:"to"`
	Message string `json:"message"`
	Sender string  `json:"sender"`
}

// func ConnectWSUser(c echo.Context) error {
// 	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
// 	userID := int(middleware.ExtractUserIdToken(token))

// 	currentconn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), c.Response().Header())

// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, "failure when connect websocket")
// 	}

// 	wsconnuser[userID] = currentconn

// 	go handleIO(currentconn, wsconndoctor, userID)
// 	return c.JSON(http.StatusAccepted, "success create connection")
// }

func ConnectWS(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	id, _ := middleware.ExtractToken(token)

	currentconn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), c.Response().Header())

	if err != nil {
		return c.JSON(http.StatusBadRequest, "failure when connect websocket")
	}

	// if role == "doctor" {
	// 	wsconndoctor[int(id)] = currentconn
	// 	go handleIO(currentconn, wsconnuser, int(id), role)
	// 	return c.JSON(http.StatusAccepted, "success create connection")
	// } else if role == "user" {
		wsconnuser[int(id)] = currentconn
		go handleIO(currentconn, wsconndoctor, int(id), "user")
		return c.JSON(http.StatusAccepted, "success create connection")
	// }
	
	// return c.JSON(http.StatusBadRequest, "failure when connect websocket")

}

func ConnectWSDoctor(c echo.Context) error {
	token := c.Param("Authorization")
	doctorID, _ := middleware.ExtractToken(token)

	currentconn, err2 := upgrader.Upgrade(c.Response().Writer, c.Request(), c.Response().Header())

	if err2 != nil {
		return c.JSON(http.StatusBadRequest, "failure when connect websocket")
	}

	wsconndoctor[int(doctorID)] = currentconn

	go handleIO(currentconn, wsconnuser, int(doctorID), "doctor")
	return c.JSON(http.StatusAccepted, "success create connection")
}

func handleIO(currentconn *websocket.Conn, connectionmapsender map[int]*websocket.Conn, from int, roles string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("error", fmt.Sprintf("%v", r))
		}
	}()

	chatroom := model.Chatroom{}
	for {
		message := Message{}

		err := currentconn.ReadJSON(&message)
		if roles == "user" {
			message.From = from
			chatroom.UserID = uint(message.From)
			chatroom.DoctorID = uint(message.To)
		} else if roles == "doctor" {
			message.From = from
			chatroom.DoctorID = uint(message.From)
			chatroom.UserID = uint(message.To)
		}
		message.From = from
		message.Sender = roles

		if err != nil {

			if strings.Contains(err.Error(), "websocket: close") {
				closeconn(currentconn, message)
				return
			}

			log.Println("error", err.Error())
			continue
		}

		if connectionmapsender[message.To] == nil {
			return
		}

		errsave := saveMessage(message, chatroom, roles)
		if !errsave {
			log.Println("error when save chat", errsave)
		}
		sendMessage(message, connectionmapsender[message.To])
	}
}

func sendMessage(message Message, destination *websocket.Conn) {
	destination.WriteJSON(message)
}

func closeconn(currentconn *websocket.Conn, message Message) {
	filtered := gubrak.From(wsconnuser).Reject(func(each *websocket.Conn) bool {
		return each == currentconn
	}).Result()
	wsconnuser[message.From] = filtered.(*websocket.Conn)
}

func createChatRoom(user model.User, doctor model.Doctor) (model.Chatroom, error) {

	err := config.DB.Model(&doctor).Where("id = ?", doctor.ID).Association("ChatwithUser").Append(&user)
	if err != nil {
		return model.Chatroom{}, err
	}

	result := config.DB.Model(model.Chatroom{}).Create(model.Chatroom{
		UserID:   user.ID,
		DoctorID: doctor.ID,
	}).Error

	if result != nil {
		return model.Chatroom{}, result
	}
	var Chatroom model.Chatroom
	config.DB.Model(model.Chatroom{}).Where("user_id = ? AND doctor_id = ?", user.ID, doctor.ID).Find(&Chatroom)
	return Chatroom, nil
}

func saveMessage(message Message, chatroom model.Chatroom, roles string) bool {
	chat := model.Chat{}
	chat.UserIDnoFK = int(chatroom.UserID)
	chat.DoctorIDnoFK = int(chatroom.DoctorID)
	chat.Content = message.Message
	chat.Sender = roles
	result := config.DB.Create(&chat)

	if result.RowsAffected < 1 {
		return false
	}
	return true
}

func GetAllChatHistory(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	id, role := middleware.ExtractToken(token)
	idDoctorUser := c.Param("id")
	var chat []model.Chat

	if role == "doctor" {
		//kayaknya error di bagian nama fieldnya
		config.DB.Model(&model.Chat{}).Where("doctor_idno_fk = ? AND user_idno_fk = ?", id, idDoctorUser).Find(&chat)
	} else if role == "user" {
		config.DB.Model(&model.Chat{}).Where("doctor_idno_fk = ? AND user_idno_fk = ?", idDoctorUser, id).Find(&chat)
	} else {
		return c.JSON(http.StatusInternalServerError, "failed get role")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all chat",
		"chat":    chat,
	})
}
