package exchangeconfig

func GetExchangeConfig(ticker string, exchange string) (*ExchangeConfigItem, error) {
	repo, err := NewExchangeConfigRepository()
	if err != nil {
		return nil, err
	}

	item, err := repo.GetItem(ticker, exchange)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func SaveExchangeConfig(item *ExchangeConfigItem) error {
	repo, err := NewExchangeConfigRepository()
	if err != nil {
		return err
	}

	err = repo.SaveItem(item)
	if err != nil {
		return err
	}

	return nil
}
