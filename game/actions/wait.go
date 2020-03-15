// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package actions

import (
	"time"
)

// Wait returns an action that lasts waitTime but does nothing.
func Wait(waitTime time.Duration) Action {
	return &_Wait{waitTime}
}

type _Wait struct {
	toEnd time.Duration
}

func (w *_Wait) Run(atMost time.Duration) Status {
	if atMost < w.toEnd {
		w.toEnd -= atMost
		return Paused()
	}
	return Done(atMost - w.toEnd)
}
