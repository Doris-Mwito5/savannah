package models

import (
	"errors"
	"fmt"
	"github/Doris-Mwito5/savannah-pos/internal/utiltime"
	"time"
)

type Filter struct {
	Page       int
	Per        int
	From       string
	To         string
	Token      string
	Term       string
	UUID       string
	Status     *string
	Type       string
	Valid      *bool
	Active     *bool
	Year       string
	Reference  string
	FromTime   *time.Time
	ToTime     *time.Time
	ShopID     *string
	CategoryID *string
	OrderID    *int64
}

func (f *Filter) ConvertTime() error {
	if f.From == "" || f.To == "" {
		return errors.New("from or to filter time is empty")
	}

	fromTime, err := utiltime.ParseTime(f.From)
	if err != nil {
		return fmt.Errorf("time_filter: parse from time [%v], err [%v]", f.From, err)
	}

	f.FromTime = &fromTime

	toTime, err := utiltime.ParseTime(f.To)
	if err != nil {
		return fmt.Errorf("time_filter: parse to time [%v], err [%v]", f.To, err)
	}

	f.ToTime = &toTime

	return nil
}

func (f *Filter) NoPagination() *Filter {
	return &Filter{
		From:       f.From,
		To:         f.To,
		Term:       f.Term,
		UUID:       f.UUID,
		Status:     f.Status,
		Type:       f.Type,
		Token:      f.Token,
		Valid:      f.Valid,
		Active:     f.Active,
		ShopID:     f.ShopID,
		CategoryID: f.CategoryID,
	}
}

func (f *Filter) TimeFilterSet() bool {

	return f.From != "" && f.To != ""
}

func (f *Filter) ExportLimit() {

	f.Page = 1
	f.Per = 1000
}
