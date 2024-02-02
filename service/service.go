package service

type Service struct {
	config Config
}

func New(config Config) Service {
	s := Service{config: config}

	return s
}

func (s Service) Start() {

}
