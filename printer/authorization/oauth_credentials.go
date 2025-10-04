package authorization

type OAuthCredential struct {
	clientKey       string
	clientSecret    string
	photoServiceURL string
	callbackURL     string
}

var oauthCredential *OAuthCredential

// OAuth client credentials (registered with Photo Service)
func NewOAuthCredential() {
	oauthCredential = &OAuthCredential{
		clientKey:       "dpf43f3p2l4k3l03",
		clientSecret:    "kd94hf93k423kf44",
		photoServiceURL: "http://172.27.19.159:8081",
		callbackURL:     "http://172.27.19.159:8080/callback",
	}
}

func GetCredential() *OAuthCredential {
	return oauthCredential
}

func (o *OAuthCredential) GetClientKey() string {
	return o.clientKey
}

func (o *OAuthCredential) GetClientSecret() string {
	return o.clientSecret
}

func (o *OAuthCredential) GetPhotoServiceURL() string {
	return o.photoServiceURL
}
func (o *OAuthCredential) GetCallbackURL() string {
	return o.callbackURL
}
