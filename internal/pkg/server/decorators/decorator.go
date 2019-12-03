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
 */

package decorators

import "github.com/nalej/derrors"

// Decorator is an interface that should be implemented by any particular decorator that we want to include
// To include new Decorator:
// 1. Create a decorator that implements Decorator interface:
// To add structures to which decorators can apply:**
//

type Decorator interface {
	// Apply method called for the execution of the decorator. Receives a list of interfaces and returns the
	// list once the decorator is executed.
	Apply(elements []interface{}) ([]interface{}, derrors.Error)
	// Validate is the function in charge of validating if the decorator can be executed on the field requested
	// or any other necessary validation. The structure to which the decorator is to be applied is passed
	// because many times the validation to be carried out will depend on this structure
	Validate(result interface{}) derrors.Error
}
