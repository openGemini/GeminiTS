/*
Copyright 2023 Huawei Cloud Computing Technologies Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httpd

import (
	"net/http"

	meta2 "github.com/openGemini/openGemini/open_src/influx/meta"
)

func (h *Handler) serveQueryLogByCursor(w http.ResponseWriter, r *http.Request, user meta2.User) {}

func (h *Handler) serveConsumeLogs(w http.ResponseWriter, r *http.Request, user meta2.User) {}

func (h *Handler) serveConsumeCursorTime(w http.ResponseWriter, r *http.Request, user meta2.User) {}

func (h *Handler) serveGetConsumeCursors(w http.ResponseWriter, r *http.Request, user meta2.User) {}

func (h *Handler) serveContextQueryLog(w http.ResponseWriter, r *http.Request, user meta2.User) {}

func (h *Handler) serveGetCursor(w http.ResponseWriter, r *http.Request, user meta2.User) {}

func (h *Handler) servePullLog(w http.ResponseWriter, r *http.Request, user meta2.User) {}
