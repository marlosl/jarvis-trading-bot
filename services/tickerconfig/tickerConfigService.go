package tickerconfig

import "fmt"

func GetTickerConfig(ticker string) (*TickerConfigItem, error) {
	repo, err := NewTickerConfigRepository()
	if err != nil {
		return nil, err
	}

	item, err := repo.GetItem(ticker)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func SaveTickerConfig(item *TickerConfigItem) error {
	repo, err := NewTickerConfigRepository()
	if err != nil {
		return err
	}

	err = repo.SaveItem(item)
	if err != nil {
		return err
	}

	return nil
}

func CreateTickerConfigIfNotExists(ticker string) error {
	item, err := GetTickerConfig(ticker)
	if err != nil {
		return err
	}

	if item == nil || item.Ticker == "" {
		item = &TickerConfigItem{
			Ticker: ticker,
		}
		err = SaveTickerConfig(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteTickerConfig(ticker string) error {
	repo, err := NewTickerConfigRepository()
	if err != nil {
		return err
	}

	err = repo.DeleteItem(ticker)
	if err != nil {
		return err
	}

	return nil
}

func ListTickers() []string {
	tickers := make([]string, 0)
	repo, err := NewTickerConfigRepository()

	if err != nil {
		fmt.Printf("Can't create ticker config repo: %v", err)
		return tickers
	}

	items, err := repo.GetItems()
	if err != nil {
		fmt.Printf("Can't get ticker config items: %v", err)
		return tickers
	}

	for _, item := range items {
		tickers = append(tickers, item.Ticker)
	}

	return tickers
}
