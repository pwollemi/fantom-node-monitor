package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"

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

func queryNode(key string) (string, error) {
	lachesisConsole := config.Config().GetString("LACHESIS_CONSOLE")
	command := fmt.Sprintf(`echo "%v" | %v | tr -d '\n'`, key, lachesisConsole)
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
	res, err := queryNode("admin.nodeInfo.id")
	if err != nil {
		log.Errorln(err.Error())
		return ""
	}
	re := regexp.MustCompile(`>\s+"([0-9a-zA-Z]+)">`)
	match := re.FindStringSubmatch(res)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func getBlockNumber() uint64 {
	res, err := queryNode("ftm.blockNumber")
	if err != nil {
		log.Errorln(err.Error())
		return 0
	}
	re := regexp.MustCompile(`>\s+([0-9]+)>`)
	match := re.FindStringSubmatch(res)
	if len(match) < 2 {
		return 0
	}
	blockNumber, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		log.Errorln(err.Error())
		return 0
	}
	return blockNumber
}

func getTopPeersBlockHeight() uint64 {
	res, err := queryNode("admin.peers")
	if err != nil {
		log.Errorln(err.Error())
		return 0
	}
	re := regexp.MustCompile(`\s+blocks:\s+([0-9]+)`)
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
