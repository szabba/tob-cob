// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package interaction

import (
	"github.com/szabba/assert/v2"
	"github.com/szabba/assert/v2/assertions/theval"
)

// A Target is something that might provide interactions to perform against it.
//
// Any zero-length slice means the target has no available interactions right now.
type Target interface {
	Interactions() []Interaction
}

// An Interaction is an instantaneous state change that can be applied to a target.
type Interaction interface {
	Interact()
}

// A DefaultTarget wraps another target and
type DefaultTarget struct {
	tgt Target
	def []Interaction
}

var _ Target = DefaultTarget{}

// NewDefaultTarget wraps tgt, creating a target that provides def as the sole possible interaction when tgt reports none.
func NewDefaultTarget(tgt Target, def Interaction) DefaultTarget {

	t := DefaultTarget{
		tgt,
		[]Interaction{def},
	}

	t.verify()
	return t
}

func (t DefaultTarget) Interactions() []Interaction {

	t.verify()

	ints := t.tgt.Interactions()
	if len(ints) == 0 {
		return t.def
	}

	return ints
}

func (t DefaultTarget) verify() {

	assert.UsingPanic().
		That(theval.NotZero(t)).
		That(theval.NotZero(t.tgt)).
		That(theval.NotZero(t.def))
}
