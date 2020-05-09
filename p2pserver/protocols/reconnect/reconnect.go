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

package reconnect

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	log4 "github.com/alecthomas/log4go"
	"github.com/ontio/ontology-tool/p2pserver/common"
	p2p "github.com/ontio/ontology-tool/p2pserver/net/protocol"
	"github.com/ontio/ontology-tool/p2pserver/peer"
	"github.com/ontio/ontology/common/config"
)

//ReconnectService contain addr need to reconnect
type ReconnectService struct {
	sync.RWMutex
	MaxRetryCount uint
	RetryAddrs    map[string]int
	net           p2p.P2P
	quit          chan bool
}

func NewReconectService(net p2p.P2P) *ReconnectService {
	return &ReconnectService{
		net:           net,
		MaxRetryCount: common.MAX_RETRY_COUNT,
		quit:          make(chan bool),
		RetryAddrs:    make(map[string]int),
	}
}

func (self *ReconnectService) Start() {
	go self.keepOnlineService()
}

func (self *ReconnectService) Stop() {
	close(self.quit)
}

func (this *ReconnectService) keepOnlineService() {
	tick := time.NewTicker(time.Second * common.CONN_MONITOR)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			this.retryInactivePeer()
		case <-this.quit:
			return
		}
	}
}

func getPeerListenAddr(p *peer.PeerInfo) (string, error) {
	addrIp, err := common.ParseIPAddr(p.Addr)
	if err != nil {
		return "", fmt.Errorf("failed to parse addr: %s", p.Addr)
	}
	nodeAddr := addrIp + ":" + strconv.Itoa(int(p.Port))
	return nodeAddr, nil
}

func (self *ReconnectService) OnAddPeer(p *peer.PeerInfo) {
	nodeAddr, err := getPeerListenAddr(p)
	if err != nil {
		log4.Error("failed to parse addr: %s", p.Addr)
		return
	}
	self.Lock()
	delete(self.RetryAddrs, nodeAddr)
	self.Unlock()
}

func (self *ReconnectService) OnDelPeer(p *peer.PeerInfo) {
	nodeAddr, err := getPeerListenAddr(p)
	if err != nil {
		log4.Error("failed to parse addr: %s", p.Addr)
		return
	}
	self.Lock()
	self.RetryAddrs[nodeAddr] = 0
	self.Unlock()
}

func (this *ReconnectService) retryInactivePeer() {
	net := this.net
	connCount := net.GetOutConnRecordLen()
	if connCount >= config.DefConfig.P2PNode.MaxConnOutBound {
		log4.Warn("[p2p]Connect: out connections(%d) reach max limit(%d)", connCount,
			config.DefConfig.P2PNode.MaxConnOutBound)
		return
	}

	//try connect
	if len(this.RetryAddrs) > 0 {
		this.Lock()

		list := make(map[string]int)
		addrs := make([]string, 0, len(this.RetryAddrs))
		for addr, v := range this.RetryAddrs {
			v += 1
			addrs = append(addrs, addr)
			if v < common.MAX_RETRY_COUNT {
				list[addr] = v
			}
		}

		this.RetryAddrs = list
		this.Unlock()
		for _, addr := range addrs {
			rand.Seed(time.Now().UnixNano())
			log4.Debug("[p2p]Try to reconnect peer, peer addr is ", addr)
			<-time.After(time.Duration(rand.Intn(common.CONN_MAX_BACK)) * time.Millisecond)
			log4.Debug("[p2p]Back off time`s up, start connect node")
			net.Connect(addr)
		}
	}
}

func (self *ReconnectService) ReconnectCount() int {
	self.RLock()
	defer self.RUnlock()
	return len(self.RetryAddrs)
}
