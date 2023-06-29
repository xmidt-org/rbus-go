// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		description string
		url         string
		want        Client
	}{
		// Success case
		{
			description: "Valid args",
			url:         "unix://file",
			want: Client{
				network: "unix",
				address: "/file",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			got := New(Config{
				URL: tc.url,
			})

			require.NotNil(got)
			assert.Equal(&tc.want, got)
		})
	}
}
