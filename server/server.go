package server

import (
	"net"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/f4814n/piston/protocol"
	"runtime"
)

type Server struct {
}

func (s *Server) Serve() {
	log.WithFields(log.Fields{
		"cpus": runtime.NumCPU(),
	}).Info("Starting Piston")

	ln, err := net.Listen("tcp", "192.168.122.107:25565")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Error(err)
		}

		c, err := protocol.NewConnection(conn.(*net.TCPConn))
		if err != nil {
			log.Error(err)
			continue
		}

		go s.waitForPackets(c)
	}
}

func (s *Server) waitForPackets(c *protocol.Connection) {
	c.GetLogger().Debug("Connection initiated")
	defer c.Close()

	for {
		p, err := c.ReadPacket();
		if err != nil {
			if err.Error() != "EOF" {
				c.GetLogger().Error(err)
			}
			return
		}
		s.handlePacket(c, p)
	}
}

func (s *Server) handlePacket(c *protocol.Connection, p protocol.Packet) {
	switch c.State {
	case protocol.Handshaking:
		handleHandshakingPacket(c, p)
	case protocol.Status:
		handleStatusPacket(c, p)
	case protocol.Login:
		handleLoginPacket(c, p)
	case protocol.Play:
		handlePlayPacket(c, p)
	}
}

func handleHandshakingPacket(c *protocol.Connection, p protocol.Packet) {
	v := p.(protocol.Handshake) // No other packets in this state
	next := v.NextState

	if next == 1 {
		c.GetLogger().WithFields(log.Fields{
			"state": protocol.Status,
		}).Debug("Switching Connection state")
		c.State = protocol.Status
	} else {
		c.GetLogger().WithFields(log.Fields{
			"state": protocol.Login,
		}).Debug("Switching Connection state")
		c.State = protocol.Login
	}
}

func handleStatusPacket(c *protocol.Connection, p protocol.Packet) {
	switch p.(type) {
	case protocol.Request:
		status := protocol.ResponseJSON{}
		status.Version.Name = "lool"
		status.Version.Protocol = 242
		status.Description.Text = "hi"
		b, _ := json.Marshal(status)

		err := c.WritePacket(protocol.Response{string(b)})

		if err != nil {
			c.GetLogger().Error(err)
		}
	}
}

func handleLoginPacket(c *protocol.Connection, p protocol.Packet) {
	switch v := p.(type) {
	case protocol.LoginStart:
		c.WritePacket(protocol.LoginSuccess{UUID: "123e4567-e89b-12d3-a456-426655440000", Username: v.Name})
		c.GetLogger().WithFields(log.Fields{
			"state": protocol.Play,
		}).Debug("Switching Connection state")
		c.State = protocol.Play
		err := c.WritePacket(protocol.JoinGame{
			EID: 6,
			Gamemode: 0,
			Dimension: 0,
			LevelType: "default",
			ViewDistance: 10,
			ReducedDebugInfo: false,
		})

		if err != nil {
			c.GetLogger().Error(err)
		}
			
	}
}

func handlePlayPacket(c *protocol.Connection, p protocol.Packet) {
}
