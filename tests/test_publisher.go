package tests

type Publisher struct {
}

func (p *Publisher) Connect() error {
	return nil
}

func (p *Publisher) Close() error {
	return nil
}

func (p *Publisher) Send(data []byte) error {
	return nil
}
