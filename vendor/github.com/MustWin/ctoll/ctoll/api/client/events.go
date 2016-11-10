package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/MustWin/ctoll/ctoll/api/v1"
)

type MeterEventClient interface {
	SendUsageSample(apiKey string, e v1.SampleMeterEvent) error
	SendStartMeter(apiKey string, e v1.StartMeterEvent) error
	SendStopMeter(apiKey string, e v1.StopMeterEvent) error
}

type meterEventClient struct {
	*Client
}

func (c *Client) MeterEvents() MeterEventClient {
	return &meterEventClient{c}
}

func (c *meterEventClient) sendMeterEvent(apiKey string, body []byte) ([]byte, error) {
	urlStr, err := c.urls().BuildMeterEvents()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.useAPIKey(r, apiKey)
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

func (c *meterEventClient) SendMachineUsageSample(apiKey string, e v1.MachineSampleMeterEvent) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = c.sendMeterEvent(apiKey, body)
	return err
}

func (c *meterEventClient) SendUsageSample(apiKey string, e v1.SampleMeterEvent) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = c.sendMeterEvent(apiKey, body)
	return err
}

func (c *meterEventClient) SendStartMeter(apiKey string, e v1.StartMeterEvent) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = c.sendMeterEvent(apiKey, body)
	return err
}

func (c *meterEventClient) SendStopMeter(apiKey string, e v1.StopMeterEvent) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = c.sendMeterEvent(apiKey, body)
	return err
}
