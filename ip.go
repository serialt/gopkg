package gopkg

import (
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
)

const ipv4_regex = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`

// 获取所在区域的公网ip的网站
var urlList = []string{
	"https://ip.tool.lu",
	"http://cip.cc",
}

// IPGet 返回客户端 IP
func IPGet(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

// 获取公网ip, 如果两个ip不同，则访问ip.tool.lu 和false, ip都相同则返回true
func GetPubIP() (ip string, all bool) {
	var ipList []string
	for _, url := range urlList {
		client := &http.Client{}
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
		request.Header.Add("Connection", "keep-alive")
		request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
		resp, err := client.Do(request)
		if resp.StatusCode != 200 && err != nil {
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		reg := regexp.MustCompile(ipv4_regex)
		ipList = reg.FindAllString(string(body), -1)

	}
	if len(ipList) > 0 {
		// fmt.Printf("my public ip is: %s\n", ipList[0])
		ip = ipList[0]
		if ipList[0] == ipList[1] {
			all = true

		} else {
			all = false
		}

	}
	return
}
