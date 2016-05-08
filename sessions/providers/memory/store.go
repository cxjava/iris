// Copyright (c) 2016, Gerasimos Maropoulos
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//	  this list of conditions and the following disclaimer
//    in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse
//    or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER AND CONTRIBUTOR, GERASIMOS MAROPOULOS
// BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package memory

import (
	"time"

	"github.com/kataras/iris/sessions"
)

// Store the memory store, contains the session id and the values
type Store struct {
	sid              string
	lastAccessedTime time.Time
	values           map[interface{}]interface{} // here is the real memory store
}

var _ sessions.IStore = &Store{}

// GetAll returns all values
func (s *Store) GetAll() map[interface{}]interface{} {
	return s.values
}

// VisitAll loop each one entry and calls the callback function func(key,value)
func (s *Store) VisitAll(cb func(k interface{}, v interface{})) {
	for key := range s.values {
		cb(key, s.values[key])
	}
}

// Get returns the value of an entry by its key
func (s *Store) Get(key interface{}) interface{} {
	provider.Update(s.sid)

	if value, found := s.values[key]; found {
		return value
	}

	return nil
}

// GetString same as Get but returns as string, if nil then returns an empty string
func (s *Store) GetString(key interface{}) string {
	if value := s.Get(key); value != nil {
		return value.(string)
	}

	return ""
}

// GetInt same as Get but returns as int, if nil then returns -1
func (s *Store) GetInt(key interface{}) int {
	if value := s.Get(key); value != nil {
		return value.(int)
	}

	return -1
}

// Set fills the session with an entry, it receives a key and a value
// returns an error, which is always nil
func (s *Store) Set(key interface{}, value interface{}) error {
	s.values[key] = value
	provider.Update(s.sid)
	return nil
}

// Delete removes an entry by its key
// returns an error, which is always nil
func (s *Store) Delete(key interface{}) error {
	delete(s.values, key)
	provider.Update(s.sid)
	return nil
}

// Clear removes all entries
// returns an error, which is always nil
func (s *Store) Clear() error {
	for key := range s.values {
		delete(s.values, key)
	}
	provider.Update(s.sid)
	return nil
}

// ID returns the session id
func (s *Store) ID() string {
	return s.sid
}

// LastAccessedTime returns the last time this session has been used
func (s *Store) LastAccessedTime() time.Time {
	return s.lastAccessedTime
}

// SetLastAccessedTime updates the last accessed time
func (s *Store) SetLastAccessedTime(lastacc time.Time) {
	s.lastAccessedTime = lastacc
}

// Destroy does nothing here, to destroy the session use the manager's .Destroy func
func (s *Store) Destroy() {
	// nothing
}
