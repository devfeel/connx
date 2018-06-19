package connx

const(
	DefaultVersion = "1.0"
	DefaultRequestCommand = "request"
	DefaultResponseCommand = "response"
)

// Message message define
type Message struct {
	Version string
	Command string
	Data    interface{}
}

func NewMessage(ver, cmd string, data interface{}) *Message{
	return &Message{
		Version:ver,
		Command:cmd,
		Data:data,
	}
}

func RequestMessage(data interface{}) *Message{
	return &Message{
		Version:DefaultVersion,
		Command:DefaultRequestCommand,
		Data:data,
	}
}

func ResponseMessage(data interface{}) *Message{
	return &Message{
		Version:DefaultVersion,
		Command:DefaultResponseCommand,
		Data:data,
	}
}