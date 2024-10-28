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

func (s *WalletService) GetWallets() ([]Wallet, error) {
	return s.Repo.GetWallets()
}

func (s *WalletService) DeleteWallet(address string) error {
	return s.Repo.DeleteWallet(address)
}
