package ws

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/olahol/melody"
)

var (
	m  = melody.New()
	am = newAgentManager()
)

// Run web socket
func Run(router *gin.Engine, db *gorm.DB, monitor bool) {
	router.GET("/monitor", func(c *gin.Context) {
		var v = struct {
			Host string
		}{
			c.Request.Host + "/ws/browser",
		}
		fmt.Println(c.Request.Host, v.Host)
		homeTempl.Execute(c.Writer, &v)
	})

	router.GET("/ws/:name", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		am.remove(s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var err error
		rmsg := &Message{}
		agent := am.get(s)

		defer func() {
			if err != nil {
				msg := NewMsg(RespRefuseMsg, fmt.Sprintf("agent %d %s", rmsg.Method, err.Error()))
				s.Write(msg.Bytes())
			} else if monitor {
				// 广播到浏览器,方便查看上报信息
				m.BroadcastFilter([]byte(fmt.Sprintf("%d: %s", agent.ID, string(rmsg.Payload))), func(q *melody.Session) bool {
					return q.Request.URL.Path == "/ws/browser"
				})
			}
		}()

		rmsg, err = MsgFromBytes(msg)
		switch rmsg.Method {
		case ReqJoinMsg:
			agent = &AgentStatus{}
			err = rmsg.DecodePayload(agent)
			if err == nil {
				if !am.add(s, agent) {
					err = fmt.Errorf("agent %d already exist", agent.ID)
				} else {
					magent := &model.Agent{}
					err = db.First(magent, agent.ID).Error
				}
				if err == nil {
					// TODO
					var chainDeployNodes []*model.ChainDeployNode
					db.Where(&model.ChainDeployNode{
						AgentID: agent.ID,
					}).Find(&chainDeployNodes)
					var nodes []*Node
					var chainDeploy model.ChainDeploy
					var chain model.Chain
					for _, chainDeployNode := range chainDeployNodes {
						id := chainDeployNode.ID
						if err := db.First(&chainDeploy, chainDeployNode.ChainDeployID).Error; err != nil {
							continue
						}
						if err := db.First(&chain, chainDeploy.ChainID).Error; err != nil {
							continue
						}
						name := chainDeploy.Name
						nodes = append(nodes, &Node{
							ID:        id,
							Name:      name,
							ChainName: chain.Name,
						})
					}
					if len(nodes) > 0 {
						msg := NewMsg(RespAddNodeMsg, nodes)
						s.Write(msg.Bytes())
					}
				}
			}
		case ReqServerMsg:
			ss := &StatusServer{}
			err = rmsg.DecodePayload(ss)
			if err == nil {
				bts, _ := json.Marshal(ss)
				db.Model(&model.Agent{
					Model: gorm.Model{
						ID: agent.ID,
					},
				}).Updates(&model.Agent{
					Status: string(bts),
				})
			}
		case ReqProcessMsg:
			sp := &StatusProcess{}
			err = rmsg.DecodePayload(sp)
			if err == nil {
				bts, _ := json.Marshal(sp)
				db.Model(&model.ChainDeployNode{
					AgentID:       agent.ID,
					ChainDeployID: sp.NodeID,
				}).Updates(&model.Agent{
					Status: string(bts),
				})
			}
		case ReqNodeInfoMsg:
			sr := &StatusRPC{}
			err = rmsg.DecodePayload(sr)
			if err == nil {
				db.Model(&model.ChainDeployNodeStatus{
					ChainDeployNodeID: sr.NodeID,
					ChainStatusID:     sr.RPCID,
				}).Updates(&model.ChainDeployNodeStatus{
					Value: sr.Value,
				})
			}
		case ReqCommandMsg:
			sc := &StatusCmd{}
			err = rmsg.DecodePayload(sc)
			if err == nil {
				//fmt.Println(sc.Cmd, sc.Value, sc.Error)
			}
		default:
			err = fmt.Errorf("not support method %d", rmsg.Method)
		}
	})
}

// HandleCommand 处理命令
func HandleCommand(agentID uint, cmd string) (string, error) {
	s := am.getsession(agentID)
	if s == nil {
		return cmd, fmt.Errorf("not connected")
	}
	msg := NewMsg(RespCommandMsg, &StatusCmd{
		Cmd: cmd,
	})
	err := s.Write(msg.Bytes())
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

// HandleAddNode 处理命令
func HandleAddNode(agentID uint, nodeID uint, nodeName string, chainName string) error {
	s := am.getsession(agentID)
	if s == nil {
		return fmt.Errorf("not connected")
	}
	msg := NewMsg(RespAddNodeMsg, &Node{
		ID:        nodeID,
		Name:      nodeName,
		ChainName: chainName,
	})
	err := s.Write(msg.Bytes())
	return err
}

// HandleRemoveNode 处理命令
func HandleRemoveNode(agentID uint, nodeID uint) error {
	s := am.getsession(agentID)
	if s == nil {
		return fmt.Errorf("not connected")
	}
	msg := NewMsg(RespRemoveNodeMsg, &Node{
		ID: nodeID,
	})
	err := s.Write(msg.Bytes())
	return err
}
func getClientIP(req *http.Request) string {
	ip := req.Header.Get("x-forwarded-for")
	if ip == "" || len(ip) == 0 || strings.ToLower(ip) == "unknown" {
		ip = req.Header.Get("Proxy-Client-IP")
	}
	if ip == "" || len(ip) == 0 || strings.ToLower(ip) == "unknown" {
		ip = req.Header.Get("WL-Proxy-Client-IP")
	}
	if ip == "" || len(ip) == 0 || strings.ToLower(ip) == "unknown" {
		ip = req.Header.Get("HTTP_CLIENT_IP")
	}
	if ip == "" || len(ip) == 0 || strings.ToLower(ip) == "unknown" {
		ip = req.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if ip == "" || len(ip) == 0 || strings.ToLower(ip) == "unknown" {
		ip = req.RemoteAddr
	}
	// addr, err := net.ResolveTCPAddr("tcp", ip)
	// if err != nil {
	// 	ip = ""
	// } else {
	// 	ip = addr.IP.String()
	// }
	return ip
}

var homeTempl = template.Must(template.New("").Parse(homeHTML))

const homeHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>Chat Example</title>
<script type="text/javascript">
window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        conn.send(msg.value);
        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://{{.Host}}");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};
</script>
<style type="text/css">
html {
    overflow: hidden;
}

body {
    overflow: hidden;
    padding: 0;
    margin: 0;
    width: 100%;
    height: 100%;
    background: gray;
}

#log {
    background: white;
    margin: 0;
    padding: 0.5em 0.5em 0.5em 0.5em;
    position: absolute;
    top: 0.5em;
    left: 0.5em;
    right: 0.5em;
    bottom: 3em;
    overflow: auto;
}

#form {
    padding: 0 0.5em 0 0.5em;
    margin: 0;
    position: absolute;
    bottom: 1em;
    left: 0px;
    width: 100%;
    overflow: hidden;
}

</style>
</head>
<body>
<div id="log"></div>
<form id="form">
    <input type="submit" value="Send" />
    <input type="text" id="msg" size="64"/>
</form>
</body>
</html>
`
