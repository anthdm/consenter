# Consenter
A pluggable blockchain consensus simulation framework written in Go.

### What?
Testing and benchmarking complex consensus algorithms in existing codebases can sometimes be a pain. Consenter is solving this problem by exposing a simple blockchain behind a p2p network where consensus engines can easily be plugged in and out. 

### How?
The server is emiting transactions each N seconds that will be relayed through the network. Each engine implementation will receive those transactions, hence engines can than operate on those transactions according to their implementation. Messages can be quickly bootstrapped and implemented with `protobuf`.

Installing required dependencies:
```
make deps
```

Start a simulation:
```
make simulation
```

```
Recreating consenter_node_four_1 ...
Recreating consenter_node_one_1 ...
Recreating consenter_node_three_1 ...
Recreating consenter_node_five_1 ...
Recreating consenter_node_two_1 ... done
Attaching to consenter_node_three_1, consenter_node_four_1, consenter_node_one_1, consenter_node_five_1, consenter_node_two_1
node_three_1  | time="2018-05-08T11:47:12Z" level=info msg="starting p2p server.."
node_three_1  | time="2018-05-08T11:47:12Z" level=info msg="server.tcp accepting new connections on 0.0.0.0:3002"
node_three_1  | time="2018-05-08T11:47:14Z" level=info msg="new peer connected" endpoint="172.21.0.6:58050"
node_four_1   | time="2018-05-08T11:47:12Z" level=info msg="starting p2p server.."
node_four_1   | time="2018-05-08T11:47:12Z" level=info msg="server.tcp accepting new connections on 0.0.0.0:3003"
node_four_1   | time="2018-05-08T11:47:14Z" level=info msg="new peer connected" endpoint="172.21.0.6:38460"
node_one_1    | time="2018-05-08T11:47:13Z" level=info msg="starting p2p server.."
node_two_1    | time="2018-05-08T11:47:14Z" level=info msg="starting p2p server.."
node_five_1   | time="2018-05-08T11:47:13Z" level=info msg="starting p2p server.."
node_two_1    | time="2018-05-08T11:47:14Z" level=info msg="server.tcp accepting new connections on 0.0.0.0:3001"
node_one_1    | time="2018-05-08T11:47:13Z" level=info msg="server.tcp accepting new connections on 0.0.0.0:3000"
node_two_1    | time="2018-05-08T11:47:14Z" level=info msg="new peer connected" endpoint="172.21.0.2:3002"
node_five_1   | time="2018-05-08T11:47:13Z" level=info msg="server.tcp accepting new connections on 0.0.0.0:3004"
node_two_1    | time="2018-05-08T11:47:14Z" level=info msg="new peer connected" endpoint="172.21.0.4:3003"
node_five_1   | time="2018-05-08T11:47:14Z" level=info msg="new peer connected" endpoint="172.21.0.6:50642"
node_two_1    | time="2018-05-08T11:47:14Z" level=info msg="new peer connected" endpoint="172.21.0.5:3004"
node_two_1    | time="2018-05-08T11:47:15Z" level=info msg="receiving new tx: 82ae127dcd6557af38f0ca946c9e65e20a2fa682ce0a5ba3a0fabe3d7677191a"
node_three_1  | time="2018-05-08T11:47:15Z" level=info msg="receiving new tx: 82ae127dcd6557af38f0ca946c9e65e20a2fa682ce0a5ba3a0fabe3d7677191a"
node_five_1   | time="2018-05-08T11:47:15Z" level=info msg="receiving new tx: 82ae127dcd6557af38f0ca946c9e65e20a2fa682ce0a5ba3a0fabe3d7677191a"
node_two_1    | time="2018-05-08T11:47:15Z" level=info msg="receiving new tx: d9734e16b52b93254fb1aa190e6739dd802208384ff711da68f9f3a63959d585"
```

### Example
There is a [solo engine example](https://github.com/anthdm/consenter/blob/master/pkg/consensus/solo/engine.go) that should cover the idea and get you up to speed. 

### Todo
- configuration
- blockchain persistance 
- implementing engines
