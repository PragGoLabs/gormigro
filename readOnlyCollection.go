package gormigro

type ReadOnlyCollection struct {
	migrations []Migration
}

func NewReadOnlyCollection(migrations []Migration) ReadOnlyCollection {
	return ReadOnlyCollection{
		migrations: migrations,
	}
}

func (rc ReadOnlyCollection) List() []Migration {
	return rc.migrations
}