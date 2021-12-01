package queue

//RService interface
type Service interface {
	Send(message string)
	Close()
}
