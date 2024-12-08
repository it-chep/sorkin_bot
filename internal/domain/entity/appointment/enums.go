package appointment

type AppointmentType string

const (
	OnlineAppointment AppointmentType = "online_appointment"
	ClinicAppointment AppointmentType = "clinic_appointment"
	HomeAppointment   AppointmentType = "home_appointment"
)
