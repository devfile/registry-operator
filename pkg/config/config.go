//
// Copyright (c) 2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation

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
