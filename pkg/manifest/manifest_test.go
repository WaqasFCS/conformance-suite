package manifest

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"encoding/json"
	"testing"
)

func TestDiscoveryEndpointsMapToManifestCorrectly(t *testing.T) {
	discoJSON := `
{
	"discoveryModel": {
		"name": "ob-v3.1-ozone",
		"description": "An Open Banking UK discovery template for v3.1 of Accounts and Payments with pre-populated model Bank (Ozone) data.",
		"discoveryVersion": "v0.3.0",
		"tokenAcquisition": "psu",
		"discoveryItems": [{
			"apiSpecification": {
				"name": "Account and Transaction API Specification",
				"url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1",
				"version": "v3.1",
				"schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json",
				"manifest": "file://manifests/ob_3.1__accounts_fca.json"
			},
			"openidConfigurationUri": "https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration",
			"resourceBaseUri": "https://modelobank2018.o3bank.co.uk:4501/open-banking/v3.1/aisp",
			"resourceIds": {
				"ConsentId": "$consent_id"
			},
			"endpoints": [{
					"method": "HEAD",
					"path": "/accounts/{AccountId}"
				},
				{
					"method": "GET",
					"path": "/accounts/{AccountId}"
				},
				{
					"method": "GET",
					"path": "/accounts/{AccountId}/statements/{StatementId}"
				}
			]
		}, {
			"apiSpecification": {
				"name": "Payment Initiation API",
				"url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937754701/Payment+Initiation+API+Specification+-+v3.1",
				"version": "v3.1",
				"schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/payment-initiation-swagger.json",
				"manifest": "file://manifests/ob_3.1__payment_fca.json"
			},
			"openidConfigurationUri": "https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration",
			"resourceBaseUri": "https://modelobank2018.o3bank.co.uk:4501/open-banking/v3.1/",
			"endpoints": [{
					"method": "GET",
					"path": "/domestic-payment-consents"
				}
			]
		}],
		"customTests": [{}]
	}
}
`
	require := test.NewRequire(t)
	mfJSON := `
{
	"scripts": [
        {
			"description": "Minimal data returned for a given account using the ReadAccountsBasic permission, status and headers.",
            "id": "OB-301-ACC-120382",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937623627/Accounts+v3.1#Accountsv3.1-PermissionCodes",
            "detail" : "Checks that the resource differs depending on the permissions (ReadAccountsBasic and ReadAccountsDetail) used to access the resource with additional schema checks on status and headers.",
			"parameters": {
				"accountAccessConsent": "basicAccountAccessConsent",
				"tokenRequestScope": "accounts",
                "accountId": "$consentedAccountId"         
            },
            "uri": "/accounts/{accountId}",
            "uriImplementation": "mandatory",
            "resource": "Account",
            "asserts": ["OB3ACCAssertOnSuccess"],
            "method":"get",
            "schemaCheck": true
        },
        {
			"description": "All data returned for a given account with ReadAccountsDetail permission, status and headers.",
            "id": "OB-301-ACC-352203",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937623627/Accounts+v3.1#Accountsv3.1-PermissionCodes",
            "detail" : "Checks that the resource returns the correct data depending on the permissions ReadAccountsDetail with additional additional schema checks on status and headers.",
			"parameters": {
				"accountAccessConsent": "detailAccountAccessConsent",
				"tokenRequestScope": "accounts",
				"accountId": "$consentedAccountId"
            },
            "uri": "/accounts/{accountId}",
            "uriImplementation": "mandatory",
			"resource": "Account",
            "asserts": ["OB3ACCAssertOnSuccess"],
            "method":"head",
            "schemaCheck": true
        },
		{
			"description": "Domestic Payment consents succeeds with minimal data set with additional schema checks.",
            "id": "OB-301-DOP-206111",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937984109/Domestic+Payments+v3.1#DomesticPaymentsv3.1-POST/domestic-payment-consents",
            "detail" : "Checks that the resource succeeds posting a domestic payment consents with a minimal data set and checks additional schema.",
			"parameters": {
                "tokenRequestScope": "payments",
                "paymentType": "domestic-payment-consents",
                "post" : "minimalDomesticPaymentConsent"    
            },
            "uri": "/domestic-payment-consents",
            "uriImplementation": "mandatory",
            "resource": "DomesticPayment",
            "asserts": ["OB3DOPAssertOnSuccess", "OB3GLOAAssertConsentId"],
            "keepContext": ["OB3GLOAAssertConsentId"],
            "method":"get",
            "schemaCheck": true
        },
		{
			"description": "Domestic Payment consents succeeds with minimal data set with additional schema checks.",
            "id": "OB-301-DOP-2061133",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937984109/Domestic+Payments+v3.1#DomesticPaymentsv3.1-POST/domestic-payment-consents",
            "detail" : "Checks that the resource succeeds posting a domestic payment consents with a minimal data set and checks additional schema.",
			"parameters": {
                "tokenRequestScope": "payments",
                "paymentType": "domestic-payment-consents",
                "post" : "minimalDomesticPaymentConsent"    
            },
            "uri": "/accounts/{accountId}/statements/{statementId}",
            "uriImplementation": "mandatory",
            "resource": "/accounts/{accountId}/statements/{statementId}",
            "asserts": ["OB3DOPAssertOnSuccess", "OB3GLOAAssertConsentId"],
            "keepContext": ["OB3GLOAAssertConsentId"],
            "method":"get",
            "schemaCheck": true
        }
	]
}
`

	var mf Scripts
	err := json.Unmarshal([]byte(mfJSON), &mf)
	require.Nil(err)

	disco, err := discovery.UnmarshalDiscoveryJSON(discoJSON)
	require.Nil(err)

	mpResults := MapDiscoveryEndpointsToManifestTestIDs(disco, mf)

	exp := DiscoveryPathsTestIDs{
		"/accounts/{accountid}": {
			"GET":  {"OB-301-ACC-120382"},
			"HEAD": {"OB-301-ACC-352203"},
		},
		"/domestic-payment-consents": {
			"GET": {"OB-301-DOP-206111"},
		},
		"/accounts/{accountid}/statements/{statementid}": {
			"GET": {"OB-301-DOP-2061133"},
		},
	}

	require.Equal(exp, mpResults)
}

func TestUnMappedManifestItemsReportedCorrectly(t *testing.T) {
	discoJSON := `
{
	"discoveryModel": {
		"name": "ob-v3.1-ozone",
		"description": "An Open Banking UK discovery template for v3.1 of Accounts and Payments with pre-populated model Bank (Ozone) data.",
		"discoveryVersion": "v0.3.0",
		"tokenAcquisition": "psu",
		"discoveryItems": [{
			"apiSpecification": {
				"name": "Account and Transaction API Specification",
				"url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1",
				"version": "v3.1",
				"schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json",
				"manifest": "file://manifests/ob_3.1__accounts_fca.json"
			},
			"openidConfigurationUri": "https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration",
			"resourceBaseUri": "https://modelobank2018.o3bank.co.uk:4501/open-banking/v3.1/aisp",
			"resourceIds": {
				"ConsentId": "$consent_id"
			},
			"endpoints": [{
					"method": "HEAD",
					"path": "/accounts/{AccountId}"
				},
				{
					"method": "GET",
					"path": "/accounts/{AccountId}"
				}
			]
		}, {
			"apiSpecification": {
				"name": "Payment Initiation API",
				"url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937754701/Payment+Initiation+API+Specification+-+v3.1",
				"version": "v3.1",
				"schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/payment-initiation-swagger.json",
				"manifest": "file://manifests/ob_3.1__payment_fca.json"
			},
			"openidConfigurationUri": "https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration",
			"resourceBaseUri": "https://modelobank2018.o3bank.co.uk:4501/open-banking/v3.1/",
			"endpoints": [{
					"method": "GET",
					"path": "/domestic-payment-consents/{ConsentId}/funds-confirmation"
				},
				{
					"method": "POST",
					"path": "/domestic-scheduled-payment-consents"
				}
			]
		}],
		"customTests": [{}]
	}
}
`
	mfJSON := `
{
	"scripts": [
        {
			"description": "Minimal data returned for a given account using the ReadAccountsBasic permission, status and headers.",
            "id": "OB-301-ACC-120382",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937623627/Accounts+v3.1#Accountsv3.1-PermissionCodes",
            "detail" : "Checks that the resource differs depending on the permissions (ReadAccountsBasic and ReadAccountsDetail) used to access the resource with additional schema checks on status and headers.",
			"parameters": {
				"accountAccessConsent": "basicAccountAccessConsent",
				"tokenRequestScope": "accounts",
                "accountId": "$consentedAccountId"         
            },
            "uri": "/accounts/{accountId}",
            "uriImplementation": "mandatory",
            "resource": "Account",
            "asserts": ["OB3ACCAssertOnSuccess"],
            "method":"get",
            "schemaCheck": true
        },
        {
			"description": "All data returned for a given account with ReadAccountsDetail permission, status and headers.",
            "id": "OB-301-ACC-352203",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937623627/Accounts+v3.1#Accountsv3.1-PermissionCodes",
            "detail" : "Checks that the resource returns the correct data depending on the permissions ReadAccountsDetail with additional additional schema checks on status and headers.",
			"parameters": {
				"accountAccessConsent": "detailAccountAccessConsent",
				"tokenRequestScope": "accounts",
				"accountId": "$consentedAccountId"
            },
            "uri": "/accounts/{accountId}",
            "uriImplementation": "mandatory",
			"resource": "Account",
            "asserts": ["OB3ACCAssertOnSuccess"],
            "method":"head",
            "schemaCheck": true
        },
		{
			"description": "",
            "id": "unmapped-test-id",
            "refURI": "",
            "detail" : "",
			"parameters": {},
            "uri": "/FOO-BAR",
            "uriImplementation": "mandatory",
			"resource": "Account",
            "asserts": [],
            "method":"head",
            "schemaCheck": true
        },
		{
			"description": "Domestic Payment consents succeeds with minimal data set with additional schema checks.",
            "id": "OB-301-DOP-206111",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937984109/Domestic+Payments+v3.1#DomesticPaymentsv3.1-POST/domestic-payment-consents",
            "detail" : "Checks that the resource succeeds posting a domestic payment consents with a minimal data set and checks additional schema.",
			"parameters": {
                "tokenRequestScope": "payments",
                "paymentType": "domestic-payment-consents",
                "post" : "minimalDomesticPaymentConsent"    
            },
            "uri": "/domestic-payment-consents",
            "uriImplementation": "mandatory",
            "resource": "DomesticPayment",
            "asserts": ["OB3DOPAssertOnSuccess", "OB3GLOAAssertConsentId"],
            "keepContext": ["OB3GLOAAssertConsentId"],
            "method":"get",
            "schemaCheck": true
        }
	]
}
`
	require := test.NewRequire(t)

	var mf Scripts
	err := json.Unmarshal([]byte(mfJSON), &mf)
	require.Nil(err)

	disco, err := discovery.UnmarshalDiscoveryJSON(discoJSON)
	require.Nil(err)

	mpResults := MapDiscoveryEndpointsToManifestTestIDs(disco, mf)

	unmatched := FindUnmatchedManifestTests(mf, mpResults)

	exp := []string{"unmapped-test-id", "OB-301-DOP-206111"}

	require.Equal(exp, unmatched)
}
