package recurse

import (
	"testing"

	"github.com/zew/exceldb/config"
	"github.com/zew/logx"
	"github.com/zew/util"

	_ "github.com/go-sql-driver/mysql"
)

type Foo struct {
	FooId int    `db:"foo_id, primarykey, autoincrement"`
	Foo1  string `db:"foo1"`
}

type Bar struct {
	BarId int    `db:"bar_id, primarykey, autoincrement"`
	FooId int    `db:"foo_id"`
	Bar1  string `db:"bar1"`
}

type FooAndBars struct {
	Foo
	Bar
	FooId int `db:"foo_id"` // explicit

	Bars []Bar // ignored; not assignable

	T2 Bar `db:"t2"` // ignored; not assignable
}

var results []FooAndBars

func Test_2(t *testing.T) {

	dbMap := DBMap(testDB)
	// defer db.Close()

	dbMap.AddTable(Foo{})
	dbMap.AddTable(Bar{})

	//
	err := dbMap.DropTablesIfExists()
	if err != nil {
		logx.Printf("mysql incapable of drop-if-exits: %v", err)
	}
	err = dbMap.CreateTablesIfNotExists()
	util.CheckErr(err)
	err = dbMap.CreateIndex()
	util.CheckErr(err)
	logx.Printf("gorp tables created")

	tables := []string{}

	if !config.Config.SQLite {
		_, err = dbMap.Select(&tables, `show tables;`)
		util.CheckErr(err)
		logx.Printf(util.IndentedDump(tables))
	}

	//
	// Data
	//
	err = dbMap.TruncateTables()
	util.CheckErr(err)

	f1 := Foo{
		FooId: 1,
		Foo1:  "Foo1",
	}
	f2 := Foo{
		FooId: 2,
		Foo1:  "Foo1",
	}
	err = dbMap.Insert(&f1, &f2)
	util.CheckErr(err)

	b1 := Bar{
		BarId: 1,
		FooId: 1,
		Bar1:  "Bar1",
	}
	b2 := Bar{
		BarId: 2,
		FooId: 1,
		Bar1:  "Bar2",
	}
	err = dbMap.Insert(&b1, &b2)
	util.CheckErr(err)

	// This - strangely - requires
	// uppper case table names.
	_, err = dbMap.Select(&results,
		`SELECT 
				t1.foo_id    foo_id ,
				t1.foo1      foo1 ,
				t2.bar_id    bar_id ,
				t2.foo_id    foo_id ,
				t2.bar1      bar1
			FROM 				`+TableName(Foo{})+` t1
					LEFT JOIN 	`+TableName(Bar{})+` t2 USING(foo_id)
			WHERE 			1=1
					AND		foo_id = 1
			`,
		// 1,
	)
	util.CheckErr(err)

	for _, x := range results {
		logx.Printf("%+v", x)
	}

	dbMap.DropTablesIfExists()

}
