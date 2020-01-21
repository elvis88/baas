package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"

	"os"
	"os/exec"
	"os/signal"

	"github.com/elvis88/baas/core/ws"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("client") // name of config file
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the feconfig.yaml file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Sprintf("read config file error: %s", err))
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	id := flag.Uint("id", 0, "agent id")
	flag.Parse()
	first := true
	for {
		conn, _, err := websocket.DefaultDialer.Dial(viper.GetString("client.ws"), nil)
		if err != nil {
			fmt.Println("dial error: ", err)
			if first {
				break
			}
			select {
			case <-time.After(time.Duration(time.Duration(viper.GetInt64("client.reconnect")) * time.Second)):
			// 自动重连
			case <-interrupt:
				return
			}
			continue
		}
		fmt.Println("dial success")

		msg := ws.NewMsg(ws.ReqJoinMsg, &ws.AgentStatus{
			ID: *id,
		})
		if err := conn.WriteMessage(websocket.TextMessage, msg.Bytes()); err != nil {
			conn.Close()
			continue
		}

		first = false
		writeMsgChan := make(chan *ws.Message, 200)
		done := make(chan struct{})
		quit := make(chan struct{})
		go func() {
			for {
				defer close(done)
				if err := handleMessage(conn, writeMsgChan, quit); err != nil {
					fmt.Println("read error: ", err)
					return
				}
			}
		}()

		ticker := time.NewTicker(time.Duration(viper.GetInt64("client.period")) * time.Second)
		for {
			select {
			case <-done:
				goto exit
			case <-quit:
				goto quit
			case <-ticker.C:
				writeMsgChan <- ws.NewMsg(ws.ReqServerMsg, serverInfo())
				m.iterNode(func(node *ws.Node) {
					ps := processInfo(node.PID)
					ps.NodeID = node.ID
					writeMsgChan <- ws.NewMsg(ws.ReqProcessMsg, ps)
				})
			case msg := <-writeMsgChan:
				err := conn.WriteMessage(websocket.TextMessage, msg.Bytes())
				if err != nil {
					fmt.Println("write error: ", err)
					goto exit
				}
			case <-interrupt:
				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					fmt.Println("write close error:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	quit:
		ticker.Stop()
		conn.Close()
		return
	exit:
		ticker.Stop()
		conn.Close()
	}

}

func handleMessage(conn *websocket.Conn, writeMsgChan chan<- *ws.Message, quit chan struct{}) error {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		msg, err := ws.MsgFromBytes(message)
		if err != nil {
			return err
		}
		switch msg.Method {
		case ws.RespRefuseMsg:
			close(quit)
		case ws.RespCommandMsg:
			sc := &ws.StatusCmd{}
			err = msg.DecodePayload(sc)
			if err == nil {
				if sc.Cmd[:3] == "rpc" {

				} else {
					command := exec.Command("/bin/bash", "-c", sc.Cmd) //初始化Cmd
					res, err := command.Output()                       //运行脚本
					payload := &ws.StatusCmd{
						Value: string(res),
					}
					if err != nil {
						payload.Error = err.Error()
					}
					rmsg := ws.NewMsg(ws.ReqCommandMsg, payload)
					writeMsgChan <- rmsg
				}
			}
		case ws.RespAddNodeMsg:
			var nodes []*ws.Node
			err = msg.DecodePayload(&nodes)
			if err == nil {
				for _, node := range nodes {
					m.addNode(node)
				}
			}
		case ws.RespRemoveNodeMsg:
			var nodes []*ws.Node
			err = msg.DecodePayload(&nodes)
			if err == nil {
				for _, node := range nodes {
					m.removeNode(node)
				}
			}
		default:
			return fmt.Errorf("not support method %d", msg.Method)
		}
		fmt.Printf("recv: %v\n", string(msg.Payload))
	}
}

var (
	m = manager{
		Nodes: make(map[uint]*ws.Node),
	}
)

type manager struct {
	Nodes  map[uint]*ws.Node
	rwNode sync.RWMutex
}

func (m *manager) addNode(node *ws.Node) {
	m.rwNode.Lock()
	defer m.rwNode.Unlock()
	m.Nodes[node.ID] = &ws.Node{
		ID:        node.ID,
		Name:      node.Name,
		ChainName: node.ChainName,
	}
}

func (m *manager) removeNode(node *ws.Node) {
	m.rwNode.Lock()
	defer m.rwNode.Unlock()
	delete(m.Nodes, node.ID)
}

func (m *manager) iterNode(function func(node *ws.Node)) {
	m.rwNode.RLock()
	defer m.rwNode.RUnlock()
	for _, node := range m.Nodes {
		if ok, _ := process.PidExists(node.PID); !ok {
			pid := int32(0)
			processes, _ := process.Processes()
			for _, process := range processes {
				cmd, _ := process.Cmdline()
				if strings.Contains(cmd, fmt.Sprintf("/.baas/%s/deploy/%s/", node.ChainName, node.Name)) {
					pid = process.Pid
					break
				}
			}
			node.PID = pid
		}
		function(node)
	}
}
