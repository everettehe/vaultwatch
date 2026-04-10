// Package notifier provides alert delivery integrations for vaultwatch.
//
// # Telegram Notifier
//
// TelegramNotifier delivers secret expiration alerts via the Telegram Bot API.
//
// # Configuration
//
// The following fields are required in the vaultwatch config:
//
//	notifiers:
//	  telegram:
//	    bot_token: "<your-bot-token>"
//	    chat_id:   "<your-chat-id>"
//
// Obtain a bot token from @BotFather and the chat ID from the
// Telegram getUpdates API or a helper bot.
package notifier
