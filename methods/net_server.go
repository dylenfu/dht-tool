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

package methods

import (
	"fmt"
	log4 "github.com/alecthomas/log4go"
	"github.com/ontio/ontology-tool/config"
	"github.com/ontio/ontology-tool/p2pserver/common"
	"github.com/ontio/ontology-tool/p2pserver/connect_controller"
	"github.com/ontio/ontology-tool/p2pserver/message/types"
	netsvr "github.com/ontio/ontology-tool/p2pserver/net/netserver"
	commCfg "github.com/ontio/ontology/common/config"
	// "github.com/ontio/ontology-tool/p2pserver/net/protocol"
	"github.com/ontio/ontology-tool/p2pserver/peer"
	"net"
	"time"
)

const Version = ""

type NetServer struct {
	base     *peer.PeerInfo
	listener net.Listener
	NetChan  chan *types.MsgPayload
	Np       *netsvr.NbrPeers

	connCtrl *connect_controller.ConnectController

	stopRecvCh chan bool // To stop sync channel
}

func NewNetServer() (*NetServer, error) {
	n := &NetServer{
		NetChan:    make(chan *types.MsgPayload, common.CHAN_CAPABILITY),
		base:       &peer.PeerInfo{},
		Np:         netsvr.NewNbrPeers(),
		stopRecvCh: make(chan bool),
	}

	err := n.init(config.DefConfig.Net, Version)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (s *NetServer) init(conf *commCfg.P2PNodeConfig, version string) error {
	keyId := common.RandPeerKeyId()

	httpInfo := conf.HttpInfoPort
	nodePort := conf.NodePort
	if nodePort == 0 {
		return fmt.Errorf("[p2p]invalid link port")
	}

	s.base = peer.NewPeerInfo(keyId.Id, common.PROTOCOL_VERSION, common.SERVICE_NODE, true, httpInfo,
		nodePort, 0, version, "")

	option, err := connect_controller.ConnCtrlOptionFromConfig(conf)
	if err != nil {
		return err
	}
	s.connCtrl = connect_controller.NewConnectController(s.base, keyId, option)

	syncPort := s.base.Port
	if syncPort == 0 {
		return fmt.Errorf("[p2p]sync port invalid")
	}
	s.listener, err = connect_controller.NewListener(syncPort, conf)
	if err != nil {
		return fmt.Errorf("[p2p]failed to create sync listener")
	}

	log4.Info("[p2p]init peer ID to %s", s.base.Id.ToHexString())

	return nil
}

func (this *NetServer) handleClientConnection(conn net.Conn) error {
	peerInfo, conn, err := this.connCtrl.AcceptConnect(conn)
	if err != nil {
		return err
	}
	remotePeer := createPeer(peerInfo, conn)
	remotePeer.AttachChan(this.NetChan)
	this.ReplacePeer(remotePeer)

	go remotePeer.Link.Rx()

	// todo
	// this.protocol.HandleSystemMessage(this, p2p.PeerConnected{Info: remotePeer.Info})
	return nil
}

func createPeer(info *peer.PeerInfo, conn net.Conn) *peer.Peer {
	remotePeer := peer.NewPeer()
	remotePeer.SetInfo(info)
	remotePeer.Link.UpdateRXTime(time.Now())
	remotePeer.Link.SetAddr(conn.RemoteAddr().String())
	remotePeer.Link.SetConn(conn)
	remotePeer.Link.SetID(info.Id)

	return remotePeer
}

func (this *NetServer) SendTo(p common.PeerId, msg types.Message) {
	// todo
	//peer := this.GetPeer(p)
	//if peer != nil {
	//	this.Send(peer, msg)
	//}
}

func (this *NetServer) GetPeer(id common.PeerId) *peer.Peer {
	return this.Np.GetPeer(id)
}

func (this *NetServer) ReplacePeer(remotePeer *peer.Peer) {
	// todo
	//old := this.Np.ReplacePeer(remotePeer, this)
	//if old != nil {
	//	old.Close()
	//}
}
