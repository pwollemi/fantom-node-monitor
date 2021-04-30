package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/flashguru-git/node-monitor/config"
	"github.com/flashguru-git/node-monitor/log"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

func sendPostRequest(serverURL string, data map[string]interface{}) {
	client := &http.Client{}

	b, _ := json.Marshal(data)
	rq, err := http.NewRequest("POST", serverURL, bytes.NewReader(b))
	if err != nil {
		log.Errorln("An error occured while posting data" + string(b))
		return
	}

	rq.Header.Set("Content-Type", "application/json")

	rp, err := client.Do(rq)
	if err != nil || rp == nil {
		log.Errorln("An error occured while posting data")
	}

	defer func(r *http.Response) {
		if r != nil && r.Body != nil {
			_, _ = ioutil.ReadAll(r.Body)
			_ = r.Body.Close()
		}
		if r := recover(); r != nil {
			log.Infoln("Recovered in postSyncEvent", r)
		}
	}(rp)
}

func queryConsole(query string) (string, error) {
	lachesisConsole := config.Config().GetString("LACHESIS_CONSOLE")
	// command := fmt.Sprintf(`echo "%v" | %v | tr -d '\n'`, key, lachesisConsole)
	command := fmt.Sprintf(`%v --exec "%v"`, lachesisConsole, query)
	cmd := exec.Command("/bin/sh", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}

	var b bytes.Buffer
	io.Copy(&b, stdout)
	return b.String(), nil
}

func getNodeId() string {
	res, err := queryConsole("admin.nodeInfo.id")
	if err != nil {
		log.Errorln(err.Error())
		return ""
	}
	res = strings.ReplaceAll(res, "\n", "")
	return strings.ReplaceAll(res, `"`, "")
}

func getBlockNumber() uint64 {
	res, err := queryConsole("ftm.blockNumber")
	if err != nil {
		log.Errorln(err.Error())
		return 0
	}
	res = strings.ReplaceAll(res, "\n", "")
	blockNumber, err := strconv.ParseUint(res, 10, 64)
	if err != nil {
		log.Errorln(err.Error())
		return 0
	}
	return blockNumber
}

func getTopPeersBlockHeight() uint64 {
	res, err := queryConsole("admin.peers")
	if err != nil {
		log.Errorln(err.Error())
		return 0
	}
	re := regexp.MustCompile(`\s+blocks:\s+([0-9]+),`)
	var topHeight uint64 = 0
	for _, match := range re.FindAllStringSubmatch(res, -1) {
		if len(match) < 2 {
			continue
		}
		if height, err := strconv.ParseUint(match[1], 10, 64); err == nil {
			if height > topHeight {
				topHeight = height
			}
		}
	}
	return topHeight
}

func getMemoryUsage() map[string]interface{} {
	memory, err := memory.Get()
	if err != nil {
		return nil
	}
	return map[string]interface{}{
		"total":  memory.Total,
		"used":   memory.Used,
		"cached": memory.Cached,
		"free":   memory.Free,
	}
}

func getCpuUsage() map[string]interface{} {
	data, err := cpu.Get()
	if err != nil {
		return nil
	}
	return map[string]interface{}{
		"user":   data.User,
		"system": data.System,
		"total":  data.Total,
		"idle":   data.Idle,
		"nice":   data.Nice,
	}
}

func getIpAddr() string {
	cmd := exec.Command("/bin/sh", "-c", "curl -4 icanhazip.com")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return ""
	}
	if err := cmd.Start(); err != nil {
		return ""
	}

	var b bytes.Buffer
	io.Copy(&b, stdout)
	return strings.ReplaceAll(b.String(), "\n", "")
}

func getLocalIpAddrs() (res []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Errorln("Oops: " + err.Error() + "\n")
		return
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				res = append(res, ipnet.IP.String())
			}
		}
	}
	return res
}
