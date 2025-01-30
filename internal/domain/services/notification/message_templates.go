package notification

var visitReminderTemplate = `Visit Reminder
%s

Dear, %s .
Please be reminded that you have an appointment on %s, at %s.

%s.
The doctor is %s.

Address: %s
Phone number for inquiries: %s

https://maps.app.goo.gl/USF1fBXPYA4rWGg88 

Please do not reply to this message, it was sent automatically.`

var cancelAppointmentTemplate = `Cancellation of visit
%s

Dear, %s .
Your visit's on %s, at %s been canceled.

%s.
The doctor is %s.

Address: %s
Phone number for inquiries: %s

https://maps.app.goo.gl/USF1fBXPYA4rWGg88 

Please do not reply to this message, it was sent automatically.`

var createAppointmentTemplate = `Dear, %s .
You have an appointment on %s, at %s.

%s.
The doctor is %s.

Address: %s
Phone number for inquiries: %s

https://maps.app.goo.gl/USF1fBXPYA4rWGg88 

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
