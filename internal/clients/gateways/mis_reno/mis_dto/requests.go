package mis_dto

const (
	DefaultClinic            = 1
	Source                   = "sorkin_telegram_bot"
	GetUsersMethod           = "getUsers"
	GetSpecialityMethod      = "getProfessions"
	GetScheduleMethod        = "getSchedule"
	GetSchedulePeriodsMethod = "getSchedulePeriods"
	GetAppointmentsMethod    = "getAppointments"
	GetAppointmentsV2Method  = "V2/getAppointments"
	GetPatientMethod         = "getPatient"
	CreatePatientMethod      = "createPatient"
	CreateAppointmentMethod  = "createAppointment"
	CancelAppointmentMethod  = "cancelAppointment"
	ConfirmAppointmentMethod = "confirmAppointment"
)

type CreateAppointmentRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ThirdName string `json:"third_name"`
	BirthDate string `json:"birth_date"`
	Mobile    string `json:"mobile"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	DoctorId  int    `json:"doctor_id"`
	TimeStart string `json:"time_start"` //dd.mm.yyyy hh:mm
	TimeEnd   string `json:"time_end"`   //dd.mm.yyyy hh:mm
	ClinicId  int    `json:"clinic_id"`
	Source    string `json:"source"`
}

type GetScheduleRequest struct {
	ClinicId      int    `json:"clinic_id"`
	DoctorId      int    `json:"user_id"`
	TimeStart     string `json:"time_start"` //dd.mm.yyyy hh:mm
	TimeEnd       string `json:"time_end"`   //dd.mm.yyyy hh:mm
	Room          string `json:"room"`
	Step          int    `json:"step"`
	UseDocAVGTime bool   `json:"use_doctor_avg_time"`
	AllClinics    bool   `json:"all_clinics"`
	ShowBusy      bool   `json:"show_busy"`
	ShowAll       bool   `json:"show_all"`
	ShowPast      bool   `json:"show_past"`
}

type GetSchedulePeriodsRequest struct {
	ClinicId   int    `json:"clinic_id"`
	DoctorId   int    `json:"user_id"`
	TimeStart  string `json:"time_start"` //dd.mm.yyyy hh:mm
	TimeEnd    string `json:"time_end"`   //dd.mm.yyyy hh:mm
	RoleId     int    `json:"role_id"`
	CategoryId int    `json:"category_id"`
	Type       string `json:"type"`
}

type GetSpecialityRequest struct {
	ShowAll        bool `json:"show_all"`
	ShowDeleted    bool `json:"show_deleted"`
	WithoutDoctors bool `json:"without_doctors"`
}

type GetUserRequest struct {
	DoctorId     int    `json:"user_id"`
	SpecialityId int    `json:"profession_id"`
	ClinicId     int    `json:"clinic_id"`
	Role         string `json:"role"`
}

type CancelAppointmentRequest struct {
	AppointmentId int    `json:"appointment_id"`
	Source        string `json:"source"`
	MovedTo       string `json:"moved_to"`
}

type ConfirmAppointmentRequest struct {
	AppointmentId int    `json:"appointment_id"`
	Source        string `json:"source"`
}

type GetAppointmentsRequest struct {
	DateCreatedFrom string `json:"date_created_from"` //dd.mm.yyyy hh:mm
	DateCreatedTo   string `json:"date_created_to"`   //dd.mm.yyyy hh:mm
	PatientId       int    `json:"patient_id"`
	StatusId        string `json:"status_id"`
	AppointmentId   int    `json:"id"`
}

type GetPatientRequest struct {
	Id              int    `json:"id"`
	LastName        string `json:"last_name"`
	FirstName       string `json:"first_name"`
	ThirdName       string `json:"third_name"`
	BirthDate       string `json:"birth_date"`
	BirthDay        string `json:"birth_day"`
	Mobile          string `json:"mobile"`
	Email           string `json:"email"`
	CarNumber       string `json:"car_number"`
	CategoryID      string `json:"category_id"`
	WithDocuments   bool   `json:"with_documents"`
	DateCreatedFrom string `json:"date_created_from"` //dd.mm.yyyy hh:mm
	DateCreatedTo   string `json:"date_created_to"`   //dd.mm.yyyy hh:mm
	DateUpdatedFrom string `json:"date_updated_from"` //dd.mm.yyyy hh:mm
	DateUpdatedTo   string `json:"date_updated_to"`   //dd.mm.yyyy hh:mm
}

type CreatePatientRequest struct {
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	ThirdName string `json:"third_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Mobile    string `json:"mobile"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}
