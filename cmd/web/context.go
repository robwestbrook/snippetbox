package main

// Define a custom type contextKey
type contextKey string

// Set the isAuthenticatedContextKey constant key
// to "isAuthenticated"
const isAuthenticatedContextKey = contextKey("isAuthenticated")