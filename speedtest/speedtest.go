package speedtest

import (
    "errors"
    "net/http"
    "encoding/xml"
    "io/ioutil"
    "strings"
    "time"
    "fmt"
)

const (
    bitCount = 8
    bitSize = 64
    ConfigURL = "http://www.speedtest.net/speedtest-config.php"
)

var (
    sizes = []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000}
)

type Server struct {
    URL      string  `xml:"url,attr"`
    Lat      string  `xml:"lat,attr"`
    Lon      string  `xml:"lon,attr"`
    Name     string  `xml:"name,attr"`
    Country  string  `xml:"country,attr"`
    Sponsor  string  `xml:"sponsor,attr"`
    ID       string  `xml:"id,attr"`
    URL2     string  `xml:"url2,attr"`
    Host     string  `xml:"host,attr"`
    Distance float64
    DLSpeed  float64
    ULSpeed  float64
}

type ServerList struct {
    Servers []Server `xml:"servers->server"`
}

type ClientConfig struct {
    IP  string `xml:"ip,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
	ISP string `xml:"isp,attr"`
}

type Config struct {
    Client ClientConfig `xml:"settings->client"`
}

// GetConfig gets client configuration XML from Speedtest.net servers.
func GetConfig() (*ConfigXML, error) {
	config, err := getXML(ConfigURL, &ConfigXML{})
	if err != nil {
		return nil, err
	}
	return config.(*ConfigXML), nil
}

// GetBestServer finds the "best server" for the client.
//
// The best server is defined as the closest server to the
// client in terms of lon/lat.
func GetBestServer(config *ConfigXML) (*Server, error) {
	var t, b float64
	var best string
	b = 2 << 63
	sl, err := getServerList()
	if err != nil {
		return nil, err
	}
	sm := make(map[string]Server, len(sl.Servers.Values))
	for _, s := range sl.Servers.Values {
		t, err = s.getDistance(config)
		if err != nil {
			continue
		}
		if t < b {
			b = t
			sm[s.ID] = s
			best = s.ID
		}
	}
	bs := sm[best]
	return &bs, nil
}

// getServerList gets the server list XML from Speedtest.net servers
func GetServerList() (*ServerListXML, error) {
    list, err := getXML(ServerListURL, &ServerListXML{})
    if err != nil {
        return nil, err
    }
    return list.(*ServerList), nil
}

func DownloadSpeed(s *Server) (mbs float64) {
    var d []byte
    var url string
    var total float64

    queue := make(chan []byte, len(sizes))
    start := time.Now()

    for _, v := range sizes {
        url = getDownloadURL(s.URL, v)
        go func(url string, ch chan<- []byte) {
            resp, err := http.Get(url)
            if err != nil {
                ch <- []byte{}
                return
            }
            defer resp.Body.Close()
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                ch <- []byte{}
                return
            }
            ch <- body
        }(url, queue)
    }

    files := make([][]byte, len(sizes))
    total = 0
    for i := 0; i < len(sizes); i++ {
        d = <-queue
        files[i] = d
        total += float64(len(d))
    }
    close(queue)

    dur := time.Since(start)
    mbs = ((total * bitCount) / (1e6)) / dur.Seconds()
    return
}

// getDownloadURL returns the download URL of an image of
// a certain size s to be downloaded from Speedtest.net servers.
func getDownloadURL(url string, s int) string {
    base := strings.Replace(url, "upload.php", "", 1)
    return fmt.Sprintf("%srandom%dx%d.jpg", base, s, s)
}

func Start() float64 {
    config, _ := GetConfig()
    serverList, _ := GetServerList()
    bestServer := GetBestServer(config, serverList)
    return DownloadSpeed(bestServer)
}
