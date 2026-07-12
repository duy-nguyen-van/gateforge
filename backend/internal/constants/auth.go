package constants

import "time"

// BcryptCost is the work factor for bcrypt password hashes.
const BcryptCost = 10

// SelectionTokenTTL is the lifetime of the JWT used for tenant picker after login.
const SelectionTokenTTL = 5 * time.Minute
