package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/webitel/webitel-fts/infra/searchengine"
	"github.com/webitel/webitel-fts/infra/sql"
	"github.com/webitel/webitel-fts/pkg/client"
	"github.com/webitel/wlog"
)

type ImportData struct {
	log    *wlog.Logger
	sql    sql.Store
	search searchengine.SearchEngine
}

func NewImport(log *wlog.Logger, s sql.Store, search searchengine.SearchEngine) *ImportData {
	return &ImportData{
		log:    log,
		sql:    s,
		search: search,
	}
}

func (i *ImportData) Import(ctx context.Context, q string, colId string, colDomainId string, objectName string) error {
	rows, err := i.sql.Query(ctx, q)
	if err != nil {
		return err
	}

	defer rows.Close()

	fields := rows.Columns()

	cols := make([]interface{}, len(fields))
	colPtrs := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		colPtrs[i] = &cols[i]
	}

	for rows.Next() {
		if err = rows.Scan(colPtrs...); err != nil {
			return err
		}
		msg := client.Message{
			Id:         "",
			DomainId:   0,
			ObjectName: objectName,
			Date:       0,
			Body:       nil,
		}
		body := make(map[string]interface{})
		for i, col := range cols {
			switch fields[i] {
			case colId:
				msg.Id = client.MessageId(fmt.Sprintf("%v", col))
			case colDomainId:
				var ok bool
				msg.DomainId, ok = col.(int64)
				if !ok {
					panic(col)
				}
			default:
				body[fields[i]] = col
			}

		}
		msg.Body, err = json.Marshal(body)
		if err != nil {
			panic(err.Error())
		}

		err = i.search.Insert(
			ctx,
			fmt.Sprintf("%v", msg.Id),
			fmt.Sprintf("%v_%v", msg.ObjectName, msg.DomainId),
			msg.Body)

		if err != nil {
			panic(err)
		} else {
			i.log.Debug(fmt.Sprintf("create doc %v", msg.Id))
		}

	}

	return nil
}
