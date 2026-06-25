package service

import infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"

var (
	ErrModelCatalogNotFound = infraerrors.NotFound("MODEL_CATALOG_NOT_FOUND", "model catalog entry not found")
	ErrModelCatalogExists   = infraerrors.Conflict("MODEL_CATALOG_EXISTS", "model catalog entry already exists for this platform and model")
	ErrModelVendorNotFound  = infraerrors.NotFound("MODEL_VENDOR_NOT_FOUND", "model vendor not found")
)
