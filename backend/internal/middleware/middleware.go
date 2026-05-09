package middleware

import (
	"fmt"
	"net/http"
)

// Tag is a string identifier for named middleware. Handlers declare which
// Tags they require; the router resolves them to concrete handler functions.
type Tag string

const (
	TagAuthenticated    Tag = "authenticated"
	TagAnyAuthenticated Tag = "any_authenticated"
)

var registry = make(map[Tag]func(http.Handler) http.Handler)

// Register associates a Tag with its middleware. Panics on duplicate registration.
func Register(tag Tag, mw func(http.Handler) http.Handler) {
	if _, exists := registry[tag]; exists {
		panic(fmt.Sprintf("middleware: tag %q already registered", tag))
	}
	registry[tag] = mw
}

// Resolve returns middleware functions for the given tags in order.
// Panics if any tag is unregistered — an unresolvable tag is a startup
// misconfiguration, not a runtime error.
func Resolve(tags []Tag) []func(http.Handler) http.Handler {
	mws := make([]func(http.Handler) http.Handler, 0, len(tags))
	for _, tag := range tags {
		mw, ok := registry[tag]
		if !ok {
			panic(fmt.Sprintf("middleware: tag %q is not registered", tag))
		}
		mws = append(mws, mw)
	}
	return mws
}
