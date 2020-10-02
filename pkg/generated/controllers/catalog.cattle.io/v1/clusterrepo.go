/*
Copyright 2020 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v1 "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ClusterRepoHandler func(string, *v1.ClusterRepo) (*v1.ClusterRepo, error)

type ClusterRepoController interface {
	generic.ControllerMeta
	ClusterRepoClient

	OnChange(ctx context.Context, name string, sync ClusterRepoHandler)
	OnRemove(ctx context.Context, name string, sync ClusterRepoHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() ClusterRepoCache
}

type ClusterRepoClient interface {
	Create(*v1.ClusterRepo) (*v1.ClusterRepo, error)
	Update(*v1.ClusterRepo) (*v1.ClusterRepo, error)
	UpdateStatus(*v1.ClusterRepo) (*v1.ClusterRepo, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1.ClusterRepo, error)
	List(opts metav1.ListOptions) (*v1.ClusterRepoList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ClusterRepo, err error)
}

type ClusterRepoCache interface {
	Get(name string) (*v1.ClusterRepo, error)
	List(selector labels.Selector) ([]*v1.ClusterRepo, error)

	AddIndexer(indexName string, indexer ClusterRepoIndexer)
	GetByIndex(indexName, key string) ([]*v1.ClusterRepo, error)
}

type ClusterRepoIndexer func(obj *v1.ClusterRepo) ([]string, error)

type clusterRepoController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewClusterRepoController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ClusterRepoController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &clusterRepoController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromClusterRepoHandlerToHandler(sync ClusterRepoHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.ClusterRepo
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.ClusterRepo))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *clusterRepoController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.ClusterRepo))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateClusterRepoDeepCopyOnChange(client ClusterRepoClient, obj *v1.ClusterRepo, handler func(obj *v1.ClusterRepo) (*v1.ClusterRepo, error)) (*v1.ClusterRepo, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *clusterRepoController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *clusterRepoController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *clusterRepoController) OnChange(ctx context.Context, name string, sync ClusterRepoHandler) {
	c.AddGenericHandler(ctx, name, FromClusterRepoHandlerToHandler(sync))
}

func (c *clusterRepoController) OnRemove(ctx context.Context, name string, sync ClusterRepoHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromClusterRepoHandlerToHandler(sync)))
}

func (c *clusterRepoController) Enqueue(name string) {
	c.controller.Enqueue("", name)
}

func (c *clusterRepoController) EnqueueAfter(name string, duration time.Duration) {
	c.controller.EnqueueAfter("", name, duration)
}

func (c *clusterRepoController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *clusterRepoController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *clusterRepoController) Cache() ClusterRepoCache {
	return &clusterRepoCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *clusterRepoController) Create(obj *v1.ClusterRepo) (*v1.ClusterRepo, error) {
	result := &v1.ClusterRepo{}
	return result, c.client.Create(context.TODO(), "", obj, result, metav1.CreateOptions{})
}

func (c *clusterRepoController) Update(obj *v1.ClusterRepo) (*v1.ClusterRepo, error) {
	result := &v1.ClusterRepo{}
	return result, c.client.Update(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *clusterRepoController) UpdateStatus(obj *v1.ClusterRepo) (*v1.ClusterRepo, error) {
	result := &v1.ClusterRepo{}
	return result, c.client.UpdateStatus(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *clusterRepoController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), "", name, *options)
}

func (c *clusterRepoController) Get(name string, options metav1.GetOptions) (*v1.ClusterRepo, error) {
	result := &v1.ClusterRepo{}
	return result, c.client.Get(context.TODO(), "", name, result, options)
}

func (c *clusterRepoController) List(opts metav1.ListOptions) (*v1.ClusterRepoList, error) {
	result := &v1.ClusterRepoList{}
	return result, c.client.List(context.TODO(), "", result, opts)
}

func (c *clusterRepoController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), "", opts)
}

func (c *clusterRepoController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*v1.ClusterRepo, error) {
	result := &v1.ClusterRepo{}
	return result, c.client.Patch(context.TODO(), "", name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type clusterRepoCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *clusterRepoCache) Get(name string) (*v1.ClusterRepo, error) {
	obj, exists, err := c.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.ClusterRepo), nil
}

func (c *clusterRepoCache) List(selector labels.Selector) (ret []*v1.ClusterRepo, err error) {

	err = cache.ListAll(c.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ClusterRepo))
	})

	return ret, err
}

func (c *clusterRepoCache) AddIndexer(indexName string, indexer ClusterRepoIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.ClusterRepo))
		},
	}))
}

func (c *clusterRepoCache) GetByIndex(indexName, key string) (result []*v1.ClusterRepo, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.ClusterRepo, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.ClusterRepo))
	}
	return result, nil
}

type ClusterRepoStatusHandler func(obj *v1.ClusterRepo, status v1.RepoStatus) (v1.RepoStatus, error)

type ClusterRepoGeneratingHandler func(obj *v1.ClusterRepo, status v1.RepoStatus) ([]runtime.Object, v1.RepoStatus, error)

func RegisterClusterRepoStatusHandler(ctx context.Context, controller ClusterRepoController, condition condition.Cond, name string, handler ClusterRepoStatusHandler) {
	statusHandler := &clusterRepoStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromClusterRepoHandlerToHandler(statusHandler.sync))
}

func RegisterClusterRepoGeneratingHandler(ctx context.Context, controller ClusterRepoController, apply apply.Apply,
	condition condition.Cond, name string, handler ClusterRepoGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &clusterRepoGeneratingHandler{
		ClusterRepoGeneratingHandler: handler,
		apply:                        apply,
		name:                         name,
		gvk:                          controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterClusterRepoStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type clusterRepoStatusHandler struct {
	client    ClusterRepoClient
	condition condition.Cond
	handler   ClusterRepoStatusHandler
}

func (a *clusterRepoStatusHandler) sync(key string, obj *v1.ClusterRepo) (*v1.ClusterRepo, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type clusterRepoGeneratingHandler struct {
	ClusterRepoGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *clusterRepoGeneratingHandler) Remove(key string, obj *v1.ClusterRepo) (*v1.ClusterRepo, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.ClusterRepo{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *clusterRepoGeneratingHandler) Handle(obj *v1.ClusterRepo, status v1.RepoStatus) (v1.RepoStatus, error) {
	objs, newStatus, err := a.ClusterRepoGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
