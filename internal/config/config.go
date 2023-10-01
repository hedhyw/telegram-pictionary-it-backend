package config

import "time"

type Config struct {
	DebugEnabled          bool          `env:"DEBUG_ENABLED" envDefault:"true"`
	ServerAddress         string        `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:8081"`
	ServerTimeout         time.Duration `env:"SERVER_TIMEOUT" envDefault:"10s"`
	ServerReadBufferSize  int           `env:"SERVER_READ_BUFFER_SIZE" envDefault:"1024"`
	ServerWriteBufferSize int           `env:"SERVER_WRITE_BUFFER_SIZE" envDefault:"1024"`
	TelegramBotToken      string        `env:"TELEGRAM_BOT_TOKEN" envDefault:"change_this_secret"`
	GameRoundTimeout      time.Duration `env:"GAME_ROUND_TIMEOUT" envDefault:"2m"`
	WorkersCount          int           `env:"SERVER_WORKERS_COUNT" envDefault:"1024"`
	AnimationURL          string        `json:"ANIMATION_URL" envDefault:"https://64.media.tumblr.com/dcfe9790135b1ff3aa08a06df31419d7/tumblr_pcw8kaJLYO1qhy6c9o1_r1_500.gifv"`
}

func (cfg Config) Sanitized() Config {
	const mask = "***"

	if cfg.TelegramBotToken != "" {
		cfg.TelegramBotToken = mask
	}

	return cfg
}
