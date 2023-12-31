// Copyright 2017 Gerasimos Maropoulos, ΓΜ. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"github.com/brucejrc/iris/cache/cfg"
	"github.com/brucejrc/iris/cache/client/rule"
	"github.com/brucejrc/iris/cache/ruleset"
	"github.com/brucejrc/iris/context"
)

// DefaultRuleSet is a list of the default pre-cache validators
// which exists in ALL handlers, local and remote.
var DefaultRuleSet = rule.Chained(
	// #1 A shared cache MUST NOT use a cached response to a request with an
	// Authorization header field
	rule.HeaderClaim(ruleset.AuthorizationRule),
	// #2 "must-revalidate" and/or
	// "s-maxage" response directives are not allowed to be served stale
	// (Section 4.2.4) by shared caches.  In particular, a response with
	// either "max-age=0, must-revalidate" or "s-maxage=0" cannot be used to
	// satisfy a subsequent request without revalidating it on the origin
	// server.
	rule.HeaderClaim(ruleset.MustRevalidateRule),
	rule.HeaderClaim(ruleset.ZeroMaxAgeRule),
	// #3 custom No-Cache header used inside this library
	// for BOTH request and response (after get-cache action)
	rule.Header(ruleset.NoCacheRule, ruleset.NoCacheRule),
)

// NoCache called when a particular handler is not valid for cache.
// If this function called inside a handler then the handler is not cached
// even if it's surrounded with the Cache/CacheFunc wrappers.
func NoCache(ctx context.Context) {
	ctx.Header(cfg.NoCacheHeader, "true")
}
