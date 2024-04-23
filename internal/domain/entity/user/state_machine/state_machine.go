package state_machine

import (
	"context"
	"github.com/looplab/fsm"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UserStateMachine struct {
	FSM *fsm.FSM
}

func NewUserStateMachine() *UserStateMachine {
	machine := &UserStateMachine{}
	machine.FSM = fsm.NewFSM(
		"",
		fsm.Events{
			{Name: "chooseLanguage", Src: []string{""}, Dst: "chooseLanguage"},
			{Name: "chooseSpeciality", Src: []string{"", "chooseLanguage"}, Dst: "chooseSpeciality"},
			{Name: "chooseClosestDoctor", Src: []string{"", "chooseSpeciallity"}, Dst: "chooseClosestDoctor"},
			{Name: "fastAppointment", Src: []string{"", "chooseClosestDoctor"}, Dst: "fastAppointment"},
			{Name: "chooseDoctor", Src: []string{"", "fastAppointment"}, Dst: "chooseDoctor"},
			{Name: "chooseSchedule", Src: []string{"", "chooseDoctor"}, Dst: "chooseSchedule"},
			{Name: "getPhone", Src: []string{"chooseSchedule"}, Dst: "getPhone"},
			{Name: "getName", Src: []string{"getPhone"}, Dst: "getName"},
			{Name: "createAppointment", Src: []string{"getName"}, Dst: "createAppointment"},
			{Name: "chooseMyAppointments", Src: []string{"chooseSpeciallity"}, Dst: "myAppointments"},
			{Name: "detailMyAppointment", Src: []string{"myAppointments"}, Dst: "detailMyAppointment"},
			{Name: "cancelAppointment", Src: []string{""}, Dst: "cancelAppointment"},
		},
		fsm.Callbacks{
			"enter_chooseLanguage":      enterChooseLanguage,
			"enter_chooseSpeciality":    enterChooseSpeciality,
			"enter_chooseClosestDoctor": enterChooseClosestDoctor,
			"enter_fastAppointment":     enterFastAppointment,
			"enter_chooseDoctor":        enterChooseDoctor,
			"enter_chooseSchedule":      enterChooseSchedule,
			"enter_getPhone":            enterGetPhone,
			"enter_getName":             enterGetName,
			"enter_createAppointment":   enterCreateAppointment,
			"enter_myAppointments":      enterMyAppointments,
			"enter_detailMyAppointment": enterDetailMyAppointment,
			"enter_cancelAppointment":   enterCancelAppointment,
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
