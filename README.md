# Go Library to assist for access to Microsoft Graph API

This library wraps the Graph API by Microsoft v1.0.  The library can authenticate as a service (NewKeyClient) or
as a delegate for a user (NewClient).  The NewClient function will launch a browser for the OAuth2 authentication 
flow.  This will obviously not work with a web application and slightly different API will be needed in such
a case.

## Application Registration
An Application must be created in Azure Active Directory in order for OAuth2 to work.  This registration
will get you a client ID and client Key to be used.
See https://docs.microsoft.com/en-us/graph/auth-register-app-v2 for steps and more information.

To get the Tenant ID, see https://docs.microsoft.com/en-us/onedrive/find-your-office-365-tenant-id.
The Tenant ID is a UUID like 8978cef5-80eb-4282-a783-642044e5f373

## Examples
There are some examples in the examples folder to assist with learning the library.