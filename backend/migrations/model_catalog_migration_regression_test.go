package migrations

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Applied migrations are identified by full filename, not numeric prefix.
func TestModelCatalogMigrationsKeepPublishedIdentity(t *testing.T) {
	checksums := map[string]string{
		"239_model_catalog.sql": "d7d9a046a68048f46bd9cc9aeff6385845a24a2d9f087e2b97a101ce0cf36104",
		"240_model_catalog_vendor_soft_delete_platform_cleanup.sql": "a68721ce923c3f9567c3e61ec8c86e5d69415a92c74d945827b0385b5dfa89af",
	}
	for name, checksum := range checksums {
		t.Run(name, func(t *testing.T) {
			content, err := FS.ReadFile(name)
			require.NoError(t, err)
			require.Equal(t, checksum, fmt.Sprintf("%x", sha256.Sum256(content)))
		})
	}
	for _, name := range []string{
		"239_channel_reasoning_effort_multipliers.sql",
		"240_affiliate_ledger_operation_id.sql",
		"241_add_payment_order_bonus_amount.sql",
		"241_add_typesafe_platform.sql",
	} {
		_, err := FS.ReadFile(name)
		require.NoError(t, err)
	}
}
