package main

import "os"

type UniversalStorage struct {
	s3  *S3Storage
	local *LocalStorage
}

func (store *UniversalStorage) List(p string) ([]*StorageDirectory, []*StorageFile, error) {
	if store.s3 != nil {
		return store.s3.List(p)
	}
	if store.local != nil {
		return store.local.List(p)
	}
	return nil, nil, nil
}

func (store *UniversalStorage) FileContentURL(p string) (string, error) {
	if store.s3 != nil {
		return store.s3.FileContentURL(p)
	}
	if store.local != nil {
		return store.local.FileContentURL(p)
	}
	return "", nil
}

func (store *UniversalStorage) OpenFile(p string) (*os.File, error) {
	if store.local != nil {
		return store.local.OpenFile(p)
	}
	return nil, nil
}


