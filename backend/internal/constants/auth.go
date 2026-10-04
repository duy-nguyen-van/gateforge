package constants

import "time"

// BcryptCost is the work factor for bcrypt password hashes.
const BcryptCost = 10

// TOTPIssuer is the name authenticator apps show. It stays GateForge when APP_NAME differs.
const TOTPIssuer = "GateForge"

// SelectionTokenTTL is the lifetime of the JWT used for tenant picker after login.
const SelectionTokenTTL = 5 * time.Minute
