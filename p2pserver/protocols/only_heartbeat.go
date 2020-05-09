/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package protocols

import (
	log4 "github.com/alecthomas/log4go"
	msgCommon "github.com/ontio/ontology-tool/p2pserver/common"
	msgTypes "github.com/ontio/ontology-tool/p2pserver/message/types"
	"github.com/ontio/ontology-tool/p2pserver/net/protocol"
	"github.com/ontio/ontology-tool/p2pserver/protocols/heatbeat"
)

type OnlyHeartbeatMsgHandler struct {
	heatBeat *heatbeat.HeartBeat
}

func NewOnlyHeartbeatMsgHandler() *OnlyHeartbeatMsgHandler {
	return &OnlyHeartbeatMsgHandler{}
}

func (self *OnlyHeartbeatMsgHandler) start(net p2p.P2P) {
	self.heatBeat = heatbeat.NewHeartBeat(net)
	go self.heatBeat.Start()
}

func (self *OnlyHeartbeatMsgHandler) stop() {
	self.heatBeat.Stop()
}

func (self *OnlyHeartbeatMsgHandler) HandleSystemMessage(net p2p.P2P, msg p2p.SystemMessage) {
	switch m := msg.(type) {
	case p2p.NetworkStart:
		self.start(net)
	case p2p.PeerConnected:
		log4.Debug("peer connected, address: %s, id %d", m.Info.Addr, m.Info.Id.ToUint64())
	case p2p.PeerDisConnected:
		log4.Debug("peer disconnected, address: %s, id %d", m.Info.Addr, m.Info.Id.ToUint64())
	case p2p.NetworkStop:
		self.stop()
	}
}

func (self *OnlyHeartbeatMsgHandler) HandlePeerMessage(ctx *p2p.Context, msg msgTypes.Message) {
	log4.Trace("[p2p]receive message, remote address %s, id %d", ctx.Sender().GetAddr(), ctx.Sender().GetID().ToUint64())
	switch m := msg.(type) {
	case *msgTypes.Ping:
		self.heatBeat.PingHandle(ctx, m)
	case *msgTypes.Pong:
		self.heatBeat.PongHandle(ctx, m)
	case *msgTypes.NotFound:
		log4.Debug("[p2p]receive notFound message, hash is %s", m.Hash.ToHexString())
	default:
		msgType := msg.CmdType()
		if msgType == msgCommon.VERACK_TYPE || msgType == msgCommon.VERSION_TYPE {
			log4.Info("receive message: %s from peer %s", msgType, ctx.Sender().GetAddr())
		} else {
			log4.Warn("unknown message handler for the msg: %s", msgType)
		}
	}
}
