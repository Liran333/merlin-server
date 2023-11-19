package controller

import (
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/utils"
)

var (
	apiConfig         APIConfig
	encryptHelper     utils.SymmetricEncryption
	encryptHelperCSRF utils.SymmetricEncryption
	log               *logrus.Entry
)

func Init(cfg *APIConfig, l *logrus.Entry) error {
	log = l
	apiConfig = *cfg

	e, err := utils.NewSymmetricEncryption(cfg.EncryptionKey, "")
	if err != nil {
		return err
	}

	csrfe, err := utils.NewSymmetricEncryption(cfg.EncryptionKeyForCSRF, "")
	if err != nil {
		return err
	}

	encryptHelper = e
	encryptHelperCSRF = csrfe

	return nil
}

type APIConfig struct {
	TokenKey             []byte `json:"token_key"                   required:"true"`
	TokenExpiry          int64  `json:"token_expiry"                required:"true"`
	EncryptionKey        []byte `json:"encryption_key"              required:"true"`
	EncryptionKeyForCSRF []byte `json:"encryption_key_csrf"         required:"true"`
}
