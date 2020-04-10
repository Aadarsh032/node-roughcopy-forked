/*
 * Copyright (C) 2018 The "MysteriumNetwork/node" Authors.
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

package session

import (
	"errors"
	"testing"

	"github.com/mysteriumnetwork/node/core/connection"
	"github.com/mysteriumnetwork/node/identity"
	"github.com/mysteriumnetwork/node/market"
	node_session "github.com/mysteriumnetwork/node/session"
	"github.com/stretchr/testify/assert"
)

var (
	stubLocation = &StubServiceDefinition{}

	errMock     = errors.New("error")
	sessionID   = node_session.ID("sessionID")
	consumerID  = identity.FromAddress("consumerID")
	providerID  = identity.FromAddress("providerID")
	serviceType = "serviceType"

	mockSession = connection.Status{
		SessionID:  sessionID,
		ConsumerID: consumerID,
		Proposal: market.ServiceProposal{
			ServiceDefinition: stubLocation,
			ServiceType:       serviceType,
			ProviderID:        providerID.Address,
		},
	}
)

func TestSessionStorageGetAll(t *testing.T) {
	storer := &StubSessionStorer{}
	storage := NewSessionStorage(storer)
	sessions, err := storage.GetAll()
	assert.Nil(t, err)
	assert.True(t, storer.GetAllCalled)
	assert.Len(t, sessions, 0)
}

func TestSessionStorageGetAllReturnsError(t *testing.T) {
	storer := &StubSessionStorer{
		GetAllError: errMock,
	}
	storage := NewSessionStorage(storer)
	sessions, err := storage.GetAll()
	assert.NotNil(t, err)
	assert.True(t, storer.GetAllCalled)
	assert.Nil(t, sessions)
}

func TestSessionStorageConsumeEventEndedOK(t *testing.T) {
	storer := &StubSessionStorer{}

	storage := NewSessionStorage(storer)
	storage.consumeSessionEvent(connection.AppEventConnectionSession{
		Status:      connection.SessionCreatedStatus,
		SessionInfo: mockSession,
	})
	storage.consumeSessionEvent(connection.AppEventConnectionSession{
		Status:      connection.SessionEndedStatus,
		SessionInfo: mockSession,
	})
	assert.True(t, storer.UpdateCalled)
}

func TestSessionStorageConsumeEventConnectedOK(t *testing.T) {
	storer := &StubSessionStorer{}

	storage := NewSessionStorage(storer)
	storage.consumeSessionEvent(connection.AppEventConnectionSession{
		Status:      connection.SessionCreatedStatus,
		SessionInfo: mockSession,
	})
	assert.True(t, storer.SaveCalled)
}

// StubSessionStorer allows us to get all sessions, save and update them
type StubSessionStorer struct {
	SaveError    error
	SaveCalled   bool
	UpdateError  error
	UpdateCalled bool
	GetAllCalled bool
	GetAllError  error
}

func (sss *StubSessionStorer) Store(from string, object interface{}) error {
	sss.SaveCalled = true
	return sss.SaveError
}

func (sss *StubSessionStorer) Update(from string, object interface{}) error {
	sss.UpdateCalled = true
	return sss.UpdateError
}

func (sss *StubSessionStorer) GetAllFrom(from string, array interface{}) error {
	sss.GetAllCalled = true
	return sss.GetAllError
}

type StubServiceDefinition struct{}

func (fs *StubServiceDefinition) GetLocation() market.Location { return market.Location{} }
