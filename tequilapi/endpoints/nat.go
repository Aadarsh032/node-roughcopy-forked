/*
 * Copyright (C) 2019 The "MysteriumNetwork/node" Authors.
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

package endpoints

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mysteriumnetwork/node/nat"
	"github.com/mysteriumnetwork/node/tequilapi/contract"
	"github.com/mysteriumnetwork/node/tequilapi/utils"
)

// NATEndpoint struct represents endpoints about NAT traversal
type NATEndpoint struct {
	stateProvider stateProvider
	natProber     natProber
}

type natProber interface {
	Probe(context.Context) (nat.NATType, error)
}

// NewNATEndpoint creates and returns nat endpoint
func NewNATEndpoint(stateProvider stateProvider, natProber natProber) *NATEndpoint {
	return &NATEndpoint{
		stateProvider: stateProvider,
		natProber:     natProber,
	}
}

// NATStatus provides NAT configuration info
// swagger:operation GET /nat/status NAT NATStatusDTO
// ---
// summary: Shows NAT status
// description: NAT status returns the last known NAT traversal status
// responses:
//   200:
//     description: NAT status ("not_finished"/"successful"/"failed") and optionally error if status is "failed"
//     schema:
//       "$ref": "#/definitions/NATStatusDTO"
func (ne *NATEndpoint) NATStatus(resp http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	utils.WriteAsJSON(ne.stateProvider.GetState().NATStatus, resp)
}

// NATStatusV2 provides NAT configuration info
// swagger:operation GET /v2/nat/status NAT
// ---
// summary: Shows NAT status
// description: NAT status returns the last known NAT traversal status
// responses:
//   200:
//     description: NAT status ("passed"/"failed"/"pending)
//     schema:
//       "$ref": "#/definitions/NATStatusDTO"
func (ne *NATEndpoint) NATStatusV2(resp http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	utils.WriteAsJSON(ne.stateProvider.GetState().Nat.Status, resp)
}

// NATType provides NAT type in terms of traversal capabilities
// swagger:operation GET /nat/type NAT NATTypeDTO
// ---
// summary: Shows NAT type in terms of traversal capabilities.
// description: Returns NAT type. May produce invalid result while VPN connection is established
// responses:
//   200:
//     description: NAT type
//     schema:
//       "$ref": "#/definitions/NATTypeDTO"
//   500:
//     description: Internal server error
//     schema:
//       "$ref": "#/definitions/ErrorMessageDTO"
func (ne *NATEndpoint) NATType(resp http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res, err := ne.natProber.Probe(req.Context())
	if err != nil {
		utils.SendError(resp, err, http.StatusInternalServerError)
		return
	}
	utils.WriteAsJSON(contract.NATTypeDTO{
		Type:  res,
		Error: "",
	}, resp)
}

// AddRoutesForNAT adds nat routes to given router
func AddRoutesForNAT(router *httprouter.Router, stateProvider stateProvider, natProber natProber) {
	natEndpoint := NewNATEndpoint(stateProvider, natProber)

	router.GET("/nat/status", natEndpoint.NATStatus)
	router.GET("/nat/type", natEndpoint.NATType)

	router.GET("/v2/nat/status", natEndpoint.NATStatusV2)
}
