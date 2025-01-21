package wau_sms

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"io/ioutil"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/clients/sms_gateway/wau_sms/sms_dto"
	"sorkin_bot/internal/config"
	"time"
)

const (
	name = "UMCL"
)

type Sender struct {
	client   http.Client
	logger   *slog.Logger
	url      string
	user     string
	password string
}

func NewSender(logger *slog.Logger, appConfig config.Config) *Sender {
	return &Sender{
		client: http.Client{
			Timeout: time.Second * 3,
		},
		logger:   logger,
		url:      appConfig.WAUSMS.URL,
		user:     appConfig.WAUSMS.User,
		password: appConfig.WAUSMS.Password,
	}
}

func (s *Sender) SendNotification(ctx context.Context, to []string, message string) error {
	var basicResponse sms_dto.BasicResponse
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	dto := sms_dto.BasicRequest{
		To:   to,
		Text: message,
		From: name,
	}

	requestBody, err := json.Marshal(dto)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error while marshalling response: %s", err.Error()))
		return err
	}

	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodPost, s.url, bytes.NewBuffer(requestBody))
	if err != nil {
		s.logger.Error(fmt.Sprintf("error creating request: %s", err.Error()))
		return err
	}

	s.setHeaders(req)

	response, err := s.client.Do(req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error while doing response: %s", err.Error()))
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error reading response: %s", err.Error()))
		return err
	}

	if !lo.Contains([]int{http.StatusAccepted, http.StatusMultiStatus}, response.StatusCode) {
		s.logger.Error(fmt.Sprintf("error while send sms notification %s", string(body)))
		return errors.New("error while send sms notification " + string(body))
	}

	err = json.Unmarshal(body, &basicResponse)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error while unmarshalling JSON sms %s", err.Error()))
		return err
	}

	s.logger.Info(fmt.Sprintf("successfully sending notification: %v", basicResponse))

	return nil
}

func (s *Sender) generateBasicAuth(user, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))
}

func (s *Sender) getWauSmsConfig() config.WAUSMSConfig {
	return config.NewConfig().WAUSMS
}

func (s *Sender) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+s.generateBasicAuth(s.user, s.password))
}
