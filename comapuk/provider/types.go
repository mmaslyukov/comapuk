package provider

type Config = map[string]string

const (
	ComEmail    = "email"
	ComPassword = "password"

	IntLogin    = "login"
	IntPassword = "password"

	ZkhViewState = "javax.faces.ViewState"
	ZkhLogin     = "login"
	ZkhPassword  = "passw"
	ZkhCertPath  = "cert_path"
	ZkhSource    = "javax.faces.source"
	ZkhExecute   = "javax.faces.partial.execute"
	ZkhRenderer  = "javax.faces.partial.render"
	ZkhFormId    = "calculationsHeadForm:j_idt46"

	PageAccount      = "page_account"
	PageSubscription = "page_subscription"
	PageLogin        = "page_login"
	PageCalc         = "page_calc"
)
