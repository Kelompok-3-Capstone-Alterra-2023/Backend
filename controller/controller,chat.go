package controller

import (
	"capstone/middleware"
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

func ConnectWSUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(middleware.ExtractUserIdToken(token))

	currentconn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), c.Response().Header())

	if err != nil {
		return c.JSON(http.StatusBadRequest, "failure when connect websocket")
	}

	wsconnuser[userID] = currentconn

	go handleIO(currentconn, wsconndoctor, userID)
	return c.JSON(http.StatusAccepted, "success create connection")
}

func ConnectWSDoctor(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed when cast jwt",
			"error":   err,
		})
	}

	currentconn, err2 := upgrader.Upgrade(c.Response().Writer, c.Request(), c.Response().Header())

	if err2 != nil {
		return c.JSON(http.StatusBadRequest, "failure when connect websocket")
	}

	wsconndoctor[int(doctorID)] = currentconn

	go handleIO(currentconn, wsconnuser, int(doctorID))
	return c.JSON(http.StatusAccepted, "success create connection")
}

func handleIO(currentconn *websocket.Conn, connectionmapsender map[int]*websocket.Conn, from int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("error", fmt.Sprintf("%v", r))
		}
	}()

	for {
		message := Message{}

		err := currentconn.ReadJSON(&message)
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

		sendMessage(currentconn, message, connectionmapsender[message.To])
	}
}

func sendMessage(currentconn *websocket.Conn, message Message, destination *websocket.Conn) {
	destination.WriteJSON(message)
}

func closeconn(currentconn *websocket.Conn, message Message) {
	filtered := gubrak.From(wsconnuser).Reject(func(each *websocket.Conn) bool {
		return each == currentconn
	}).Result()
	wsconnuser[message.From] = filtered.(*websocket.Conn)
}
