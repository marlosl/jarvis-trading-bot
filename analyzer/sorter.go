package analyzer

import "jarvis-trading-bot/structs"

type SortByInstance []structs.TradingStatusKey

func (t SortByInstance) Len() int {
	return len(t)
}

func (t SortByInstance) Less(i, j int) bool {
	if t[i].UserId < t[j].UserId {
		return true
	}

	if t[i].UserId == t[j].UserId && t[i].BotParameterId < t[j].BotParameterId {
		return true
	}

	if t[i].UserId == t[j].UserId &&
		t[i].BotParameterId == t[j].BotParameterId &&
		t[i].InstanceId < t[j].InstanceId {
		return true
	}

	return false
}

func (t SortByInstance) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
