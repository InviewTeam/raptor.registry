package tests

type Publisher struct {
}

func (p *Publisher) Connect() error {
	return nil
}

func (p *Publisher) Close() error {
	return nil
}

func (p *Publisher) DeclareQueue(queue string) error {
	return nil
}

func (p *Publisher) Send(data []byte, queue string) error {
	return nil
}
