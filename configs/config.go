package configs

import "time"

type BackOffConfig struct {
	baseDelay  time.Duration
	multiplier float64
	jitter     float64
	maxDelay   time.Duration
	retries    uint
}

var (
	baseDelay  = 1 * time.Second
	multiplier = 1.1
	jitter     = 0.2
	maxDelay   = 120 * time.Second
	retries    = 10
)

func NewBackOffConfig(defaultCfg bool) *BackOffConfig {
	if defaultCfg {
		return &BackOffConfig{
			baseDelay:  baseDelay,
			multiplier: multiplier,
			jitter:     jitter,
			maxDelay:   maxDelay,
			retries:    uint(retries),
		}
	}

	return &BackOffConfig{}
}

func (b *BackOffConfig) SetBaseDelay(baseDelay time.Duration) {
	b.baseDelay = baseDelay
}

func (b *BackOffConfig) GetBaseDelay() time.Duration {
	return b.baseDelay
}

func (b *BackOffConfig) SetMultiplier(multiplier float64) {
	b.multiplier = multiplier
}

func (b *BackOffConfig) GetMultiplier() float64 {
	return b.multiplier
}

func (b *BackOffConfig) SetJitter(jitter float64) {
	b.jitter = jitter
}

func (b *BackOffConfig) GetJitter() float64 {
	return b.jitter
}

func (b *BackOffConfig) SetMaxDelay(maxDelay time.Duration) {
	b.maxDelay = maxDelay
}

func (b *BackOffConfig) GetMaxDelay() time.Duration {
	return b.maxDelay
}

func (b *BackOffConfig) SetRetries(retries uint) {
	b.retries = retries
}

func (b *BackOffConfig) GetRetries() uint {
	return b.retries
}
