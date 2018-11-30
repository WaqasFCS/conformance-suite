package model

var (
	// Get /accounts example json response from ozone
	specificationStaticData = []byte(
		`[
			{
			  "identifier": "account-transaction-v3.0",
			  "name": "Account and Transaction API Specification",
			  "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
			  "version": "v3.0",
			  "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
			},
			{
			  "identifier": "payment-initiation-v3.0",
			  "name": "Payment Initiation API",
			  "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/645367011/Payment+Initiation+API+Specification+-+v3.0",
			  "version": "v3.0",
			  "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/payment-initiation-swagger.json"
			},
			{
			  "identifier": "confirmation-funds-v3.0",
			  "name": "Confirmation of Funds API Specification",
			  "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/645203467/Confirmation+of+Funds+API+Specification+-+v3.0",
			  "version": "v3.0",
			  "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/confirmation-funds-swagger.json"
			},
			{
			  "identifier": "event-notification-aspsp-v3.0",
			  "name": "Event Notification API Specification - ASPSP Endpoints",
			  "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/645367055/Event+Notification+API+Specification+-+v3.0",
			  "version": "v3.0",
			  "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/callback-urls-swagger.yaml"
			},
			{
			  "identifier": "event-notification-tpp-v3.0",
			  "name": "Event Notification API Specification - TPP Endpoints",
			  "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/645367055/Event+Notification+API+Specification+-+v3.0",
			  "version": "v3.0",
			  "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/event-notifications-swagger.yaml"
			}
		  ]
		  `)
)