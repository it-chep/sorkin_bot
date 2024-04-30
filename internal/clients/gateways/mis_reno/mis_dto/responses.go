package mis_dto

import "sorkin_bot/internal/domain/entity/appointment"

type BaseResponse struct {
	Error int `json:"error"`
	Data  struct {
		Code             int    `json:"code"`
		ErrorDescription string `json:"desc"`
	} `json:"data"`
}

type MisUserDTO struct {
	ID                     int      `json:"id"`
	Avatar                 string   `json:"avatar"`
	AvatarSmall            string   `json:"avatar_small"`
	Name                   string   `json:"name"`
	BirthDate              string   `json:"birth_date"`
	Gender                 int      `json:"gender"`
	Role                   []string `json:"role"`
	RoleTitles             string   `json:"role_titles"`
	DocumentNumber         string   `json:"document_number"`
	DocumentDate           string   `json:"document_date"`
	Phone                  string   `json:"phone"`
	Email                  string   `json:"email"`
	Contacts               []string `json:"contacts"`
	Profession             []string `json:"profession"`
	ProfessionTitles       string   `json:"profession_titles"`
	SecondProfession       []string `json:"second_profession"`
	SecondProfessionTitles string   `json:"second_profession_titles"`
	Clinic                 []string `json:"clinic"`
	ClinicTitles           string   `json:"clinic_titles"`
	AvgTime                string   `json:"avg_time"`
	HasCompany             bool     `json:"has_company"`
	AvgTimeCompany         string   `json:"avg_time_company"`
	AvgTimeRepeat          string   `json:"avg_time_repeat"`
	AvgTimeRepeatCompany   string   `json:"avg_time_repeat_company"`
	DefaultClinic          string   `json:"default_clinic"`
	DefaultRoom            string   `json:"default_room"`
	IsChildDoctor          bool     `json:"is_child_doctor"`
	IsAdultDoctor          bool     `json:"is_adult_doctor"`
	PatientAgeFrom         int      `json:"patient_age_from"`
	PatientAgeTo           int      `json:"patient_age_to"`
	IsOutside              bool     `json:"is_outside"`
	IsTelemedicine         bool     `json:"is_telemedicine"`
	DateWorkFrom           string   `json:"date_work_from"`
	WorkPeriod             string   `json:"work_period"`
	WorkDegree             string   `json:"work_degree"`
	WorkRank               string   `json:"work_rank"`
	WorkAcademyStatus      string   `json:"work_academy_status"`
	Qualification          string   `json:"qualification"`
	Education              string   `json:"education"`
	EducationCourses       string   `json:"education_courses"`
	Services               []string `json:"services"`
	IsDeleted              bool     `json:"is_deleted"`
}

type GetUsersResponse struct {
	Error int `json:"error"`
	Data  []struct {
		MisUserDTO
	} `json:"data"`
}

func (u MisUserDTO) ToDomain() appointment.Doctor {
	return appointment.NewDoctor(
		u.ID, u.Name, u.Phone, u.Email, u.ProfessionTitles, u.SecondProfessionTitles, u.IsDeleted,
	)
}

type SpecialityDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	DoctorName string `json:"doctor_name"`
	IsDeleted  bool   `json:"is_deleted"`
}

type GetSpecialityResponse struct {
	Error int `json:"error"`
	Data  []struct {
		SpecialityDTO
	} `json:"data"`
}

func (s SpecialityDTO) ToDomain() appointment.Speciality {
	return appointment.NewSpeciality(
		s.ID, s.Name, s.DoctorName, s.IsDeleted,
	)
}

type ScheduleDTO struct {
	ClinicId       int    `json:"clinic_id"`
	Date           string `json:"date"`
	TimeStart      string `json:"time_start"`       //dd.mm.yyyy hh:mm
	TimeStartShort string `json:"time_start_short"` //dd.mm.yyyy hh:mm
	TimeEnd        string `json:"time_end"`         //dd.mm.yyyy hh:mm
	TimeEndShort   string `json:"time_end_short"`   //dd.mm.yyyy hh:mm
	Category       string `json:"category"`
	CategoryId     int    `json:"category_id"`
	Room           string `json:"room"`
	IsBusy         bool   `json:"is_busy"`
	IsPast         bool   `json:"is_past"`
}

type GetScheduleResponse struct {
	Error int `json:"error"`
	Data  []struct {
		ScheduleItem []struct {
			ScheduleDTO
		} // Массив объектов, где ключом является id сотрудника
	} `json:"data"`
}

func (sch ScheduleDTO) ToDomain() {

}

type SchedulePeriodDTO struct {
	Date            string `json:"date"`
	TimeStart       string `json:"time_start"` //dd.mm.yyyy hh:mm
	TimeEnd         string `json:"time_end"`   //dd.mm.yyyy hh:mm
	Type            string `json:"type"`
	ClinicId        int    `json:"clinic_id"`
	DoctorId        int    `json:"user_id"`
	CategoryId      int    `json:"category_id"`
	Room            string `json:"room"`
	WithoutCrossing bool   `json:"without_crossing"`
	DisableInSalary bool   `json:"disable_in_salary"`
}

type GetSchedulePeriodsResponse struct {
	Error int `json:"error"`
	Data  []struct {
		SchedulePeriodDTO
	} `json:"data"`
}

func (schPer SchedulePeriodDTO) ToDomain() {

}

type CreateAppointmentResponse struct {
	Error int `json:"error"`
	Data  struct {
		ID int `json:"id"`
	} `json:"data"`
}

type ConfirmAndCancelAppointmentResponse struct {
	Error int `json:"error"`
	Data  struct {
		True bool `json:"true"`
	} `json:"data"`
}

type AppointmentDTO struct {
	Id               int    `json:"id"`
	TimeStart        string `json:"time_start"` //dd.mm.yyyy hh:mm
	TimeEnd          string `json:"time_end"`   //dd.mm.yyyy hh:mm
	ClinicId         int    `json:"clinic_id"`
	Clinic           string `json:"clinic"`
	DoctorId         int    `json:"doctor_id"`
	Doctor           string `json:"doctor"`
	PatientId        int    `json:"patient_id"`
	PatientName      string `json:"patient_name"`
	PatientBirthDate string `json:"patient_birth_date"`
	PatientGender    string `json:"patient_gender"`
	PatientPhone     string `json:"patient_phone"`
	PatientEmail     string `json:"patient_email"`
	DateCreated      string `json:"date_created"` //dd.mm.yyyy hh:mm
	DateUpdated      string `json:"date_updated"` //dd.mm.yyyy hh:mm
	Status           string `json:"status"`
	StatusId         int    `json:"status_id"`
	ConfirmStatus    string `json:"confirm_status"`
	Source           string `json:"source"`
	MovedTo          string `json:"moved_to"`
	MovedFrom        string `json:"moved_from"`
}

type GetAppointmentsResponse struct {
	Error int `json:"error"`
	Data  []struct {
		AppointmentDTO
	} `json:"data"`
}

func (a AppointmentDTO) ToDomain() appointment.Appointment {
	return appointment.NewAppointment(
		a.Id, a.ClinicId, a.DoctorId, a.PatientId, a.StatusId, a.TimeStart, a.TimeEnd, a.Clinic, a.Doctor,
		a.PatientName, a.PatientBirthDate, a.PatientGender, a.PatientPhone, a.PatientEmail, a.DateCreated, a.DateUpdated,
		a.Status, a.ConfirmStatus, a.Source, a.MovedTo, a.MovedFrom,
	)
}
