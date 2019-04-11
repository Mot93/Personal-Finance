package financetest

import (
	"fmt"
	"os"
	"testing"

	"github.com/Mot93/personalfinance/financedatabase"
)

func TestStartDatabase(t *testing.T) {
	financedatabase.InitDB("testfinance.db")
	// Erasing the test database after it's usage
	defer func() {
		err := os.Remove("testfinance.db")
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}()
}
