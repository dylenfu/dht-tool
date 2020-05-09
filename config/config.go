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

package config

import (
	"encoding/json"
	"fmt"
	log4 "github.com/alecthomas/log4go"
	cmf "github.com/ontio/ontology/common/config"
	"io/ioutil"
	"os"
)

var DefConfig = NewDHTConfig()

type DHTConfig struct {
	Seed []string
	Sync []string
	Net  *cmf.P2PNodeConfig
	Sdk  *SDKConfig
}

type SDKConfig struct {
	JsonRpcAddress   string
	RestfulAddress   string
	WebSocketAddress string

	//Gas Price of transaction
	GasPrice uint64
	//Gas Limit of invoke transaction
	GasLimit uint64
	//Gas Limit of deploy transaction
	GasDeployLimit uint64
}

func NewDHTConfig() *DHTConfig {
	return &DHTConfig{}
}

func (c *DHTConfig) Init(fileName string) error {
	err := c.loadConfig(fileName)
	if err != nil {
		return fmt.Errorf("loadConfig error:%s", err)
	}
	return nil
}

func (this *DHTConfig) loadConfig(fileName string) error {
	data, err := this.readFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, this)
	if err != nil {
		return fmt.Errorf("json.Unmarshal TestConfig:%s error:%s", data, err)
	}
	cmf.DefConfig.P2PNode = this.Net
	return nil
}

func (this *DHTConfig) readFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("OpenFile %s error %s", fileName, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log4.Error("File %s close error %s", fileName, err)
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll %s error %s", fileName, err)
	}
	return data, nil
}
