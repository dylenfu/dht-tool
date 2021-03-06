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
package heatbeat

import (
	"sync/atomic"
	"time"

	log4 "github.com/alecthomas/log4go"
	tcm "github.com/ontio/ontology-tool/common"
	"github.com/ontio/ontology-tool/p2pserver/common"
	"github.com/ontio/ontology-tool/p2pserver/message/msg_pack"
	"github.com/ontio/ontology-tool/p2pserver/message/types"
	"github.com/ontio/ontology-tool/p2pserver/net/protocol"
	"github.com/ontio/ontology/common/config"
)

type HeartBeat struct {
	net    p2p.P2P
	id     common.PeerId
	quit   chan bool
	height uint64
	start  int64
}

func NewHeartBeat(net p2p.P2P) *HeartBeat {
	return &HeartBeat{
		id:     net.GetID(),
		net:    net,
		quit:   make(chan bool),
		height: tcm.HeartbeatBlockHeight,
		start:  time.Now().Unix(),
	}
}

func (self *HeartBeat) Start() {
	go self.heartBeatService()
}

func (self *HeartBeat) Stop() {
	close(self.quit)
}

func (this *HeartBeat) heartBeatService() {
	var periodTime uint = config.DEFAULT_GEN_BLOCK_TIME / common.UPDATE_RATE_PER_BLOCK
	t := time.NewTicker(time.Second * (time.Duration(periodTime)))

	for {
		select {
		case <-t.C:
			this.ping()
			this.timeout()
		case <-this.quit:
			t.Stop()
			return
		}
	}
}

// mark:
func (this *HeartBeat) ping() {
	if this.NeedInterrupt(true) {
		log4.Debug("[p2p]interrupt ping...")
		return
	}

	height := this.height
	ping := msgpack.NewPingMsg(uint64(height))
	go this.net.Broadcast(ping)
	log4.Debug("[p2p]send ping msg height %d", height)
}

//timeout trace whether some peer be long time no response
func (this *HeartBeat) timeout() {
	peers := this.net.GetNeighbors()
	var periodTime uint = config.DEFAULT_GEN_BLOCK_TIME / common.UPDATE_RATE_PER_BLOCK
	for _, p := range peers {
		t := p.GetContactTime()
		if t.Before(time.Now().Add(-1 * time.Second *
			time.Duration(periodTime) * common.KEEPALIVE_TIMEOUT)) {
			log4.Warn("[p2p]keep alive timeout!!!lost remote peer %d - %s from %s", p.GetID(), p.Link.GetAddr(), t.String())
			p.Close()
		}
	}
}

func (this *HeartBeat) PingHandle(ctx *p2p.Context, ping *types.Ping) {
	// mark:
	if this.NeedInterrupt(false) {
		log4.Debug("[p2p]interrupt pong...")
		return
	}

	remotePeer := ctx.Sender()
	remotePeer.SetHeight(ping.Height)
	p2p := ctx.Network()

	height := this.height
	p2p.SetHeight(uint64(height))
	msg := msgpack.NewPongMsg(uint64(height))

	err := remotePeer.Send(msg)
	if err != nil {
		log4.Warn(err)
	}
	log4.Debug("[p2p]recv pong msg height %d", height)
}

func (this *HeartBeat) PongHandle(ctx *p2p.Context, pong *types.Pong) {
	remotePeer := ctx.Network()
	remotePeer.SetHeight(pong.Height)
	// mark:
	atomic.AddUint64(&this.height, 1)
}

func (this *HeartBeat) NeedInterrupt(iscli bool) bool {
	lastTime := tcm.HeartbeatInterruptPingLastTime
	if !iscli {
		lastTime = tcm.HeartbeatInterruptPongLastTime
	}

	breakAfterStart := tcm.HeartbeatInterruptAfterStartTime
	now := time.Now().Unix()
	start := this.start + breakAfterStart
	end := this.start + breakAfterStart + lastTime

	if breakAfterStart > 0 && now > start && now < end {
		return true
	}
	return false
}
