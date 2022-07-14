/*
Copyright 2020-2022 Red Hat, Inc.

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

package config

// ControllerCfg logic borrowed from https://github.com/devfile/devworkspace-operator/blob/master/pkg/config/config.go
var ControllerCfg ControllerConfig

type ControllerConfig struct {
	isOpenShift bool
}

func (c *ControllerConfig) IsOpenShift() bool {
	return c.isOpenShift
}

func (c *ControllerConfig) SetIsOpenShift(isOpenShift bool) {
	c.isOpenShift = isOpenShift
}
