package state_machine

import (
	"context"
	"github.com/looplab/fsm"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/services/appointment"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/pkg/client/telegram"
)

type UserStateMachine struct {
	FSM                *fsm.FSM
	userService        userService
	appointmentService appointment.AppointmentService
	bot                telegram.Bot
}

const (
	Start               = ""
	ChooseLanguage      = "chooseLanguage"
	ChooseSpeciality    = "chooseSpeciality"
	FastAppointment     = "fastAppointment"
	ChooseDoctor        = "chooseDoctor"
	ChooseCalendar      = "chooseCalendar"
	ChooseSchedule      = "chooseSchedule"
	GetPhone            = "getPhone"
	GetBirthDate        = "getBirthDate"
	GetName             = "getName"
	CreateAppointment   = "createAppointment"
	DetailMyAppointment = "detailMyAppointment"
	CancelAppointment   = "cancelAppointment"
	ChooseMyAppointment = "chooseMyAppointment"
	ChooseAppointment   = "chooseAppointment"
	GetDoctorInfo       = "getDoctorInfo"
	ClinicAppointment   = "clinicAppointment"
	HomeAppointment     = "homeAppointment"
	OnlineAppointment   = "onlineAppointment"
	Pediatrician        = "pediatrician"
	Therapist           = "therapist"
	SetAddress          = "setAddress"
)

func NewUserStateMachine(userService user.UserService) *UserStateMachine {
	machine := &UserStateMachine{
		userService: userService,
	}
	machine.FSM = fsm.NewFSM(
		Start,
		fsm.Events{
			{Name: ChooseLanguage, Src: []string{Start}, Dst: ChooseLanguage},
			{Name: ChooseSpeciality, Src: []string{Start, ChooseLanguage}, Dst: ChooseSpeciality},
			{Name: FastAppointment, Src: []string{Start}, Dst: FastAppointment},
			{Name: ChooseDoctor, Src: []string{DetailMyAppointment, Start, ChooseSpeciality, FastAppointment, CreateAppointment}, Dst: ChooseDoctor},
			{Name: ChooseSchedule, Src: []string{Start, ChooseCalendar}, Dst: ChooseSchedule},
			{Name: GetPhone, Src: []string{ChooseSchedule}, Dst: GetPhone},
			{Name: GetName, Src: []string{GetPhone}, Dst: GetName},
			{Name: GetBirthDate, Src: []string{GetName}, Dst: GetBirthDate},
			{Name: CreateAppointment, Src: []string{GetName, ChooseSchedule, FastAppointment, GetDoctorInfo}, Dst: CreateAppointment},
			{Name: DetailMyAppointment, Src: []string{Start, GetDoctorInfo, ChooseMyAppointment}, Dst: DetailMyAppointment},
			{Name: ChooseMyAppointment, Src: []string{Start, GetDoctorInfo}, Dst: ChooseMyAppointment},
			{Name: CancelAppointment, Src: []string{ChooseMyAppointment}, Dst: CancelAppointment},
			{Name: Start, Src: []string{DetailMyAppointment, Start, ChooseLanguage, ChooseSpeciality, FastAppointment,
				ChooseDoctor, ChooseSchedule, CreateAppointment, CancelAppointment, ChooseMyAppointment, GetDoctorInfo,
			}, Dst: Start},
			{Name: ChooseCalendar, Src: []string{ChooseDoctor, Start}, Dst: ChooseCalendar},
			{Name: GetDoctorInfo, Src: []string{GetDoctorInfo, CreateAppointment, ChooseMyAppointment, DetailMyAppointment, ChooseDoctor, ChooseSchedule}, Dst: GetDoctorInfo},
			{Name: ChooseAppointment, Src: []string{Start}, Dst: ChooseAppointment},
		},
		fsm.Callbacks{},
	)
	return machine
}

func (machine *UserStateMachine) SetState(user entity.User, to string) {
	ctx := context.Background()
	err := changeState(ctx, to, user, machine.userService)
	if err != nil {
		return
	}
}

func changeState(ctx context.Context, to string, user entity.User, userService userService) error {
	_, err := userService.ChangeState(ctx, user.GetTgId(), to)
	if err != nil {
		return err
	}
	return nil
}
