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
	log4 "github.com/alecthomas/log4go"
	"github.com/ontio/ontology-tool/common"
	"github.com/ontio/ontology-tool/config"
	"github.com/ontio/ontology-tool/p2pserver/net/netserver"
	"github.com/ontio/ontology-tool/p2pserver/net/protocol"
	"github.com/ontio/ontology-tool/p2pserver/protocols"
	"github.com/ontio/ontology-tool/utils/timer"
)

var (
	ns *netserver.NetServer
	tr *timer.Timer
)

func Demo() bool {
	log4.Info("hello, dht demo")
	return true
}

func setup(protocol p2p.Protocol) {
	var err error

	if ns, err = netserver.NewNetServer(protocol, config.DefConfig.Net); err != nil {
		log4.Crashf("[NewNetServer] crashed, err %s", err)
	}
	if err = ns.Start(); err != nil {
		log4.Crashf("start netserver failed, err %s", err)
	}

	tr = timer.NewTimer(2)
}

func Handshake() bool {

	// 1. get params from json file
	var params struct {
		Remote        string
		HeartbeatTime int
	}
	if err := getParamsFromJsonFile("./params/Handshake.json", &params); err != nil {
		log4.Error("%s", err)
		return false
	}

	// 2. set common params
	common.SetHandshakeDuraion(10)
	common.SetHandshakeLevel(common.HandshakeNormal)
	common.SetHeartbeatBlockHeight(358)

	// 3. setup p2p.protocols
	protocol := protocols.NewOnlyHeartbeatMsgHandler()
	setup(protocol)

	// 4. connect and handshake
	if err := ns.Connect(params.Remote); err != nil {
		log4.Debug("connecting to %s failed, err: %s", params.Remote, err)
		return false
	}

	// 5. dispatch
	dispatch(params.HeartbeatTime)
	log4.Info("handshake end!")

	return true
}
