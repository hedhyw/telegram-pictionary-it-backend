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
	AnimationURL          string        `env:"ANIMATION_URL" envDefault:"https://pusheen.com/wp-content/uploads/2018/11/tumblr_pcw8kaJLYO1qhy6c9o1_r1_1280.gif"`
	ServerCheckOrigin     string        `env:"SERVER_CHECK_ORIGIN" envDefault:"http://localhost:3000"`
}

func (cfg Config) Sanitized() Config {
	const mask = "***"

	if cfg.TelegramBotToken != "" {
		cfg.TelegramBotToken = mask
	}

	return cfg
}
