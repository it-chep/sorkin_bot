package appointment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sorkin_bot/internal/domain/services/appointment/mis_dto"
)

type AppointmentService struct {
	logger *slog.Logger
	client http.Client
}

func NewAppointmentService(logger *slog.Logger, client http.Client) AppointmentService {
	return AppointmentService{
		logger: logger,
		client: client,
	}
}

func (as *AppointmentService) sendToMIS(ctx context.Context, method string, body io.Reader) []byte {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.sendToMIS"

	// Создание запроса с учетом API_KEY в QueryParams
	var responseBody = make([]byte, 0, 512)
	queryParams := url.Values{}
	queryParams.Set("api_key", os.Getenv("MIS_API_KEY"))

	urlWithMethod, _ := url.JoinPath(os.Getenv("MIS_API_URL"), method)
	urlWithParams := fmt.Sprintf("%s?%s", urlWithMethod, queryParams.Encode())

	request, err := http.NewRequest(http.MethodPost, urlWithParams, body)
	if err != nil {
		// Обработка ошибки создания запроса
		as.logger.Error(fmt.Sprintf("error while create request entity, op: %s", op))
		return responseBody
	}

	// Выполнение запроса
	result, err := as.client.Do(request)
	if err != nil {
		// Обработка ошибки выполнения запроса
		as.logger.Error(fmt.Sprintf("error while do request, op: %s", op))
		return responseBody
	}
	defer result.Body.Close()

	// Чтение тела ответа
	responseBody, err = ioutil.ReadAll(result.Body)
	if err != nil {
		// Обработка ошибки чтения тела ответа
		as.logger.Error(fmt.Sprintf("error while reading response body, op: %s", op))
		return responseBody
	}

	// Базовая превалидация ответа, если 400 или 500
	var baseResponse mis_dto.BaseResponse
	err = json.Unmarshal(responseBody, &baseResponse)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshaling data to BaseResponse struct, op: %s, method %s", op, method))
		return responseBody
	}
	if baseResponse.Error == 1 {
		as.logger.Error(fmt.Sprintf("error while sending request to MIS: \ncode: %d, \ndescription: %s \nop: %s", baseResponse.Data.Code, baseResponse.Data.ErrorDescription, op))
		return responseBody
	}

	return responseBody

}

func (as *AppointmentService) FastAppointment(ctx context.Context) {

}

func (as *AppointmentService) CreateAppointment(ctx context.Context) (err error) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CreateAppointment"
	var Request mis_dto.CreateAppointmentRequest
	var Response mis_dto.CreateAppointmentResponse

	Request = mis_dto.CreateAppointmentRequest{}
	jsonBody, err := json.Marshal(Request)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err
	}
	body := bytes.NewReader(jsonBody)
	responseBody := as.sendToMIS(ctx, mis_dto.CreateAppointmentMethod, body)

	err = json.Unmarshal(responseBody, &Response)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshaling data to BaseResponse struct, op: %s", op))
		return err
	}

	return nil
}

func (as *AppointmentService) CancelAppointment(ctx context.Context, movedTo string, appointmentId int) (err error, result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CancelAppointment"
	var Request mis_dto.CancelAppointmentRequest
	var Response mis_dto.ConfirmAndCancelAppointmentResponse

	Request = mis_dto.CancelAppointmentRequest{
		AppointmentId: appointmentId,
		Source:        mis_dto.Source,
		MovedTo:       movedTo,
	}

	jsonBody, err := json.Marshal(Request)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, false
	}
	body := bytes.NewReader(jsonBody)
	responseBody := as.sendToMIS(ctx, mis_dto.CancelAppointmentMethod, body)

	err = json.Unmarshal(responseBody, &Response)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while unmarshalling data err: %s \nop: %s", err, op))
		return err, Response.Data.True
	}

	return nil, Response.Data.True
}

func (as *AppointmentService) ConfirmAppointment(ctx context.Context, appointmentId int) (err error) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.ConfirmAppointment"
	var Request mis_dto.ConfirmAppointmentRequest
	var Response mis_dto.ConfirmAndCancelAppointmentResponse

	Request = mis_dto.ConfirmAppointmentRequest{
		AppointmentId: appointmentId,
		Source:        mis_dto.Source,
	}
	jsonBody, err := json.Marshal(Request)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err
	}
	body := bytes.NewReader(jsonBody)
	responseBody := as.sendToMIS(ctx, mis_dto.ConfirmAppointmentMethod, body)

	err = json.Unmarshal(responseBody, &Response)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while unmarshalling data err: %s \nop: %s", err, op))
		return err
	}
	return nil
}

func (as *AppointmentService) RescheduleAppointment(ctx context.Context, movedTo string) (err error) {
	err, appointment := as.DetailAppointment(ctx)
	if err != nil {
		return err
	}

	err, _ = as.CancelAppointment(ctx, movedTo, appointment.Data[0].Id)
	if err != nil {
		return err
	}

	return nil
}

func (as *AppointmentService) MyAppointments(ctx context.Context) (err error, Response mis_dto.GetAppointmentsResponse) {
	// Полуаем записи по пользователю и отдаем ему только даты записей
	op := "sorkin_bot.internal.domain.services.appointment.appointment.MyAppointments"
	var Request mis_dto.GetAppointmentsRequest

	Request = mis_dto.GetAppointmentsRequest{}
	jsonBody, err := json.Marshal(Request)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, Response
	}
	body := bytes.NewReader(jsonBody)
	responseBody := as.sendToMIS(ctx, mis_dto.GetAppointmentsMethod, body)

	err = json.Unmarshal(responseBody, &Response)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while unmarshalling data err: %s \nop: %s", err, op))
		return err, Response
	}
	return nil, Response
}

func (as *AppointmentService) DetailAppointment(ctx context.Context) (err error, Response mis_dto.GetAppointmentsResponse) {
	// Полуаем запись по id записи, отдаем данные о записи
	op := "sorkin_bot.internal.domain.services.appointment.appointment.DetailAppointment"
	var Request mis_dto.GetAppointmentsRequest

	Request = mis_dto.GetAppointmentsRequest{}
	jsonBody, err := json.Marshal(Request)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, mis_dto.GetAppointmentsResponse{}
	}
	body := bytes.NewReader(jsonBody)
	responseBody := as.sendToMIS(ctx, mis_dto.GetAppointmentsMethod, body)

	err = json.Unmarshal(responseBody, &Response)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while unmarshalling data err: %s \nop: %s", err, op))
		return err, mis_dto.GetAppointmentsResponse{}
	}
	return nil, Response
}
