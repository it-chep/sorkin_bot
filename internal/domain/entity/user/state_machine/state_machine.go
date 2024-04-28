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
	userService        user.UserService
	appointmentService appointment.AppointmentService
	bot                telegram.Bot
}

const (
	ChooseLanguage      = "chooseLanguage"
	ChooseSpeciality    = "chooseSpeciality"
	FastAppointment     = "fastAppointment"
	ChooseDoctor        = "chooseDoctor"
	ChooseSchedule      = "chooseSchedule"
	GetPhone            = "getPhone"
	GetName             = "getName"
	CreateAppointment   = "createAppointment"
	MyAppointments      = "myAppointments"
	DetailMyAppointment = "detailMyAppointment"
	CancelAppointment   = "cancelAppointment"
	ChooseAppointment   = "chooseAppointment"
)

func NewUserStateMachine(userService user.UserService) *UserStateMachine {
	machine := &UserStateMachine{
		userService: userService,
	}
	machine.FSM = fsm.NewFSM(
		"",
		fsm.Events{
			{Name: ChooseLanguage, Src: []string{""}, Dst: ChooseLanguage},
			{Name: ChooseSpeciality, Src: []string{"", ChooseLanguage}, Dst: ChooseSpeciality},
			{Name: FastAppointment, Src: []string{""}, Dst: FastAppointment},
			{Name: ChooseDoctor, Src: []string{"", FastAppointment}, Dst: ChooseDoctor},
			{Name: ChooseSchedule, Src: []string{"", ChooseDoctor}, Dst: ChooseSchedule},
			{Name: GetPhone, Src: []string{ChooseSchedule}, Dst: GetPhone},
			{Name: GetName, Src: []string{GetPhone}, Dst: GetName},
			{Name: CreateAppointment, Src: []string{GetName}, Dst: CreateAppointment},
			{Name: MyAppointments, Src: []string{ChooseSpeciality}, Dst: MyAppointments},
			{Name: DetailMyAppointment, Src: []string{MyAppointments}, Dst: DetailMyAppointment},
			{Name: ChooseAppointment, Src: []string{""}, Dst: ChooseAppointment},
			{Name: CancelAppointment, Src: []string{ChooseAppointment}, Dst: CancelAppointment},
		},
		fsm.Callbacks{
			fmt.Sprintf("enter_%s", ChooseLanguage):      enterChooseLanguage,
			fmt.Sprintf("enter_%s", ChooseSpeciality):    enterChooseSpeciality,
			fmt.Sprintf("enter_%s", FastAppointment):     enterFastAppointment,
			fmt.Sprintf("enter_%s", ChooseDoctor):        enterChooseDoctor,
			fmt.Sprintf("enter_%s", ChooseSchedule):      enterChooseSchedule,
			fmt.Sprintf("enter_%s", GetPhone):            enterGetPhone,
			fmt.Sprintf("enter_%s", GetName):             enterGetName,
			fmt.Sprintf("enter_%s", CreateAppointment):   enterCreateAppointment,
			fmt.Sprintf("enter_%s", MyAppointments):      enterMyAppointments,
			fmt.Sprintf("enter_%s", DetailMyAppointment): enterDetailMyAppointment,
			fmt.Sprintf("enter_%s", CancelAppointment):   enterCancelAppointment,
			fmt.Sprintf("enter_%s", ChooseAppointment):   enterChooseAppointment,
		},
	)

	return machine
}

func (machine *UserStateMachine) SetState(user entity.User, from, to string) {
	user.SetState(to)

	// todo проверить чтобы не брались чужие состояния

	machine.FSM.SetState(from)
	err := machine.FSM.Event(context.TODO(), to)
	if err != nil {
		return
	}

}
