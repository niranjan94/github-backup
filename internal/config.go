package internal

type Config struct {
	Owners        []string `json:"owners"`
	Token         string   `json:"token"`
	Directory     string   `json:"directory"`
	Concurrency   int      `json:"concurrency"`
	PeriodSeconds int      `json:"period_seconds"`
}
