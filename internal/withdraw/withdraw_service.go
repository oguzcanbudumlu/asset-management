package withdraw

type withdrawService struct {
	withdrawRepository WithdrawRepository
}

type WithdrawService interface {
	Deposit() error
}

func NewWithdrawService(wr WithdrawRepository) WithdrawService {
	return &withdrawService{withdrawRepository: wr}
}

func (s *withdrawService) Deposit() error {
	return nil
}
