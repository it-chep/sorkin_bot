package state_machine

import (
	"context"
	"fmt"
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
	ChooseSchedule      = "chooseSchedule"
	GetPhone            = "getPhone"
	GetBirthDate        = "getBirthDate"
	GetName             = "getName"
	CreateAppointment   = "createAppointment"
	MyAppointments      = "myAppointments"
	DetailMyAppointment = "detailMyAppointment"
	CancelAppointment   = "cancelAppointment"
	ChooseAppointment   = "chooseAppointment"
	MoveAppointment     = "moveAppointment"
	GetDoctorInfo       = "getDoctorInfo"
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
			{Name: ChooseDoctor, Src: []string{Start, ChooseSpeciality, FastAppointment}, Dst: ChooseDoctor},
			{Name: ChooseSchedule, Src: []string{Start, ChooseDoctor}, Dst: ChooseSchedule},
			{Name: GetPhone, Src: []string{ChooseSchedule}, Dst: GetPhone},
			{Name: GetName, Src: []string{GetPhone}, Dst: GetName},
			{Name: GetBirthDate, Src: []string{GetName}, Dst: GetBirthDate},
			{Name: CreateAppointment, Src: []string{GetName, ChooseSchedule, FastAppointment}, Dst: CreateAppointment},
			{Name: MyAppointments, Src: []string{ChooseSpeciality}, Dst: MyAppointments},
			{Name: DetailMyAppointment, Src: []string{ChooseAppointment}, Dst: DetailMyAppointment},
			{Name: ChooseAppointment, Src: []string{Start}, Dst: ChooseAppointment},
			{Name: CancelAppointment, Src: []string{ChooseAppointment}, Dst: CancelAppointment},
			{Name: Start, Src: []string{Start, ChooseLanguage, CreateAppointment, ChooseAppointment, MyAppointments, DetailMyAppointment, MoveAppointment}, Dst: Start},
			{Name: MoveAppointment, Src: []string{DetailMyAppointment}, Dst: MoveAppointment},
			{Name: GetDoctorInfo, Src: []string{GetDoctorInfo, CreateAppointment, DetailMyAppointment, ChooseDoctor, ChooseSchedule}, Dst: GetDoctorInfo},
		},
		fsm.Callbacks{
			fmt.Sprintf("enter_%s", ChooseLanguage):      enterChooseLanguage,
			fmt.Sprintf("enter_%s", ChooseSpeciality):    enterChooseSpeciality,
			fmt.Sprintf("enter_%s", FastAppointment):     enterFastAppointment,
			fmt.Sprintf("enter_%s", ChooseDoctor):        enterChooseDoctor,
			fmt.Sprintf("enter_%s", ChooseSchedule):      enterChooseSchedule,
			fmt.Sprintf("enter_%s", GetPhone):            enterGetPhone,
			fmt.Sprintf("enter_%s", GetName):             enterGetName,
			fmt.Sprintf("enter_%s", GetBirthDate):        enterGetBirthDate,
			fmt.Sprintf("enter_%s", CreateAppointment):   enterCreateAppointment,
			fmt.Sprintf("enter_%s", MyAppointments):      enterMyAppointments,
			fmt.Sprintf("enter_%s", DetailMyAppointment): enterDetailMyAppointment,
			fmt.Sprintf("enter_%s", CancelAppointment):   enterCancelAppointment,
			fmt.Sprintf("enter_%s", ChooseAppointment):   enterChooseAppointment,
			fmt.Sprintf("enter_%s", Start):               enterStart,
			fmt.Sprintf("enter_%s", MoveAppointment):     enterMoveAppointment,
			fmt.Sprintf("enter_%s", GetDoctorInfo):       enterGetDoctorInfo,
		},
	)
	return machine
}

func (machine *UserStateMachine) SetState(user entity.User, to string) {
	ctx := context.Background()
	err := machine.FSM.Event(ctx, to, user, machine.userService)
	if err != nil {
		return
	}
}
