# dth-tool 
ontology p2pserver 测试工具

## 使用方式
下载项目ontology, 并切换到net-review分支: 
```bash
https://github.com/dylenfu/ontology.git
```
下载项目dht-tool
```bash
https://github.com/dylenfu/dht-tool.git
```
根据环境make得到可执行文件

在config.json文件配置节点列表以及magic等参数。
具体测试时，不同的测试用例在params路径下修改用例参数。
然后运行命令如下:
```bash
./dht-tool -t=handshake
```
也支持批量测试
```bash
./dht-tool -t=handshake,heartbeat
```

## 测试用例
```dtd
handshake                           // 握手
handshakeTimeout                    // 握手超时测试
handshakeWrongMsg                   // 握手客户端发送错误信息
heartbeat                           // 心跳持续测试
heartbeatInterruptPing              // p2p ping中断测试
heartbeatInterruptPong              // p2p pong中断测试
resetPeerID                         // 重置peerID
ddos                                // ddos 攻击单一节点，阻断流量(无法同步块)
invalidBlockHeight                  // 模拟节点持续快高异常 
attackRoutable                      // 路由表攻击
attackTxPool                        // 交易池攻击
doubleSpend                         // 双花攻击
```

## 测试条件及结果预期
#### 1、握手测试
```dtd
条件:
a、正常握手，或者在握手时停止于某个步骤
结果:
a、正常握手连接应该成功
b、握手中断连接应该失败
```

#### 2、握手时发送非法version
```dtd
条件:
a、使用参数构造虚假version，并发送到某个目标节点
结果:
a、连接失败
```

#### 3、网络超时重试模拟
```dtd
条件:
a、握手时在某个步骤延时
结果:
a、第一次握手失败
b、第二次握手成功
```

#### 4、心跳测试
```dtd
条件:
a、保持正常心跳
结果:
a、连接正常，模拟块高持续增加
```

#### 5、心跳中断ping
```dtd
条件:
a、心跳过程中，主动中断ping，持续n sec
结果:
a、连接正常，块高保持一定高度后持续增长
b、连接断开
```

#### 6、心跳中断pong
```dtd
条件:
a、心跳过程中，主动中断pong，持续n sec
结果:
a、连接正常，块高保持一定高度后持续增长
b、连接断开
```

#### 7、更新peerID
```dtd
条件:
a、建立连接保持心跳后，变更peerID
结果:
a、连接断开
```

#### 8、网络流量攻击
```dtd
条件:
a、构造多个虚假peerID
b、与单个目标sync节点距离较近
c、虚假peer主动发起连接，并持续ping
结果:
a、节点不能正常出块(出块慢甚至不出块)
```

#### 9.路由表攻击
```dtd
条件:
1.根据目标节点ID，构造大量距离目标seed节点很近的虚拟节点，
2.主动连接并ping目标节点，比如说节点最多允许接收1024个链接，其中4个链接为正常连接，
    1020个链接为恶意连接，当连接数超过1024时，是否会有恶意连接挤出正常连接
结果:
a、正常连接是否会被挤出
b、节点路由被恶意节点占满
c、节点重启后仍然被恶意节点占满
d、节点主动发起连接大概率会连接上恶意节点
```

#### 10.块同步攻击
```dtd
条件:
a.测试工具所模拟的节点块高度始终高于正常节点
结果:
a.造成同步异常或者延时
```

#### 11.交易池攻击
```dtd
条件:
a、多个恶意节点持续对多个目标seed节点发送不合法交易(比如余额不足)
结果:
a、查询交易池，不同目标节点的交易持应该相同，而且都不包含不合法交易
b、测试前后查询余额，账户余额不变
```

#### 12.双花攻击
```dtd
条件:
a、单个恶意节点，对多个目标seed节点发送连续的4笔交易，其中1笔能成功，另外3笔不能成功，
   比如只有2块钱的情况下，转账4次，1.1， 1.2， 1.3，1.4
结果:
a、目标seed节点交易池应该相同，都只有1笔正常的交易
b、测试前后查询余额账户，只转出一笔
```

## 测试参数:
#### 1.握手测试参数
 params/Handshake.json

```dtd
HandshakeNormal = 0                 // 正常握手
StopClientAfterSendVersion = 1      // 握手时客户端发送version后停止
StopClientAfterReceiveVersion = 2   // 握手时客户端接收version后停止
StopClientAfterUpdateKad = 3        // 握手时客户端更新kad后停止
StopClientAfterReadKad = 4          // 握手时客户端读取kad后停止
StopClientAfterSendAck = 5          // 握手时客户端发送ack后停止
StopServerAfterSendVersion = 6      // 握手时服务端发送version后停止
StopServerAfterReceiveVersion = 7   // 握手时服务端结束到version后停止
StopServerAfterUpdateKad = 8        // 握手时服务端更新kad后停止
StopServerAfterReadKad = 9          // 握手时服务端读取kad后停止
StopServerAfterReadAck = 10         // 握手时服务端接收ack后停止
```

## TODO
. dump dht
. dump txnPool
. utils for sending tx
. utils for getting block height & account balance