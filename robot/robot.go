package robot

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"qnmahjong/pf"
	"net/http"
	"net/url"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

// Start robot server
func Start() {

	// login

	loginURL := "http://localhost:5001/"
	contentType := "application/x-www-form-urlencoded"

	loginSend := pf.LoginSend{
		LoginType: 1,
		MachineID: "test1",
	}
	msg, err := proto.Marshal(&loginSend)
	if err != nil {
		return
	}

	absMessage := pf.AbsMessage{
		MsgID:   int32(pf.Login),
		MsgBody: msg,
	}
	msg, err = proto.Marshal(&absMessage)
	if err != nil {
		return
	}

	resp, err := http.Post(loginURL, contentType, bytes.NewBuffer(msg))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	buf, err = base64.StdEncoding.DecodeString(string(buf))
	absMessage = pf.AbsMessage{}
	err = proto.Unmarshal(buf, &absMessage)
	if err != nil {
		return
	}

	token := absMessage.GetToken()
	loginRecv := pf.LoginRecv{}
	err = proto.Unmarshal(absMessage.GetMsgBody(), &loginRecv)
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", loginRecv)

	// auth

	authSend := pf.AuthSend{}
	msg, err = proto.Marshal(&authSend)
	if err != nil {
		return
	}

	absMessage = pf.AbsMessage{
		Token:   token,
		MsgID:   int32(pf.Auth),
		MsgBody: msg,
	}
	msg, err = proto.Marshal(&absMessage)
	if err != nil {
		return
	}

	u := url.URL{Scheme: "ws", Host: "localhost:5002", Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	defer c.Close()

	err = c.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		return
	}

	_, buf, err = c.ReadMessage()
	absMessage = pf.AbsMessage{}
	err = proto.Unmarshal(buf, &absMessage)
	if err != nil {
		return
	}

	authRecv := pf.AuthRecv{}
	err = proto.Unmarshal(absMessage.GetMsgBody(), &authRecv)
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", authRecv)

	// dirty

	var dirtySend = pf.DirtySend{}
	msg, err = proto.Marshal(&dirtySend)
	if err != nil {
		return
	}

	absMessage = pf.AbsMessage{
		Token:   token,
		MsgID:   int32(pf.Dirty),
		MsgBody: msg,
	}
	msg, err = proto.Marshal(&absMessage)
	if err != nil {
		return
	}

	err = c.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		return
	}

	_, buf, err = c.ReadMessage()
	absMessage = pf.AbsMessage{}
	err = proto.Unmarshal(buf, &absMessage)
	if err != nil {
		return
	}

	dirtyRecv := pf.DirtyRecv{}
	err = proto.Unmarshal(absMessage.GetMsgBody(), &dirtyRecv)
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", dirtyRecv)

}

// Shutdown robot server
func Shutdown() {

}

// ConnectLogin connect login http server
func ConnectLogin() {

}

// ConnectLogic connect logic ws server
func ConnectLogic() {

}
