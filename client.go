package fronius

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Host    string
	Timeout time.Duration
}

type Client struct {
	*http.Client
	Config Config
}

func NewClient(c Config) *Client {
	return &Client{
		Client: newTimeoutClient(c.Timeout, c.Timeout),
		Config: c,
	}
}

func parseResponse(res *http.Response, v interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (c *Client) GetInverterSystemRealtimeData() (*InverterSystemResponse, error) {
	res, err := c.Client.Get(InverterSystemRealtimeDataRequestURL(c.Config.Host))
	if err != nil {
		return nil, err
	}

	var inv InverterSystemResponse
	err = parseResponse(res, &inv)
	return &inv, err
}

func (c *Client) GetInverterDeviceRealtimeData(deviceId int, collection string) (*InverterDeviceResponse, error) {
	res, err := c.Client.Get(InverterDeviceRealtimeDataRequestURL(c.Config.Host, deviceId, collection))
	if err != nil {
		return nil, err
	}

	var inv InverterDeviceResponse
	err = parseResponse(res, &inv)
	return &inv, err
}

func (c *Client) get(url string, resp interface{}) error {
	res, err := c.Client.Get(url)
	if err != nil {
		return err
	}

	err = parseResponse(res, resp)
	return err
}

func (c *Client) Get3PInverterData(deviceId int) (inv Inverter3PhasesDeviceResponse, err error) {
	url := InverterDeviceRealtimeDataRequestURL(c.Config.Host, deviceId, Collection3Phases)
	err = c.get(url, &inv)
	return
}

func (c *Client) GetCommonInverterData(deviceId int) (inv InverterCommonDeviceResponse, err error) {
	url := InverterDeviceRealtimeDataRequestURL(c.Config.Host, deviceId, CollectionCommon)
	err = c.get(url, &inv)
	return
}

func (c *Client) requestURL(endpoint string, arguments ...string) string {
	base := path(toHTTPHost(c.Config.Host), solarAPIFormat, endpoint)
	arg := args(arguments...)
	return append(base, arg)
}

func (c *Client) GetMeterSystemRealtimeData() (*MeterSystemResponse, error) {
	res, err := c.Client.Get(c.requestURL("GetMeterRealtimeData.cgi"))
	if err != nil {
		return nil, err
	}

	var inv MeterSystemResponse
	err = parseResponse(res, &inv)
	return &inv, err
}

func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		log.Println("prettyPring:", err)
	} else {
		log.Printf("%s\n", b)
	}
}

// from http://stackoverflow.com/a/16930649/424814
func timeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func newTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
}
