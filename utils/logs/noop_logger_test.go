/*
 * Copyright (C) 2020-2021 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package logs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoopLogger(t *testing.T) {
	loggers, err := NewNoopLogger("Test")
	require.Nil(t, err)
	_testLog(t, loggers)
}
