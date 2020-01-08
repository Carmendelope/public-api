/*
 * Copyright 2020 Nalej
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

package common

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	DefaultTimeout = time.Minute
	UserID         = "userid"
)

// GetContext returns a context with a default timeout for internal communications. Notice that the context does not
// have any security related information attached to it.
func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

func GetContextWithUser(userId string) (context.Context, context.CancelFunc) {
	md := metadata.New(map[string]string{UserID: userId})
	log.Debug().Interface("md", md).Msg("metadata has been created")
	baseContext, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	return metadata.NewOutgoingContext(baseContext, md), cancel
}
