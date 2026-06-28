/**
 * Copyright 2024-2026 Su Yang (soulteary)
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

package fn_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type mockShutdownApp struct {
	called bool
}

func (m *mockShutdownApp) Shutdown() error {
	m.called = true
	return nil
}

func TestMockShutdownApp(t *testing.T) {
	mock := &mockShutdownApp{}
	require.NoError(t, mock.Shutdown())
	require.True(t, mock.called)
}
