package models

type ApiSpec struct {
	HttpVerb          string
	Path              string
	CurrentApiName    string
	CurrentApiComment string
	Calls             []ApiCall
}
