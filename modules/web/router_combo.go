// Copyright (C) Kumo inc. and its affiliates.
// Author: Jeff.li lijippy@163.com
// All rights reserved.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package web

// Combo represents a tiny group routes with same pattern
type Combo struct {
	r       *Router
	pattern string
	h       []any
}

// Get delegates Get method
func (c *Combo) Get(h ...any) *Combo {
	c.r.Get(c.pattern, append(c.h, h...)...)
	return c
}

// Post delegates Post method
func (c *Combo) Post(h ...any) *Combo {
	c.r.Post(c.pattern, append(c.h, h...)...)
	return c
}

// Delete delegates Delete method
func (c *Combo) Delete(h ...any) *Combo {
	c.r.Delete(c.pattern, append(c.h, h...)...)
	return c
}

// Put delegates Put method
func (c *Combo) Put(h ...any) *Combo {
	c.r.Put(c.pattern, append(c.h, h...)...)
	return c
}

// Patch delegates Patch method
func (c *Combo) Patch(h ...any) *Combo {
	c.r.Patch(c.pattern, append(c.h, h...)...)
	return c
}
