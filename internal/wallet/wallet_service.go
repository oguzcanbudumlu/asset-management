package wallet

type walletService struct {
	repo WalletRepository
}
type WalletService interface {
	CreateWallet(wallet *Wallet) error
	GetWallet(network, address string) (*Wallet, error)
	DeleteWallet(network, address string) error
}

func NewWalletService(repo WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) CreateWallet(wallet *Wallet) error {
	return s.repo.CreateWallet(wallet)
}

func (s *walletService) DeleteWallet(network, address string) error {
	return s.repo.DeleteWallet(network, address)
}

func (s *walletService) GetWallet(network, address string) (*Wallet, error) {
	return s.repo.GetWallet(network, address)
}
