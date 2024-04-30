package mis_reno

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
	mis_dto "sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"time"
)

type MisRenoGateway struct {
	logger *slog.Logger
	client http.Client
}

//todo добавить кеш, потому что запрос долго идет

func NewMisRenoGateway(logger *slog.Logger, client http.Client) MisRenoGateway {
	return MisRenoGateway{
		logger: logger,
		client: client,
	}
}

func (mg *MisRenoGateway) sendToMIS(ctx context.Context, method string, body io.Reader) []byte {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.sendToMIS"

	// Создание запроса с учетом API_KEY в QueryParams
	var responseBody = make([]byte, 0, 512)
	queryParams := url.Values{}
	queryParams.Set("api_key", os.Getenv("MIS_API_KEY"))

	urlWithMethod, _ := url.JoinPath(os.Getenv("MIS_API_URL"), method)
	urlWithParams := fmt.Sprintf("%s?%s", urlWithMethod, queryParams.Encode())

	var data map[string]interface{}
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while decoding JSON: %s \nplace: %s", err, op))
		return responseBody
	}

	formValues := url.Values{}
	for key, value := range data {
		strValue := fmt.Sprintf("%v", value)
		formValues.Add(key, strValue)
	}
	mg.logger.Info(fmt.Sprintf("formValues %s. Data: %s", formValues, data))

	requestBody := bytes.NewBufferString(formValues.Encode())
	mg.logger.Info(fmt.Sprintf("REQUEST BODY %s", requestBody))
	request, err := http.NewRequest(http.MethodPost, urlWithParams, requestBody)
	if err != nil {
		// Обработка ошибки создания запроса
		mg.logger.Error(fmt.Sprintf("error while create request entity, op: %s", op))
		return responseBody
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Выполнение запроса
	result, err := mg.client.Do(request)
	if err != nil {
		// Обработка ошибки выполнения запроса
		mg.logger.Error(fmt.Sprintf("error while do request, op: %s", op))
		return responseBody
	}
	defer result.Body.Close()

	// Чтение тела ответа
	responseBody, err = ioutil.ReadAll(result.Body)
	if err != nil {
		// Обработка ошибки чтения тела ответа
		mg.logger.Error(fmt.Sprintf("error while reading response body, op: %s", op))
		return responseBody
	}
	mg.logger.Info(fmt.Sprintf("RESPONSE BODY %s", responseBody))

	// Базовая превалидация ответа, если 400 или 500
	var baseResponse mis_dto.BaseResponse
	err = json.Unmarshal(responseBody, &baseResponse)
	if baseResponse.Error == 1 {
		mg.logger.Error(fmt.Sprintf("error while sending request to MIS: \ncode: %d, \ndescription: %s \nop: %s", baseResponse.Data.Code, baseResponse.Data.ErrorDescription, op))
		return responseBody
	}

	return responseBody

}

func (mg *MisRenoGateway) FastAppointment(ctx context.Context) {

}

func (mg *MisRenoGateway) CreateAppointment(ctx context.Context) (err error, appointmentId int) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CreateAppointment"
	var response mis_dto.CreateAppointmentResponse
	var request = mis_dto.CreateAppointmentRequest{}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, -1
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.CreateAppointmentMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshaling data to BaseResponse struct, op: %s", op))
		return err, -1
	}

	return nil, response.Data.ID
}

func (mg *MisRenoGateway) CancelAppointment(ctx context.Context, movedTo string, appointmentId int) (err error, result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CancelAppointment"
	var request = mis_dto.CancelAppointmentRequest{
		AppointmentId: appointmentId,
		Source:        mis_dto.Source,
		MovedTo:       movedTo,
	}
	var response mis_dto.ConfirmAndCancelAppointmentResponse

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, false
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.CancelAppointmentMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while unmarshalling data err: %s \nop: %s", err, op))
		return err, response.Data.True
	}

	return nil, response.Data.True
}

func (mg *MisRenoGateway) ConfirmAppointment(ctx context.Context, appointmentId int) (err error, result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.ConfirmAppointment"
	var request = mis_dto.ConfirmAppointmentRequest{
		AppointmentId: appointmentId,
		Source:        mis_dto.Source,
	}
	var response mis_dto.ConfirmAndCancelAppointmentResponse

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, false
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.ConfirmAppointmentMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while unmarshalling data err: %s \nop: %s", err, op))
		return err, response.Data.True
	}
	return nil, response.Data.True
}

func (mg *MisRenoGateway) MyAppointments(ctx context.Context, user entity.User) (err error, appointments []appointment.Appointment) {
	// Полуаем записи по пользователю и отдаем ему только даты записей
	op := "sorkin_bot.internal.domain.services.appointment.appointment.MyAppointments"
	var response mis_dto.GetAppointmentsResponse
	currentTime := time.Now()

	// todo рассматриваем только записи из бота, то есть человек будет получать только доступ к тем записям, которые были созданы им из бота

	var request = mis_dto.GetAppointmentsRequest{
		DateCreatedFrom: user.GetRegistrationTime(),
		DateCreatedTo:   fmt.Sprintf("%02d.%02d.%d %02d:%02d", currentTime.Day(), currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute()),
		PatientId:       user.GetPatientId(),
		StatusId:        "1, 2",
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, appointments
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetAppointmentsMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while unmarshalling data err: %s \nop: %s", err, op))
		return err, appointments
	}
	for _, appointmentDTO := range response.Data {
		appointments = append(appointments, appointmentDTO.ToDomain())
	}
	return nil, appointments
}

func (mg *MisRenoGateway) DetailAppointment(ctx context.Context, user entity.User, appointmentId int) (err error, appointmentEntity appointment.Appointment) {
	// Полуаем запись по id записи, отдаем данные о записи
	op := "sorkin_bot.internal.domain.services.appointment.appointment.DetailAppointment"
	var response mis_dto.GetAppointmentsResponse
	currentTime := time.Now()

	// todo мб сделать конструктор для GetAppointmentsRequest

	var request = mis_dto.GetAppointmentsRequest{
		AppointmentId:   appointmentId,
		DateCreatedFrom: user.GetRegistrationTime(),
		DateCreatedTo:   fmt.Sprintf("%02d.%02d.%d %02d:%02d", currentTime.Day(), currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute()),
		PatientId:       user.GetPatientId(),
		StatusId:        "1, 2",
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, appointmentEntity
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetAppointmentsMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while unmarshalling data err: %s \nop: %s", err, op))
		return err, appointmentEntity
	}

	appointmentEntity = response.Data[0].ToDomain()

	return nil, appointmentEntity
}
