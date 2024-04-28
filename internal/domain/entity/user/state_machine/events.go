package state_machine

import (
	"context"
	"fmt"
	"github.com/looplab/fsm"
)

// запуск событий откуда куда

func enterChooseLanguage(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: chooseLanguage")
}

func enterChooseSpeciality(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: chooseSpeciallity")
}

func enterChooseClosestDoctor(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: chooseClosestDoctor")
}

func enterFastAppointment(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: fastAppointment")
}

func enterChooseDoctor(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: chooseDoctor")
}

func enterChooseSchedule(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: chooseSchedule")
}

func enterGetPhone(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: getPhone")
}

func enterGetName(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: getName")
}

func enterCreateAppointment(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: createAppointment")
}

func enterMyAppointments(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: myAppointments")
}

func enterDetailMyAppointment(ctx context.Context, e *fsm.Event) {
	fmt.Println("Entered state: detailMyAppointment")
}

func enterCancelAppointment(ctx context.Context, e *fsm.Event) {
	fmt.Println("I have canceled your appointment")
}

func enterChooseAppointment(ctx context.Context, e *fsm.Event) {
	fmt.Println("Please choose appointment")
}
