package transfer

type transferService struct {
	transferRepository TransferRepository
}

type TransferService interface {
	Deposit() error
}

func NewTransferService(tr TransferRepository) TransferService {
	return &transferService{transferRepository: tr}
}

func (s *transferService) Deposit() error {
	return nil
}
