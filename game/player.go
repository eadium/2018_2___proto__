package game

import (
	ws "github.com/gorilla/websocket"
	m "proto-game-server/models"
)

type Player struct {
	session *m.Session
	conn    *ws.Conn
}

func NewPlayer(session *m.Session, conn *ws.Conn) *Player {
	return &Player{session, conn}
}

func (p *Player) ReadCommand() (*Command, error){
	command := &Command{}
	err := p.conn.ReadJSON(command)

	return command, err
}