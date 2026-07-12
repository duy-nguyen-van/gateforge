package middlewares

import (
	"os"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/testutil"
)

func TestMain(m *testing.M) {
	testutil.InitLogger()
	os.Exit(m.Run())
}
