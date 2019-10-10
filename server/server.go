package server

import (
	"encoding/json"
	"github.com/f4814n/piston/protocol"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"runtime"
	"runtime/debug"
)

type Server struct {
}

func (s *Server) Serve(address string) {
	buildInfo, _ := debug.ReadBuildInfo()
	log.WithFields(log.Fields{
		"cpus":    runtime.NumCPU(),
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"go":      runtime.Version(),
		"version": buildInfo.Main.Version,
	}).Info("Starting Piston")

	log.WithFields(log.Fields{
		"address": address,
	}).Info("Listening for connections")
	ln, err := protocol.Listen(address)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Error(err)
		}

		if err != nil {
			log.Error(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn *protocol.Conn) {
	defer conn.Close()

	p, err := conn.ReadPacket()
	if err != nil {
		conn.GetLogger().Error(err)
		return
	}

	v := p.(protocol.Handshake) // No other packets in this state
	if v.NextState == 1 {       // Client only attempts SLP
		conn.GetLogger().WithFields(log.Fields{
			"state": protocol.Status,
		}).Trace("Switching Connection state")
		conn.State = protocol.Status
		s.status(conn)
		return // No further communication
	}

	// Client wants to login
	conn.GetLogger().WithFields(log.Fields{
		"state": protocol.Login,
	}).Trace("Switching Connection state")
	conn.State = protocol.Login
	player, err := s.login(conn)
	if err != nil {
		conn.GetLogger().Error(err)
		return
	}

	conn.GetLogger().WithFields(log.Fields{
		"state": protocol.Play,
	}).Trace("Switching connection state")
	conn.State = protocol.Play
	player.Start() // Blocks until player disconnects
}

// Handle Server list ping
func (s *Server) status(conn *protocol.Conn) {
	p, err := conn.ReadPacket()
	if err != nil {
		conn.GetLogger().Error(err)
		return
	}

	// protocol.Request has no fields. So we cannot assign a variable
	if _, ok := p.(protocol.Request); !ok {
		conn.GetLogger().Error(err)
		return
	}
	status := protocol.ResponseJSON{}
	status.Version.Name = "piston"
	status.Version.Protocol = 498
	status.Description.Text = "hi"
	b, _ := json.Marshal(status)
	err = conn.WritePacket(&protocol.Response{JSONResponse: string(b)})
	if err != nil {
		conn.GetLogger().Error(err)
		return
	}

	p, err = conn.ReadPacket()
	if err != nil {
		conn.GetLogger().Error(err)
		return
	}

	v := p.(protocol.Ping)
	err = conn.WritePacket(protocol.Pong(v))
	if err != nil {
		conn.GetLogger().Error(err)
	}
}

// Authenticate a client
func (s *Server) login(conn *protocol.Conn) (Player, error) {
	p, err := conn.ReadPacket()
	if err != nil {
		return Player{}, err
	}

	v := p.(protocol.LoginStart)
	namespace, _ := uuid.Parse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
	uuid := uuid.NewSHA1(namespace, []byte(v.Name))
	err = conn.WritePacket(protocol.LoginSuccess{UUID: uuid.String(), Username: v.Name})
	if err != nil {
		return Player{}, err
	}

	player := NewPlayer(v.Name, uuid, conn)

	return player, nil
}
