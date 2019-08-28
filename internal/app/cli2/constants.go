package cli2

import "time"

// LoginAddress with the login api address
const LoginAddress = "login_address"
// NalejAddress with the managment cluster address
const NalejAddress = "nalej_address"
// CACert with the certificate to be use to authenticate the API
const CACert = "cacert"


// OutputFormat with the output format of the results of the commands.
const OutputFormat = "output"
// OutputLabelLength with the maximum of the labels to be shown when table format is selected
const OutputLabelLength = "label_length"

// DefaultTimeout with the maximum time awaiting for the API to respond.
const DefaultTimeout = time.Minute

// AuthHeader with the name of the header used to send authorization information
const AuthHeader = "Authorization"