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
 * You should contains received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
package recent_peers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	log4 "github.com/alecthomas/log4go"
	"github.com/ontio/ontology-tool/p2pserver/common"
	p2p "github.com/ontio/ontology-tool/p2pserver/net/protocol"
	common2 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/config"
)

type PersistRecentPeerService struct {
	net         p2p.P2P
	quit        chan bool
	recentPeers map[uint32][]*RecentPeer
	lock        sync.RWMutex
}

func (this *PersistRecentPeerService) contains(addr string) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	netID := config.DefConfig.P2PNode.NetworkMagic
	for i := 0; i < len(this.recentPeers[netID]); i++ {
		if this.recentPeers[netID][i].Addr == addr {
			return true
		}
	}
	return false
}

func (this *PersistRecentPeerService) AddNodeAddr(addr string) {
	if this.contains(addr) {
		return
	}
	this.lock.Lock()
	netID := config.DefConfig.P2PNode.NetworkMagic
	this.recentPeers[netID] = append(this.recentPeers[netID],
		&RecentPeer{
			Addr:  addr,
			Birth: time.Now().Unix(),
		})
	this.lock.Unlock()
}

func (this *PersistRecentPeerService) DelNodeAddr(addr string) {
	this.lock.Lock()
	netID := config.DefConfig.P2PNode.NetworkMagic
	for i := 0; i < len(this.recentPeers[netID]); i++ {
		if this.recentPeers[netID][i].Addr == addr {
			this.recentPeers[netID] = append(this.recentPeers[netID][:i], this.recentPeers[netID][i+1:]...)
		}
	}
	this.lock.Unlock()
}

type RecentPeer struct {
	Addr  string
	Birth int64
}

func (this *PersistRecentPeerService) saveToFile() {
	temp := make(map[uint32][]string)
	for networkId, rps := range this.recentPeers {
		temp[networkId] = make([]string, 0)
		for _, rp := range rps {
			elapse := time.Now().Unix() - rp.Birth
			if elapse > common.RecentPeerElapseLimit {
				temp[networkId] = append(temp[networkId], rp.Addr)
			}
		}
	}
	buf, err := json.Marshal(temp)
	if err != nil {
		log4.Warn("[p2p]package recent peer fail: ", err)
		return
	}
	err = ioutil.WriteFile(common.RECENT_FILE_NAME, buf, os.ModePerm)
	if err != nil {
		log4.Warn("[p2p]write recent peer fail: ", err)
	}
}

func NewPersistRecentPeerService(net p2p.P2P) *PersistRecentPeerService {
	return &PersistRecentPeerService{
		net:  net,
		quit: make(chan bool),
	}
}

func (self *PersistRecentPeerService) Stop() {
	close(self.quit)
}

func (this *PersistRecentPeerService) loadRecentPeers() {
	this.recentPeers = make(map[uint32][]*RecentPeer)
	if common2.FileExisted(common.RECENT_FILE_NAME) {
		buf, err := ioutil.ReadFile(common.RECENT_FILE_NAME)
		if err != nil {
			log4.Warn("[p2p]read %s fail:%s, connect recent peers cancel", common.RECENT_FILE_NAME, err.Error())
			return
		}

		temp := make(map[uint32][]string)
		err = json.Unmarshal(buf, &temp)
		if err != nil {
			log4.Warn("[p2p]parse recent peer file fail: ", err)
			return
		}
		for networkId, addrs := range temp {
			for _, addr := range addrs {
				this.recentPeers[networkId] = append(this.recentPeers[networkId], &RecentPeer{
					Addr:  addr,
					Birth: time.Now().Unix(),
				})
			}
		}
	}
}

func (this *PersistRecentPeerService) Start() {
	this.loadRecentPeers()
	this.tryRecentPeers()
	go this.syncUpRecentPeers()
}

//tryRecentPeers try connect recent contact peer when service start
func (this *PersistRecentPeerService) tryRecentPeers() {
	netID := config.DefConfig.P2PNode.NetworkMagic
	if len(this.recentPeers[netID]) > 0 {
		log4.Info("[p2p] try to connect recent peer")
	}
	for _, v := range this.recentPeers[netID] {
		go this.net.Connect(v.Addr)
	}
}

//syncUpRecentPeers sync up recent peers periodically
func (this *PersistRecentPeerService) syncUpRecentPeers() {
	periodTime := common.RECENT_TIMEOUT
	t := time.NewTicker(time.Second * (time.Duration(periodTime)))
	for {
		select {
		case <-t.C:
			this.saveToFile()
		case <-this.quit:
			t.Stop()
			return
		}
	}
}
