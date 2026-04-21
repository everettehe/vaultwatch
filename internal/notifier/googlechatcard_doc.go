// Package notifier provides notification implementations for VaultWatch.
//
// # Google Chat Card Notifier
//
// The GoogleChatCardNotifier sends rich card-formatted messages to a Google
// Chat webhook using the Cards v2 API format. Unlike the basic GoogleChatNotifier
// which sends plain text messages, this notifier renders structured cards with
// headers, sections, and colour-coded severity indicators.
//
// Configuration example:
//
//	notifiers:
//	  googlechatcard:
//	    webhook_url: "https://chat.googleapis.com/v1/spaces/.../messages?key=..."
//
// The card header colour is determined by the secret's expiration status:
//   - RED    — secret is already expired
//   - YELLOW — secret is expiring within the warning threshold
//   - GREEN  — secret is healthy (informational only)
package notifier
