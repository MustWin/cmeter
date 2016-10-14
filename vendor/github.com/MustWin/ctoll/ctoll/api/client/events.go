package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/MustWin/ctoll/ctoll/api/v1"
)

type MeterEventClient interface {
	SendUsageSample(e v1.SampleMeterEvent) error
	SendStartMeter(e v1.StartMeterEvent) error
	SendStopMeter(e v1.StopMeterEvent) error
}

type meterEventClient struct {
	*Client
}

func (c *Client) MeterEvents() MeterEventClient {
	return &meterEventClient{c}
}

func (c *meterEventClient) sendMeterEvent(body []byte) ([]byte, error) {
	urlStr, err := c.urls().BuildMeterEvents()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, err
}

func (c *meterEventClient) SendUsageSample(e v1.SampleMeterEvent) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = c.sendMeterEvent(body)
	return err
}

func (c *meterEventClient) SendStartMeter(e v1.StartMeterEvent) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = c.sendMeterEvent(body)
	return err
}

func (c *meterEventClient) SendStopMeter(e v1.StopMeterEvent) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = c.sendMeterEvent(body)
	return err
}
