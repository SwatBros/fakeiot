package connectors

type Connector interface {
	Send(data any) error
}
