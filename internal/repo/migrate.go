package repo

import (
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/gen"
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func MayBeMigrate(engine *xorm.Engine) error {
	m := migrate.New(engine, migrate.DefaultOptions, []*migrate.Migration{})
	m.InitSchema(func(tx *xorm.Engine) error {
		return tx.Sync(&gen.Book{})
	})
	return m.Migrate()
}
