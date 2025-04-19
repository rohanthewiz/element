package comps

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// genRandId generates a random id for an element
func genRandId() string {
	bytes := make([]byte, 12)

	if _, err := rand.Read(bytes); err != nil {
		return hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixMilli())))
	}

	return hex.EncodeToString(bytes)[:12]
}
