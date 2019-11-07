/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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
