package mis_reno

import (
	"context"
	"fmt"
	"math/rand"
	"sorkin_bot/internal/clients/gateways/dto"
	"sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (mg *MisRenoGateway) GetSchedulesByDoctorId(ctx context.Context, doctorId int, timeStart, timeEnd string) (schedulesMap map[int][]dto.ScheduleDTO, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedulesByDoctorId"
	var (
		request = mis_dto.GetScheduleRequest{
			DoctorId:   doctorId,
			ShowBusy:   false,
			ShowPast:   false,
			ShowAll:    false,
			AllClinics: false,
			Step:       ScheduleStepHour,
		}
		response  mis_dto.GetScheduleResponse
		schedules []dto.ScheduleDTO
	)

	// Если мы хотим получить по конкретному дню, то добавляем параметр timeStart
	if timeStart != "" && timeEnd != "" {
		request.TimeStart = timeStart
		request.TimeEnd = timeEnd
	}

	schedulesMap = make(map[int][]dto.ScheduleDTO)

	responseBody := mg.sendToMIS(ctx, mis_dto.GetScheduleMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return schedulesMap, err
	}

	// Если мы берем расписание у всех докторов
	for strResponseDoctorId, doctorSchedule := range response.Data {
		for _, schedule := range doctorSchedule {
			schedules = append(schedules, schedule.ToDTO())
		}
		if doctorId == 0 {
			responseDoctorId, _ := strconv.Atoi(strResponseDoctorId)
			schedulesMap[responseDoctorId] = schedules
		} else {
			schedulesMap[doctorId] = schedules
		}
	}

	return schedulesMap, nil
}

func listShuffle[T any](array []T) {
	rand.Seed(time.Now().UnixNano())
	for i := len(array) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

func (mg *MisRenoGateway) GetSchedulesManyDoctors(ctx context.Context, doctorIds []int, timeStart, timeEnd string) (schedules []dto.ScheduleDTO, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedulesByDoctorId"

	var (
		request = mis_dto.GetScheduleManyDoctorsRequest{
			DoctorIds:  mg.doctorIdsToString(doctorIds),
			ShowBusy:   false,
			ShowPast:   false,
			ShowAll:    false,
			AllClinics: false,
			Step:       ScheduleStepHour,
		}
		response mis_dto.GetScheduleResponse
	)

	// Если мы хотим получить по конкретному дню, то добавляем параметр timeStart
	if timeStart != "" && timeEnd != "" {
		request.TimeStart = timeStart
		request.TimeEnd = timeEnd
	}
	responseBody := mg.sendToMIS(ctx, mis_dto.GetScheduleMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return schedules, err
	}

	// Парсим ответ
	for _, doctorSchedule := range response.Data {
		for _, schedule := range doctorSchedule {
			schedules = append(schedules, schedule.ToDTO())
		}
	}

	// Перемешиваем массив с расписаниями докторов, потому что в апиху к нам приходят упорядоченные расписания по каждому врачу
	listShuffle(schedules)
	doctorsTime := make(map[int]dto.ScheduleDTO, 20)
	shuffleSchedules := make([]dto.ScheduleDTO, 0, len(schedules))
	keys := make([]int, 0, len(doctorsTime))

	//Формируем уникальные слоты с расписанием
	for _, schedule := range schedules {
		if _, ok := doctorsTime[schedule.GetIntHourStart()]; ok {
			continue
		}
		keys = append(keys, schedule.GetIntHourStart())
		doctorsTime[schedule.GetIntHourStart()] = schedule
	}

	//Сортируем рандомный список слотов, потому что в map рандомный вывод
	sort.Ints(keys)
	for _, key := range keys {
		shuffleSchedules = append(shuffleSchedules, doctorsTime[key])
	}

	return shuffleSchedules, nil
}

func (mg *MisRenoGateway) doctorIdsToString(doctorIds []int) string {
	var builder strings.Builder
	for i, id := range doctorIds {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(strconv.Itoa(id))
	}

	return builder.String()
}

// GetAvailableDoctorIdsFromSchedulePeriods метод получает id докторов по специальностям и день работы и возвращает только тех, докторов кто работает в этот день
func (mg *MisRenoGateway) GetAvailableDoctorIdsFromSchedulePeriods(ctx context.Context, doctorIdsBySpeciality []int, timeStart, timeEnd string) (availableDoctorIds []int, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedulesByDoctorId"
	var (
		request = mis_dto.GetSchedulePeriodsRequest{
			DoctorId: mg.doctorIdsToString(doctorIdsBySpeciality),
			Type:     1,
		}
		response mis_dto.GetSchedulePeriodsResponse
	)

	doctorIdSet := map[int]struct{}{}

	// Если мы хотим получить по конкретному дню, то добавляем параметр timeStart
	if timeStart != "" && timeEnd != "" {
		request.TimeStart = timeStart
		request.TimeEnd = timeEnd
	}

	availableDoctorIds = make([]int, 0, len(doctorIdSet))
	responseBody := mg.sendToMIS(ctx, mis_dto.GetSchedulePeriodsMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return availableDoctorIds, err
	}

	for _, schedulePeriodsItem := range response.Data {
		doctorIdSet[schedulePeriodsItem.DoctorId] = struct{}{}
	}

	for doctorId, _ := range doctorIdSet {
		availableDoctorIds = append(availableDoctorIds, doctorId)
	}

	return availableDoctorIds, nil
}

// GetSchedulePeriodsByDoctorId метод получает рабочие дни доктора
func (mg *MisRenoGateway) GetSchedulePeriodsByDoctorId(ctx context.Context, doctorId int, timeStart, timeEnd string) (schedulePeriods []dto.SchedulePeriodDTO, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedulesByDoctorId"
	var (
		request = mis_dto.GetSchedulePeriodsRequest{
			DoctorId:  fmt.Sprintf("%d", doctorId),
			TimeStart: timeStart,
			TimeEnd:   timeEnd,
			Type:      1,
		}
		response mis_dto.GetSchedulePeriodsResponse
	)

	responseBody := mg.sendToMIS(ctx, mis_dto.GetSchedulePeriodsMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return nil, err
	}

	// Парсим ответ
	for _, schedulePeriod := range response.Data {
		schedulePeriods = append(schedulePeriods, schedulePeriod.ToDTO())
	}
	return schedulePeriods, nil
}
