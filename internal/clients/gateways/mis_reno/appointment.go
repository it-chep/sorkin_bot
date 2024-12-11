package mis_reno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sorkin_bot/internal/clients/gateways/dto"
	"sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	"strconv"
	"time"
)

const (
	ScheduleStep10Min    = 10
	ScheduleStep15Min    = 15
	ScheduleStep20Min    = 20
	ScheduleStepHalfHour = 30
	ScheduleStepHour     = 60
)

func JsonMarshaller[T any](req T, op string, logger *slog.Logger) []byte {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return nil
	}
	return jsonBody
}

func JsonUnMarshaller[T any](resp T, respBody []byte, op string, logger *slog.Logger) (T, error) {
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		logger.Error(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return resp, err
	}
	return resp, nil
}

type MisRenoGateway struct {
	logger *slog.Logger
	client http.Client
}

func NewMisRenoGateway(logger *slog.Logger, client http.Client) MisRenoGateway {
	return MisRenoGateway{
		logger: logger,
		client: client,
	}
}

func (mg *MisRenoGateway) sendToMIS(ctx context.Context, method string, body []byte) []byte {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.sendToMIS"

	// Создание запроса с учетом API_KEY в QueryParams
	var responseBody = make([]byte, 0, 512)
	queryParams := url.Values{}
	queryParams.Set("api_key", os.Getenv("MIS_API_KEY"))

	urlWithMethod, _ := url.JoinPath(os.Getenv("MIS_API_URL"), method)
	urlWithParams := fmt.Sprintf("%s?%s", urlWithMethod, queryParams.Encode())

	var data map[string]interface{}
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&data)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while decoding JSON: %s place: %s", err, op))
		return responseBody
	}

	formValues := url.Values{}
	for key, value := range data {
		strValue := fmt.Sprintf("%v", value)
		formValues.Add(key, strValue)
	}

	requestBody := bytes.NewBufferString(formValues.Encode())
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, urlWithParams, requestBody)
	if err != nil {
		// Обработка ошибки создания запроса
		mg.logger.Error(fmt.Sprintf("error while create request entity, %v op: %s", request, op))
		return responseBody
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Выполнение запроса
	result, err := mg.client.Do(request)
	if err != nil {
		// Обработка ошибки выполнения запроса
		mg.logger.Error(fmt.Sprintf("error while do request %v, op: %s", result, op))
		return responseBody
	}
	defer result.Body.Close()

	// Чтение тела ответа
	responseBody, err = ioutil.ReadAll(result.Body)
	if err != nil {
		// Обработка ошибки чтения тела ответа
		mg.logger.Error(fmt.Sprintf("error while reading response body, %v, op: %s", responseBody, op))
		return responseBody
	}

	// Базовая превалидация ответа, если 400 или 500
	var baseResponse mis_dto.BaseResponse
	err = json.Unmarshal(responseBody, &baseResponse)
	if baseResponse.Error == 1 {
		mg.logger.Error(fmt.Sprintf("error while sending request to MIS: \ncode: %d, \ndescription: %s \nop: %s", baseResponse.Data.Code, baseResponse.Data.ErrorDescription, op))
		return responseBody
	}

	return responseBody

}

func (mg *MisRenoGateway) CreateAppointment(ctx context.Context, createAppointmentDTO dto.CreateAppointmentDTO) (appointmentId *int, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CreateAppointment"
	var response mis_dto.CreateAppointmentResponse
	var request = mis_dto.CreateAppointmentRequest{
		PatientId: createAppointmentDTO.PatientId,
		DoctorId:  createAppointmentDTO.DoctorId,
		ClinicId:  mis_dto.DefaultClinicId,
		TimeStart: createAppointmentDTO.TimeStart,
		TimeEnd:   createAppointmentDTO.TimeEnd,
	}

	if createAppointmentDTO.OnlineAppointment {
		isTelemedicine := true
		request.IsTelemedicine = &isTelemedicine
	}
	if createAppointmentDTO.HomeVisit {
		homeAddress := fmt.Sprintf("Address is %s", createAppointmentDTO.HomeAddress)
		request.Comment = &homeAddress
		isOutside := true
		request.IsOutside = &isOutside
	}

	responseBody := mg.sendToMIS(ctx, mis_dto.CreateAppointmentMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return nil, err
	}

	id, _ := strconv.Atoi(response.Data)

	return &id, nil
}

func (mg *MisRenoGateway) CancelAppointment(ctx context.Context, movedTo string, appointmentId int) (result bool, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CancelAppointment"
	var request = mis_dto.CancelAppointmentRequest{
		AppointmentId: appointmentId,
		MovedTo:       movedTo,
	}
	var response mis_dto.ConfirmAndCancelAppointmentResponse

	responseBody := mg.sendToMIS(ctx, mis_dto.CancelAppointmentMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return false, err
	}
	return true, err
}

func (mg *MisRenoGateway) ConfirmAppointment(ctx context.Context, appointmentId int) (result bool, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.ConfirmAppointment"
	var request = mis_dto.ConfirmAppointmentRequest{
		AppointmentId: appointmentId,
	}
	var response mis_dto.ConfirmAndCancelAppointmentResponse

	responseBody := mg.sendToMIS(ctx, mis_dto.ConfirmAppointmentMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (mg *MisRenoGateway) MyAppointments(ctx context.Context, patientId int, registrationTime string) (appointments []dto.AppointmentDTO, err error) {
	// Полуаем записи по пользователю и отдаем ему только даты записей
	op := "sorkin_bot.internal.domain.services.appointment.appointment.MyAppointments"
	var response mis_dto.GetAppointmentsResponse
	// Получаем текущее время
	currentTimeUTC := time.Now().UTC()

	// Получаем время в лисабоне(возможно хардкод, но гибкости пока не требуется)
	location, err := time.LoadLocation("Europe/Lisbon")
	if err != nil {
		mg.logger.Error("Error loading location:", err)
		return
	}

	// Преобразуем текущее время в локацию GMT+1
	currentTime := currentTimeUTC.In(location)

	if currentTime.Hour() == currentTimeUTC.Hour() {
		currentTime = currentTime.Add(time.Hour)
	}

	var request = mis_dto.GetAppointmentsRequest{
		DateCreatedFrom: registrationTime,
		DateCreatedTo:   fmt.Sprintf("%02d.%02d.%d %02d:%02d", currentTime.Day(), currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute()),
		PatientId:       patientId,
		StatusId:        mis_dto.ActiveStatusIDs,
	}
	mg.logger.Info("REQUEST", request)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetAppointmentsMethod, JsonMarshaller(request, op, mg.logger))
	mg.logger.Info("RESPONSE BODY", responseBody)

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	mg.logger.Info("RESPONSE", response)

	if err != nil {
		return appointments, err
	}

	for _, appointment := range response.Data {
		appointments = append(appointments, appointment.ToDTO())
	}

	return appointments, nil
}

func (mg *MisRenoGateway) DetailAppointment(ctx context.Context, patientId, appointmentId int, registrationTime string) (appointmentDTO dto.AppointmentDTO, err error) {
	// Полуаем запись по id записи, отдаем данные о записи
	op := "sorkin_bot.internal.domain.services.appointment.appointment.DetailAppointment"
	currentTime := time.Now()
	var response mis_dto.GetAppointmentsResponse
	var request = mis_dto.GetAppointmentsRequest{
		AppointmentId:   appointmentId,
		DateCreatedFrom: registrationTime,
		DateCreatedTo:   fmt.Sprintf("%02d.%02d.%d %02d:%02d", currentTime.Day(), currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute()),
		PatientId:       patientId,
		StatusId:        mis_dto.ActiveStatusIDs,
	}

	responseBody := mg.sendToMIS(ctx, mis_dto.GetAppointmentsMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return appointmentDTO, err
	}

	for _, misDTO := range response.Data {
		appointmentDTO = misDTO.ToDTO()
		if appointmentDTO.Id == appointmentId {
			break
		}
	}
	return appointmentDTO, nil
}
