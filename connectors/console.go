package connectors

import "fmt"

type ConsoleConnector struct {
}

func (cc *ConsoleConnector) Send(data any) error {
	fmt.Println(data)
	return nil
}
