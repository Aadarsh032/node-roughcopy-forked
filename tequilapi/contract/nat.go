/*
 * Copyright (C) 2020 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package contract

import (
	"github.com/mysteriumnetwork/node/nat"
)

// NATStatusDTO gives information about NAT traversal success or failure
// swagger:model NATStatusDTO
type NATStatusDTO struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// NATTypeDTO gives information about NAT type in terms of traversal capabilities
// swagger:model NATTypeDTO
type NATTypeDTO struct {
	Type  nat.NATType `json:"type"`
	Error string      `json:"error,omitempty"`
}

// Nat nat related information
type Nat struct {
	Status NATStatusDTO `json:"status"`
}
