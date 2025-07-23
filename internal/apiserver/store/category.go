package store

import (
	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	genericstore "github.com/clin211/miniblog-v2/pkg/store"
)

// CategoryStore 定义了 category 模块在 store 层所实现的方法
type CategoryStore interface {
	genericstore.IStore[model.CategoryM]
}

// categoryStore 是 CategoryStore 接口的实现
type categoryStore struct {
	*genericstore.Store[model.CategoryM]
}

// 确保 categoryStore 实现了 CategoryStore 接口
var _ CategoryStore = (*categoryStore)(nil)

// newCategoryStore 创建 categoryStore 的实例
func newCategoryStore(store *datastore) *categoryStore {
	return &categoryStore{
		Store: genericstore.NewStore[model.CategoryM](store, genericstore.NewLogger()),
	}
}
