package config

// TelegramConfig holds configuration for the Telegram notifier.
type TelegramConfig struct {
	// BotToken is the Telegram Bot API token issued by @BotFather.
	BotToken string `yaml:"bot_token" mapstructure:"bot_token"`

	// ChatID is the target chat, group, or channel ID.
	// For channels use the format "-100xxxxxxxxxx".
	ChatID string `yaml:"chat_id" mapstructure:"chat_id"`
}
