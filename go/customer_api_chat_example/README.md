# Chat.io example Customer API integration

This is a simple example of chat.io integration that creates a customer (by using Customer Accounts) on a license and use it to start a new chat. 

# What does it exactly do?

* First of all the integration must handle the installation process on a license. To do it it binds to local port to parse `HTTP GET` request that is a result of agent-sso redirect to receive a `OAuth code` and exchange it to `access_token` of an agent (see https://www.chat.io/docs/authorization/#public-server-side-apps).

* After getting the `access_token` it gets `access_token` information from agent SSO (`license` number in this case only).

* Once successfully obtaining `access_token` and `license` number, it sets up the necessary resources in the license and creates Customer.

* Customer is created by Customer Accounts service. To obtain a new customer it does `HTTP POST` request to `https://accounts.chat.io/customer/` using (link to customer-sso)