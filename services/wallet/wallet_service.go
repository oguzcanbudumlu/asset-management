package main

type WalletService struct {
	Repo *WalletRepository
}

func NewWalletService(repo *WalletRepository) *WalletService {
	return &WalletService{Repo: repo}
}

func (s *WalletService) CreateWallet(wallet *Wallet) error {
	return s.Repo.CreateWallet(wallet)
}

func (s *WalletService) DeleteWallet(network, address string) error {
	return s.Repo.DeleteWallet(network, address)
}

func (s *WalletService) GetWallet(network, address string) (*Wallet, error) {
	return s.Repo.GetWallet(network, address)
}
