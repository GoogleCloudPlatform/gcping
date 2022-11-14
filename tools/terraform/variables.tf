// Copyright 2022 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

variable "image" {
  type = string
  default = ""
}

variable "repository" {
  type = string
  default = "ping-b5e9c300f5e9cdafa118e623a88e6b97"
}
variable "project" {
  type    = string
  default = "gcping-devrel"
}

variable "domain" {
  type    = string
  default = "gcping.com"
}

variable "domain_alias_flag" {
  type    = bool
  default = true
}

variable "domain_alias" {
  type    = string
  default = "gcpping.com" // two p's
}

variable "release_bucket" {
  type    = string
  default = "gcping-release"
}
