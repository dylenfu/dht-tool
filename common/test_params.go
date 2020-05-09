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

package common

import "time"

// handshake test cases
const (
	HandshakeNormal = iota
	Handshake_StopAfterSendVersion
	Handshake_StopAfterReceiveVersion
	Handshake_StopAfterUpdateKad
	Handshake_StopAfterReadKad
	Handshake_StopAfterSendAck
	Handshake_StopAfterReadAck
)

var (
	HandshakeLevel    uint8         = HandshakeNormal                 // default normal
	HandshakeDuration time.Duration = time.Duration(10) * time.Second // default value: 10 sec
)

func SetHandshakeLevel(lvl uint8) {
	HandshakeLevel = lvl
}
func StopHandshake(lvl uint8) bool {
	return HandshakeLevel == lvl
}
func SetHandshakeDuraion(sec int) {
	HandshakeDuration = time.Duration(sec) * time.Second
}

// heartbeat test cases
var HeartbeatBlockHeight uint64 = 358 // default 100000
func SetHeartbeatBlockHeight(height uint64) {
	HeartbeatBlockHeight = height
}
