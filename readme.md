# dth-tool ontology p2pserver 测试工具

## 测试用例
```dtd
handshake               // 握手
handshakeTimeout        // 握手超时测试
handshakeWrongMsg       // 握手客户端发送错误信息
heartbeat               // 心跳持续测试
heartbeatBreak          // 心跳断开测试
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
StopServerAfterReadKad = 9          // 握手时服务端更新kad后停止
StopServerAfterReadAck = 10         // 握手时服务端接收ack后停止
```
