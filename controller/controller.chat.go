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
	id, role := middleware.ExtractToken(token)

	currentconn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), c.Response().Header())

	if err != nil {
		return c.JSON(http.StatusBadRequest, "failure when connect websocket")
	}

	if role == "doctor" {
		wsconndoctor[int(id)] = currentconn
		go handleIO(currentconn, wsconnuser, int(id), role)
		return c.JSON(http.StatusAccepted, "success create connection")
	} else if role == "user" {
		wsconnuser[int(id)] = currentconn
		go handleIO(currentconn, wsconndoctor, int(id), role)
		return c.JSON(http.StatusAccepted, "success create connection")
	}
	return c.JSON(http.StatusBadRequest, "failure when connect websocket")

}

// func ConnectWSDoctor(c echo.Context) error {
// 	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
// 	doctorID, err := middleware.ExtractDocterIdToken(token)

// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]interface{}{
// 			"message": "failed when cast jwt",
// 			"error":   err,
// 		})
// 	}

// 	currentconn, err2 := upgrader.Upgrade(c.Response().Writer, c.Request(), c.Response().Header())

// 	if err2 != nil {
// 		return c.JSON(http.StatusBadRequest, "failure when connect websocket")
// 	}

// 	wsconndoctor[int(doctorID)] = currentconn

// 	go handleIO(currentconn, wsconnuser, int(doctorID), role)
// 	return c.JSON(http.StatusAccepted, "success create connection")
// }

func handleIO(currentconn *websocket.Conn, connectionmapsender map[int]*websocket.Conn, from int, roles string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("error", fmt.Sprintf("%v", r))
		}
	}()

	chatroom := model.ChatRoom{}
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

func createChatRoom(user model.User, doctor model.Doctor) (model.ChatRoom, error) {

	err := config.DB.Model(&doctor).Where("id = ?", doctor.ID).Association("ChatwithUser").Append(&user)
	if err != nil {
		return model.ChatRoom{}, err
	}

	ChatRoom := model.ChatRoom{}
	result := config.DB.Model(&ChatRoom).Create(model.ChatRoom{
		UserID:   user.ID,
		DoctorID: doctor.ID,
	})

	if result.RowsAffected < 1 {
		return model.ChatRoom{}, nil
	}
	var Chatroom model.ChatRoom
	config.DB.Model(model.ChatRoom{}).Where("user_id = ? AND doctor_id = ?", user.ID, doctor.ID).Find(&Chatroom)
	return Chatroom, nil
}

func saveMessage(message Message, chatroom model.ChatRoom, roles string) bool {
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
	var chat model.Chat

	if role == "doctor" {
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
