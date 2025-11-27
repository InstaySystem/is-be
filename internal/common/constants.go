package common

const (
	ExchangeEmail       = "email.send"
	QueueNameAuthEmail  = "email.send.auth"
	RoutingKeyAuthEmail = "email.send.auth"

	ExchangeFile         = "file.action"
	QueueNameDeleteFile  = "file.action.delete"
	RoutingKeyDeleteFile = "file.action.delete"

	ExchangeNotification          = "notification.send"
	QueueNameServiceNotification  = "notification.send.service"
	RoutingKeyServiceNotification = "notification.send.service"
	QueueNameRequestNotification  = "notification.send.request"
	RoutingKeyRequestNotification = "notification.send.request"

	RoleAdmin            = "admin"
	RoleAdminDisplayName = "Quản trị viên"
	RoleStaff            = "staff"
	RoleStaffDisplayName = "Nhân viên"
)
