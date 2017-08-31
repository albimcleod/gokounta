# gokounta #

A simple go library for Kounta Applications

## Usage ##

**Create a client**
code := "" //  Code from site
clientId := "" // Client ID provided by Kounta
clientSecret := "" //  Client Secret provided by Kounta
redirectUrl := "" // RedirectUrl as defined in Kounta
v := gokounta.NewClient(code, clientId, clientSecret, redirectUrl)

**Get an access token**
at, rt, err := v.AccessToken()

**Get a new access token, from your refresh token**
at, rt, err := v.RefreshToken(rt)

**Get Company**
company, err := v.GetCompany(at)

**Get all Categories**
categories, err := v.GetCategories(at, company.ID)