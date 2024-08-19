package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	c   *http.Client
	key *string
}

type Building struct {
	ID        int
	Name      string
	Structure struct {
		Floors  []*Floor
		Devices []*Device
		Areas   []*Area
	}
}

type Floor struct {
	ID         int
	Name       string
	BuildingId int
	Devices    []*Device
}

type Area struct {
	ID             int
	Name           string
	BuildingId     int
	FloorId        int
	AccessLevel    int
	DirectAccess   bool
	EndDate        string
	MinTemperature float64
	MaxTemperature float64
	Expanded       bool
	Devices        []*Device
}

type Device struct {
	DeviceID   int
	DeviceName string
	BuildingID int
	Device     struct {
		RoomTemperature float64
		SetTemperature  float64
		Power           bool
		OperationMode   int
	}
}

type LoginRequest struct {
	AppVersion      string
	CaptchaResponse *string
	Email           string
	Language        int
	Password        string
	Persist         bool
}

type LoginResponse struct {
	ErrorId   *int
	LoginData struct {
		ContextKey string
	}
}

var (
	baseURL = "https://app.melcloud.com/Mitsubishi.Wifi.Client"
)

func NewClient() *Client {
	return &Client{c: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) Login(email, password string) error {
	body := LoginRequest{
		AppVersion:      "1.21.6.0",
		CaptchaResponse: nil,
		Email:           email,
		Language:        0,
		Password:        password,
		Persist:         false,
	}

	var resBody LoginResponse
	if err := c.req("/Login/ClientLogin", &body, &resBody, false); err != nil {
		return err
	}

	if resBody.ErrorId != nil {
		return fmt.Errorf("bad email/password combo, error id: %d", *resBody.ErrorId)
	}

	c.key = &resBody.LoginData.ContextKey
	return nil
}

func (c *Client) Devices() ([]*Device, error) {
	var buildings []*Building
	if err := c.req("/User/ListDevices", nil, &buildings, true); err != nil {
		return nil, err
	}

	var devices []*Device
	for _, building := range buildings {

		devices = append(devices, building.Structure.Devices...)

		if building.Structure.Floors != nil {
			for _, floor := range building.Structure.Floors {
				devices = append(devices, floor.Devices...)
			}
		}

		if building.Structure.Areas != nil {
			for _, area := range building.Structure.Areas {
				devices = append(devices, area.Devices...)
			}
		}
	}

	return devices, nil
}

func (c *Client) req(path string, body interface{}, resBody interface{}, auth bool) error {
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}

		r = bytes.NewBuffer(data)
	}

	method := "GET"
	if body != nil {
		method = "POST"
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, path), r)
	if err != nil {
		return err
	}
	if r != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth {
		req.Header.Set("X-MitsContextKey", *c.key)
	}

	res, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(resBody); err != nil {
		return err
	}

	return nil
}
