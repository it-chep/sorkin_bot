package state_machine

import (
	"context"
	"github.com/looplab/fsm"
	entity "sorkin_bot/internal/domain/entity/user"
)

// todo переписать, а то какая-то херня
// запуск событий откуда куда

func changeState(ctx context.Context, e *fsm.Event) {
	user := e.Args[0].(entity.User)
	service := e.Args[1].(userService)
	_, err := service.ChangeState(ctx, user.GetTgId(), e.Dst)
	if err != nil {
		return
	}
}

func enterChooseLanguage(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterChooseSpeciality(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterChooseClosestDoctor(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterFastAppointment(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterChooseDoctor(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterChooseSchedule(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterGetPhone(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterGetName(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterCreateAppointment(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterMyAppointments(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterDetailMyAppointment(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterCancelAppointment(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterChooseAppointment(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterGetBirthDate(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}

func enterStart(ctx context.Context, e *fsm.Event) {
	if len(e.Args) > 1 {
		changeState(ctx, e)
	}
}
