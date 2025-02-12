package notification

var InClinicVisitReminderTemplate = `Dear %s,

We would like to remind you that you have an appointment on %s, at %s with the doctor %s

Address: %s
https://maps.app.goo.gl/USF1fBXPYA4rWGg88 

If there are any changes, text us via  Whatsapp: https://wa.me/+351915013427 or Telegram: https://t.me/Unitedmedclinic.
Or call us: %s

Best regards,
United Medical Clinic Lisbon

Please do not reply to this message, it was sent automatically.`

var OnlineVisitReminderTemplate = `Dear %s,

We would like to remind you that you have an online consultation on %s, at %s with the doctor %s

Before the consultation you will receive a link to join it to your email.

If there are any changes, text us via  Whatsapp: https://wa.me/+351915013427 or Telegram: https://t.me/Unitedmedclinic.
Or call us: %s

Best regards,
United Medical Clinic Lisbon

Please do not reply to this message, it was sent automatically.`

var HomeVisitReminderTemplate = `Dear %s,

We would like to remind you that you have a home visit on %s, with the doctor %s

If there are any changes, text us via  Whatsapp: https://wa.me/+351915013427 or Telegram: https://t.me/Unitedmedclinic.
Or call us: %s

Best regards,
United Medical Clinic Lisbon

Please do not reply to this message, it was sent automatically.`

var cancelAppointmentTemplate = `Cancellation of visit
%s

Dear %s .
Your visit's on %s, at %s been canceled.

%s.
The doctor is %s.

Phone number for inquiries: %s

Please do not reply to this message, it was sent automatically.`

var CreateHouseCallAppointmentTemplate = `Dear %s,
You have booked a home visit on %s with the doctor %s

Our team will contact you to confirm all the details.

If there are any changes, text us via  Whatsapp: https://wa.me/+351915013427 or Telegram: https://t.me/Unitedmedclinic.
Or call us: %s

Best regards,
United Medical Clinic Lisbon

Please do not reply to this message, it was sent automatically.`

var CreateInClinicAppointmentTemplate = `Dear %s ,
You have an appointment on %s, at %s with the doctor %s

Address: %s
https://maps.app.goo.gl/USF1fBXPYA4rWGg88 

If there are any changes, text us via Whatsapp: https://wa.me/+351915013427 or Telegram: https://t.me/Unitedmedclinic
Or call us: %s

Best regards,
United Medical Clinic Lisbon

Please do not reply to this message, it was sent automatically.`

var CreateOnlineAppointmentTemplate = `Dear %s ,
You have booked an online consultation on %s, at %s with the doctor %s

Before the consultation you will receive a link to join it to your email.

If there are any changes, text us via Whatsapp: https://wa.me/+351915013427 or Telegram: https://t.me/Unitedmedclinic 
Or call us: %s

Best regards,
United Medical Clinic Lisbon

Please do not reply to this message, it was sent automatically.`

type clinicData struct {
	phone   string
	address string
}

var clinicDataMap = map[int]clinicData{
	1: {
		phone:   "+351300502822",
		address: "Av. Defensores de Chaves 73B, 1000-114 Lisboa, Portugal",
	},
}
