package logic

import (
	"fmt"
	"qnmahjong/def"
	"qnmahjong/msg"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/olahol/melody"
)

var (
	e *echo.Echo
	m *melody.Melody
)

// Start logic server
func Start() {
	// Echo instance
	e = echo.New()
	// Melody instance
	m = melody.New()

	m.Config.MaxMessageSize = 1024 * 1024
	m.Config.MessageBufferSize = 1024 * 1024

	// change allow all origin hosts for 403 error
	m.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// Debug mode
	// e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Logic Websocket Routes => Handler
	e.GET("/", func(c echo.Context) error {
		return m.HandleRequest(c.Response(), c.Request())
	})

	// HandleConnect fires fn when a session connects.
	m.HandleConnect(func(s *melody.Session) {
		log.WithFields(log.Fields{}).Info("ws connect")
	})

	// HandleDisconnect fires fn when a session disconnects.
	m.HandleDisconnect(func(s *melody.Session) {
		playerID, _ := s.Get("PlayerID")
		log.WithFields(log.Fields{
			"playerID": playerID,
		}).Info("ws disconnect")
		s.Close()
	})

	// HandlePong fires fn when a pong is received from a session.
	m.HandlePong(func(s *melody.Session) {
		// playerID, _ := s.Get("PlayerID")
		// log.WithFields(log.Fields{
		// 	"playerID": playerID,
		// }).Info("ws pong")
	})

	// HandleMessage fires fn when a text message comes in.
	m.HandleMessage(func(s *melody.Session, b []byte) {
		// playerID, _ := s.Get("PlayerID")
		// log.WithFields(log.Fields{
		// 	"playerID": playerID,
		// 	"recv":     string(b),
		// }).Info("ws recv text")
	})

	// HandleMessageBinary fires fn when a binary message comes in.
	m.HandleMessageBinary(func(s *melody.Session, b []byte) {
		// playerID, _ := s.Get("PlayerID")
		// log.WithFields(log.Fields{
		// 	"playerID": playerID,
		// 	"recv":     string(b),
		// }).Info("ws recv binary")

		if string(b) == "ping" {
			s.WriteBinary([]byte("pong"))
			return
		}

		err := msg.LogicHandle(b, s)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error(def.ErrHandleLogic)
			return
		}
	})

	// HandleSentMessage fires fn when a text message is successfully sent.
	m.HandleSentMessage(func(s *melody.Session, b []byte) {
		// playerID, _ := s.Get("PlayerID")
		// log.WithFields(log.Fields{
		// 	"playerID": playerID,
		// 	"send":     string(b),
		// }).Info("ws send text")
	})

	// HandleSentMessageBinary fires fn when a binary message is successfully sent.
	m.HandleSentMessageBinary(func(s *melody.Session, b []byte) {
		// playerID, _ := s.Get("PlayerID")
		// log.WithFields(log.Fields{
		// 	"playerID": playerID,
		// 	"send":     string(b),
		// }).Info("ws send binary")
	})

	// HandleError fires fn when a session has an error.
	m.HandleError(func(s *melody.Session, err error) {
		playerID, _ := s.Get("PlayerID")
		log.WithFields(log.Fields{
			"playerID": playerID,
			"error":    err,
		}).Info("ws error")
		s.Close()
	})

	// Start server
	err := e.Start(":5002")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrStartLogic)
		fmt.Println("start logic failed")
		os.Exit(-1)
	}
}

// Shutdown logic server
func Shutdown() {
	if m != nil {
		err := m.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error(def.ErrShutdownLogic)
		}
		fmt.Println("logic shut down")
	}
}
