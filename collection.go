package gormigro

import (
	"reflect"
	"sort"
)

type SortingType string

const (
	MigrationId SortingType = "migrationId"
)

type OrderType string

const (
	Asc  OrderType = "asc"
	Desc OrderType = "desc"
)

type Collection struct {
	migrations []Migration
}

func NewCollection() *Collection {
	return &Collection{}
}

func NewCollectionWithMigrations(migrations []Migration) *Collection {
	return &Collection{
		migrations: migrations,
	}
}

func (c *Collection) Append(m Migration) error {
	if c.Contains(m) {
		return CreateMigrationAlreadyExistsError(m.ID)
	}

	c.migrations = append(c.migrations, m)

	return nil
}

func (c *Collection) Contains(m Migration) bool {
	for _, em := range c.migrations {
		if reflect.DeepEqual(em, m) {
			return true
		}
	}

	return false
}

func (c *Collection) SliceFrom(migrationId string) *Collection {
	subCollection := NewCollection()

	appendable := false
	for _, mc := range c.migrations {
		if mc.ID == migrationId {
			appendable = true
			continue
		}
		if !appendable {
			continue
		}

		subCollection.Append(mc)
	}

	return subCollection
}

func (c *Collection) SortBy(sortingType SortingType, order OrderType) *Collection {
	switch sortingType {
	case MigrationId:
		return c.sortByMigrationId(order)
		break
	}

	panic("Undefined sorting type")

	return c
}

func (c *Collection) List() []Migration {
	return c.migrations
}

func (c *Collection) Empty() bool {
	return len(c.migrations) == 0
}

func (c *Collection) Count() int {
	return len(c.migrations)
}

func (c *Collection) ReadOnly() ReadOnlyCollection {
	return NewReadOnlyCollection(c.migrations)
}

func (c *Collection) sortByMigrationId(orderType OrderType) *Collection {
	exportedMigrations := map[string]Migration{}
	var keysForSort []string

	for _, m := range c.migrations {
		exportedMigrations[m.ID] = m
		keysForSort = append(keysForSort, m.ID)
	}

	if orderType == Desc {
		sort.Sort(sort.Reverse(sort.StringSlice(keysForSort)))

	} else {
		sort.Strings(keysForSort)
	}

	sortedCollection := NewCollection()
	for _, k := range keysForSort {
		sortedCollection.Append(exportedMigrations[k])
	}

	return sortedCollection
}
