// Package notifier provides alert delivery integrations for VaultWatch.
//
// # Twilio SMS Notifier
//
// The Twilio notifier sends SMS alerts via the Twilio REST API.
//
// Required configuration:
//
//	notifiers:
//	  twilio:
//	    account_sid: "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
//	    auth_token: "your_auth_token"
//	    from: "+15005550006"
//	    to: "+15005550010"
//
// The notifier sends a plain-text SMS message containing the secret path,
// expiration status, and days remaining.
package notifier
