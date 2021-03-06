/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package auth

import (
	"github.com/apache/incubator-servicecomb-service-center/pkg/chain"
	"github.com/apache/incubator-servicecomb-service-center/pkg/log"
	"github.com/apache/incubator-servicecomb-service-center/pkg/rest"
	scerr "github.com/apache/incubator-servicecomb-service-center/server/error"
	"github.com/apache/incubator-servicecomb-service-center/server/plugin"
	"github.com/apache/incubator-servicecomb-service-center/server/rest/controller"
	"net/http"
)

type AuthRequest struct {
}

func (h *AuthRequest) Handle(i *chain.Invocation) {
	r := i.Context().Value(rest.CTX_REQUEST).(*http.Request)
	err := plugin.Plugins().Auth().Identify(r)
	if err == nil {
		i.Next()
		return
	}

	log.Errorf(err, "authenticate request failed, %s %s", r.Method, r.RequestURI)

	w := i.Context().Value(rest.CTX_RESPONSE).(http.ResponseWriter)
	controller.WriteError(w, scerr.ErrUnauthorized, err.Error())

	i.Fail(nil)
}

func RegisterHandlers() {
	chain.RegisterHandler(rest.SERVER_CHAIN_NAME, &AuthRequest{})
}
