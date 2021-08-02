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

package port

import (
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Pool hands out ports for service use
type Pool struct {
	start, capacity int
	rand            *rand.Rand
}

// ServicePortSupplier provides port needed to run a service on
type ServicePortSupplier interface {
	Acquire() (Port, error)
	AcquireMultiple(n int) (ports []Port, err error)
}

// NewFixedRangePool creates a fixed size pool from port.Range
func NewFixedRangePool(r Range) *Pool {
	return &Pool{
		start:    r.Start,
		capacity: r.Capacity(),
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Acquire returns an unused port in pool's range
func (pool *Pool) Acquire() (port Port, err error) {
	p := pool.randomPort()
	available, err := available(p)
	if err != nil {
		return 0, errors.Wrap(err, "could not acquire port")
	}
	if !available {
		p, err = pool.seekAvailablePort()
	}
	log.Info().Err(err).Msgf("Supplying port %d", p)
	return Port(p), errors.Wrap(err, "could not acquire port")
}

func (pool *Pool) randomPort() int {
	return pool.start + pool.rand.Intn(pool.capacity)
}

func (pool *Pool) seekAvailablePort() (int, error) {
	for i := 0; i < pool.capacity; i++ {
		p := pool.start + i
		available, err := available(p)
		if available || err != nil {
			return p, err
		}
	}
	return 0, errors.New("port pool is exhausted")
}

// AcquireMultiple returns n unused ports from pool's range.
func (pool *Pool) AcquireMultiple(n int) (ports []Port, err error) {
	for i := 0; i < n; i++ {
		p, err := pool.Acquire()
		if err != nil {
			return ports, err
		}

		ports = append(ports, p)
	}

	return ports, nil
}
