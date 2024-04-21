package state_machine

import (
	"context"
	"fmt"
	"github.com/looplab/fsm"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UserStateMachine struct {
	To  string
	FSM *fsm.FSM
}

func NewUserStateMachine(to string) *UserStateMachine {
	machine := &UserStateMachine{
		To: to,
	}
	machine.FSM = fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "chooseLanguage", Src: []string{"start"}, Dst: "chooseLanguage"},
			{Name: "chooseSpeciallity", Src: []string{"chooseLanguage"}, Dst: "chooseSpeciallity"},
			{Name: "chooseClosestDoctor", Src: []string{"chooseSpeciallity"}, Dst: "chooseClosestDoctor"},
			{Name: "fastAppointment", Src: []string{"chooseClosestDoctor"}, Dst: "fastAppointment"},
			{Name: "chooseDoctor", Src: []string{"fastAppointment"}, Dst: "chooseDoctor"},
			{Name: "chooseSchedule", Src: []string{"chooseDoctor"}, Dst: "chooseSchedule"},
			{Name: "getPhone", Src: []string{"chooseSchedule"}, Dst: "getPhone"},
			{Name: "getName", Src: []string{"getPhone"}, Dst: "getName"},
			{Name: "createAppointment", Src: []string{"getName"}, Dst: "createAppointment"},
			{Name: "chooseMyAppointments", Src: []string{"chooseSpeciallity"}, Dst: "myAppointments"},
			{Name: "detailMyAppointment", Src: []string{"myAppointments"}, Dst: "detailMyAppointment"},
		},
		fsm.Callbacks{
			"enter_chooseLanguage": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Entered state: chooseLanguage")
			},
			"enter_chooseSpeciallity": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: chooseSpeciallity")
			},
			"enter_chooseClosestDoctor": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: chooseClosestDoctor")
			},
			"enter_fastAppointment": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: fastAppointment")
			},
			"enter_chooseDoctor": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: chooseDoctor")
			},
			"enter_chooseSchedule": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: chooseSchedule")
			},
			"enter_getPhone": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: getPhone")
			},
			"enter_getName": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: getName")
			},
			"enter_createAppointment": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: createAppointment")
			},
			"enter_myAppointments": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: myAppointments")
			},
			"enter_detailMyAppointment": func(_ context.Context, e *fsm.Event) {
				fmt.Println("Entered state: detailMyAppointment")
			},
		},
	)

	return machine
}

func (machine *UserStateMachine) enterState(e *fsm.Event) {
	fmt.Printf("The door to %s is %s\n", machine.To, e.Dst)
}

func (machine *UserStateMachine) SetState(user entity.User, to string) {
	user.SetState(to)
}
