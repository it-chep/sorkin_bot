package mis_dto

import (
	"encoding/json"
	"sorkin_bot/internal/clients/gateways/dto"
	"strconv"
)

type BaseResponse struct {
	Error int `json:"error"`
	Data  struct {
		Code             int    `json:"code"`
		ErrorDescription string `json:"desc"`
	} `json:"data"`
}

type MisUser struct {
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
	Profession             []string `json:"profession"`
	ProfessionTitles       string   `json:"profession_titles"`
	SecondProfession       []string `json:"second_profession"`
	SecondProfessionTitles string   `json:"second_profession_titles"`
	Clinic                 []string `json:"clinic"`
	ClinicTitles           *string  `json:"clinic_titles"`
	AvgTime                int      `json:"avg_time"`
	HasCompany             bool     `json:"has_company"`
	DefaultClinic          int      `json:"default_clinic"`
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
		MisUser
	} `json:"data"`
}

func (u MisUser) ToDTO() dto.DoctorDTO {
	professionIds := make([]int, len(u.SecondProfession))
	if len(u.SecondProfession) > 0 {
		for _, profession := range u.SecondProfession {
			professionId, _ := strconv.Atoi(profession)
			professionIds = append(professionIds, professionId)
		}
	}
	return dto.NewDoctorDTO(
		u.ID, u.Name, u.Phone, u.Email, u.ProfessionTitles, u.SecondProfessionTitles, u.IsDeleted, professionIds,
	)
}

type MisSpeciality struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	DoctorName string `json:"doctor_name"`
	IsDeleted  bool   `json:"is_deleted"`
}

type GetSpecialityResponse struct {
	Error int `json:"error"`
	Data  []struct {
		MisSpeciality
	} `json:"data"`
}

func (s MisSpeciality) ToDTO() dto.SpecialityDTO {
	return dto.NewSpecialityDTO(
		s.ID, s.Name, s.DoctorName, s.IsDeleted,
	)
}

type MisSchedule struct {
	ClinicId       int    `json:"clinic_id"`
	DoctorId       int    `json:"user_id"`
	Date           string `json:"date"`
	TimeStart      string `json:"time_start"`       //dd.mm.yyyy hh:mm
	TimeStartShort string `json:"time_start_short"` //dd.mm.yyyy hh:mm
	TimeEnd        string `json:"time_end"`         //dd.mm.yyyy hh:mm
	TimeEndShort   string `json:"time_end_short"`   //dd.mm.yyyy hh:mm
	Category       string `json:"category"`
	CategoryId     int    `json:"category_id"`
	Profession     string `json:"profession"`
	Room           string `json:"room"`
	User           string `json:"user"`
	IsBusy         bool   `json:"is_busy"`
	IsPast         bool   `json:"is_past"`
}

type GetScheduleResponse struct {
	Error int                      `json:"error"`
	Data  map[string][]MisSchedule `json:"data"`
}

func (sch MisSchedule) ToDTO() dto.ScheduleDTO {
	return dto.NewScheduleDTO(
		sch.ClinicId, sch.DoctorId, sch.CategoryId, sch.Date, sch.TimeStart, sch.TimeStartShort, sch.TimeEnd,
		sch.TimeEndShort, sch.Category, sch.Profession, sch.Room, sch.User, sch.IsBusy, sch.IsPast,
	)
}

type MisSchedulePeriod struct {
	Date      string `json:"date"`
	TimeStart string `json:"time_start"` //dd.mm.yyyy hh:mm
	TimeEnd   string `json:"time_end"`   //dd.mm.yyyy hh:mm
	DoctorId  int    `json:"user_id"`
}

type GetSchedulePeriodsResponse struct {
	Error int `json:"error"`
	Data  []struct {
		MisSchedulePeriod
	} `json:"data"`
}

func (schPer MisSchedulePeriod) ToDTO() dto.SchedulePeriodDTO {
	return dto.NewSchedulePeriodDTO(
		schPer.Date, schPer.TimeStart, schPer.TimeEnd, schPer.DoctorId,
	)
}

type CreateAppointmentResponse struct {
	Error int    `json:"error"`
	Data  string `json:"data"`
}

type ConfirmAndCancelAppointmentResponse struct {
	Error int  `json:"error"`
	Data  bool `json:"data"`
}

type MisAppointment struct {
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
	ConfirmStatus    int    `json:"confirm_status"`
	Source           string `json:"source"`
	MovedTo          int    `json:"moved_to"`
	MovedFrom        int    `json:"moved_from"`
	IsTelemedicine   bool   `json:"is_telemedicine"`
	IsOutside        bool   `json:"is_outside"`
}

type GetAppointmentsResponse struct {
	Error int `json:"error"`
	Data  []struct {
		MisAppointment
	} `json:"data"`
}

func (a MisAppointment) ToDTO() dto.AppointmentDTO {
	outside := 0
	if a.IsOutside {
		outside = 1
	}

	telemedicine := 0
	if a.IsTelemedicine {
		telemedicine = 1
	}

	return dto.NewAppointmentDTO(
		a.Id, a.ClinicId, a.DoctorId, a.PatientId, a.StatusId, a.MovedTo, a.MovedFrom, a.ConfirmStatus, outside,
		telemedicine, a.TimeStart, a.TimeEnd, a.Clinic, a.Doctor, a.PatientName, a.PatientBirthDate,
		a.PatientGender, a.PatientPhone, a.PatientEmail, a.DateCreated, a.DateUpdated, a.Status, a.Source,
	)
}

type MisCreatePatientData struct {
	PatientID int    `json:"patient_id"`
	Number    int    `json:"number"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	ThirdName string `json:"third_name"`
	BirthDate string `json:"birth_date"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
	Mobile    string `json:"mobile"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type CreatePatientResponse struct {
	Error int `json:"error"`
	Data  struct {
		MisCreatePatientData
	} `json:"data"`
}

type MisPatientAddress struct {
	City        string `json:"city"`
	Street      string `json:"street"`
	House       string `json:"house"`
	Building    string `json:"building"`
	Flat        string `json:"flat"`
	District    string `json:"district"`
	ZipCode     string `json:"zip_code"`
	Metro       string `json:"metro"`
	FullAddress string `json:"fullAddress"`
}

type GetPatientData struct {
	Number         int               `json:"number"`
	LastName       string            `json:"last_name"`
	FirstName      string            `json:"first_name"`
	ThirdName      string            `json:"third_name"`
	BirthDate      string            `json:"birth_date"`
	BirthPlace     string            `json:"birth_place"`
	Age            string            `json:"age"`
	Gender         string            `json:"gender"`
	Mobile         string            `json:"mobile"`
	Phone          string            `json:"phone"`
	Email          string            `json:"email"`
	CarNumber      string            `json:"car_number"`
	ExternalID     string            `json:"external_id"`
	ExternalNumber string            `json:"external_number"`
	JobPlace       string            `json:"job_place"`
	JobProfession  string            `json:"job_profession"`
	JobPosition    string            `json:"job_position"`
	Factory        string            `json:"factory"`
	Unit           string            `json:"unit"`
	Workshop       string            `json:"workshop"`
	Area           string            `json:"area"`
	ServiceNumber  string            `json:"service_number"`
	JobAttitude    string            `json:"job_attitude"`
	SendSMS        string            `json:"send_sms"`
	SendEmail      string            `json:"send_email"`
	SendSMSLab     string            `json:"send_sms_lab"`
	SendEmailLab   string            `json:"send_email_lab"`
	DateCreated    string            `json:"date_created"`
	DateUpdated    string            `json:"date_updated"`
	TimeCreated    string            `json:"time_created"`
	TimeUpdated    string            `json:"time_updated"`
	CategoryIDs    string            `json:"category_ids"`
	ParentID       string            `json:"parent_id"`
	Groups         string            `json:"groups"`
	HasAccount     string            `json:"has_account"`
	AdvChannelID   string            `json:"adv_channel_id"`
	Desc           string            `json:"desc"`
	Address        MisPatientAddress `json:"address"`
	IsDeleted      bool              `json:"is_deleted"`
	DateDeleted    string            `json:"date_deleted"`
	PatientID      int               `json:"patient_id"`
	Documents      json.RawMessage   `json:"documents"`
}

type MisGetPatientResponse struct {
	Error int `json:"error"`
	Data  struct {
		GetPatientData
	} `json:"data"`
}

type MisGetPatientsResponse struct {
	Error int `json:"error"`
	Data  []struct {
		GetPatientData
	} `json:"data"`
}

func (p GetPatientData) ToDTO() dto.CreatedPatientDTO {
	return dto.NewCreatePatientDTO(
		p.PatientID, p.Number, p.Age, p.Gender, p.Mobile, p.Email, p.Phone,
	)
}
