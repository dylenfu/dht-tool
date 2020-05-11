# dth-tool ontology p2pserver 测试工具

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
