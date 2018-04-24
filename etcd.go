// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://coreos.com/blog/etcd
// https://coreos.com/blog/transactional-memory-with-etcd3.html

package etcdfs

import (
	"context"
	"strings"
	"time"

	"github.com/chai2010/jsonmap"
	"github.com/coreos/etcd/clientv3"
)

type EtcdClient struct {
	*clientv3.Client
}

func NewEtcdClient(endpoints []string, timeout time.Duration) (*EtcdClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, err
	}

	return &EtcdClient{Client: cli}, nil
}
func NewEtcdClientWithConfig(cfg clientv3.Config) (*EtcdClient, error) {
	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdClient{Client: cli}, nil
}

func (p *EtcdClient) Close() error {
	return p.Client.Close()
}

func (p *EtcdClient) Get(key string) (val string, ok bool) {
	kv := clientv3.NewKV(p.Client)

	resp, err := kv.Get(context.Background(), key)
	if err != nil {
		return "", false
	}

	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), true
	}

	return "", false
}

func (p *EtcdClient) Set(key, val string) error {
	kv := clientv3.NewKV(p.Client)

	_, err := kv.Put(context.Background(), key, val)
	if err != nil {
		return err
	}

	return nil
}

func (p *EtcdClient) GetValues(keys ...string) (map[string]string, error) {
	kvc := clientv3.NewKV(p.Client)

	var opts []clientv3.Op
	for _, k := range keys {
		opts = append(opts, clientv3.OpGet(k))
	}

	resp, err := kvc.Txn(context.Background()).Then(opts...).Commit()
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, resp_i := range resp.Responses {
		if respRange := resp_i.GetResponseRange(); respRange != nil {
			for _, kv := range respRange.Kvs {
				m[string(kv.Key)] = string(kv.Value)
			}
		}
	}

	return m, nil
}

func (p *EtcdClient) GetValuesByPrefix(keyPrefix string) (map[string]string, error) {
	m := make(map[string]string)
	kv := clientv3.NewKV(p.Client)

	resp, err := kv.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, v := range resp.Kvs {
		m[string(v.Key)] = string(v.Value)
	}

	return m, nil
}

func (p *EtcdClient) SetValues(m map[string]string) error {
	kvc := clientv3.NewKV(p.Client)

	var opts []clientv3.Op
	for k, v := range m {
		opts = append(opts, clientv3.OpPut(k, v))
	}

	_, err := kvc.Txn(context.Background()).Then(opts...).Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *EtcdClient) GetStructValue(keyPrefix string, out interface{}) error {
	m := make(map[string]interface{})
	kv := clientv3.NewKV(p.Client)

	resp, err := kv.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, v := range resp.Kvs {
		key, val := string(v.Key), string(v.Value)
		key = strings.TrimPrefix(key, keyPrefix)
		m[key] = string(val)
	}

	jsonMap := jsonmap.NewJsonMapFromKV(m, "/")
	return jsonMap.ToStruct(out)
}

func (p *EtcdClient) SetStructValue(keyPrefix string, val interface{}) error {
	kvc := clientv3.NewKV(p.Client)

	m := make(map[string]string)
	for k, v := range jsonmap.NewJsonMapFromStruct(val).ToMapString("/") {
		m[keyPrefix+k] = v
	}

	var opts []clientv3.Op
	for k, v := range m {
		opts = append(opts, clientv3.OpPut(k, v))
	}

	_, err := kvc.Txn(context.Background()).Then(opts...).Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *EtcdClient) DelValues(keys ...string) error {
	kvc := clientv3.NewKV(p.Client)

	var opts []clientv3.Op
	for _, k := range keys {
		opts = append(opts, clientv3.OpDelete(k))
	}

	_, err := kvc.Txn(context.Background()).Then(opts...).Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *EtcdClient) DelValuesWithPrefix(keyPrefixs ...string) error {
	kvc := clientv3.NewKV(p.Client)

	var opts []clientv3.Op
	for _, k := range keyPrefixs {
		opts = append(opts, clientv3.OpDelete(k, clientv3.WithPrefix()))
	}

	_, err := kvc.Txn(context.Background()).Then(opts...).Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *EtcdClient) GetAllValues() (map[string]string, error) {
	kvc := clientv3.NewKV(p.Client)

	resp, err := kvc.Get(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, v := range resp.Kvs {
		m[string(v.Key)] = string(v.Value)
	}

	return m, nil
}
func (p *EtcdClient) Clear() error {
	kvc := clientv3.NewKV(p.Client)

	_, err := kvc.Delete(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}
