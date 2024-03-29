/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package securestoragadapter provides interfaces for defining secure manager for variable and secret.
package securestoragadapter

import (
	"context"
	"errors"

	vaultApi "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/space/domain/securestorage"
)

// SecureStorageAdapter is a function return secure storage
func SecureStorageAdapter(client *vaultApi.Client, basePath string) *vaultAdapter {
	return &vaultAdapter{
		client:   client,
		basePath: basePath,
	}
}

type vaultAdapter struct {
	client   *vaultApi.Client
	basePath string
}

func (v vaultAdapter) SaveSpaceEnvSecret(es securestorage.SpaceEnvSecret) error {
	storageData := map[string]interface{}{
		es.Name: es.Value,
	}
	storageValueList, err := v.client.KVv2(v.basePath).Get(context.Background(), es.Path)
	if err != nil && !errors.Is(err, vaultApi.ErrSecretNotFound) {
		logrus.Errorf("get storage value failed: %v", err)
		return err
	}
	if storageValueList == nil || len(storageValueList.Data) == 0 {
		return v.setSpaceEnvSecret(es, storageData)
	}
	for key, value := range storageData {
		storageValueList.Data[key] = value
	}
	_, err = v.client.KVv2(v.basePath).Patch(context.Background(), es.Path, storageValueList.Data)
	if err != nil && !errors.Is(err, vaultApi.ErrSecretNotFound) {
		logrus.Errorf("unable to patch storage value: %v", err)
		return err
	}
	return nil
}

func (v vaultAdapter) setSpaceEnvSecret(es securestorage.SpaceEnvSecret, storageData map[string]interface{}) error {
	_, err := v.client.KVv2(v.basePath).Put(context.Background(), es.Path, storageData)
	if err != nil && !errors.Is(err, vaultApi.ErrSecretNotFound) {
		logrus.Errorf("unable to write storage value: %v", err)
		return err
	}
	return nil
}

func (v vaultAdapter) DeleteSpaceEnvSecret(path string, key string) error {
	storageValueList, err := v.client.KVv2(v.basePath).Get(context.Background(), path)
	if err != nil {
		logrus.Errorf("get storage value failed: %v", err)
		return err
	}
	if storageValueList == nil {
		return errors.New("get storage value fail")
	}
	delete(storageValueList.Data, key)
	if len(storageValueList.Data) == 0 {
		return v.client.KVv2(v.basePath).Delete(context.Background(), path)
	}
	_, err = v.client.KVv2(v.basePath).Put(context.Background(), path, storageValueList.Data)
	if err != nil && !errors.Is(err, vaultApi.ErrSecretNotFound) {
		logrus.Errorf("unable to delete storage value: %v", err)
		return err
	}
	return nil
}

func (v vaultAdapter) GetAllSpaceEnvSecret(es securestorage.SpaceEnvSecret) (string, error) {
	storageValueList, err := v.client.KVv2(v.basePath).Get(context.Background(), es.Path)
	if err != nil {
		logrus.Errorf("unable to get storage value: %v", err)
		return "", err
	}
	value, ok := storageValueList.Data[es.Name].(string)
	if !ok {
		logrus.Errorf("value type assertion failed: %T", storageValueList.Data[es.Name])
		return "", err
	}
	return value, nil
}
