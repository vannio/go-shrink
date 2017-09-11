package db

import (
	"errors"
	"reflect"
	"testing"
)

type TestRow struct {
	Slug string
	URL  string
}

func (t TestRow) FindBySlug(slug string) (Row, error) {
	return Row{Slug: slug, URL: "https://www.original.url"}, nil
}

func (t TestRow) FindByURL(url string) (Row, error) {
	return Row{Slug: "abc123", URL: url}, nil
}

func (t TestRow) Inject(url string) error {
	return errors.New("Something went badly wrong, abandon ship!")
}

func (t TestRow) IncrementAccessCount() error {
	return errors.New("Something went a little wrong, let's pretend it didn't.")
}

func TestFindRowBySlug(t *testing.T) {
	row, _ := FindRowBySlug(TestRow{}, "abc123")

	if !reflect.DeepEqual(row.URL, "https://www.original.url") {
		t.Fail()
	}
}

func TestFindRowByURL(t *testing.T) {
	row, _ := FindRowByURL(TestRow{}, "https://www.original.url")

	if !reflect.DeepEqual(row.Slug, "abc123") {
		t.Fail()
	}
}

func TestInjectRow(t *testing.T) {
	err := InjectRow(TestRow{}, "https://www.original.url")

	if !reflect.DeepEqual(err.Error(), "Something went badly wrong, abandon ship!") {
		t.Fail()
	}
}

func TestIncrementAccessCount(t *testing.T) {
	err := IncrementAccessCount(TestRow{})

	if !reflect.DeepEqual(err.Error(), "Something went a little wrong, let's pretend it didn't.") {
		t.Fail()
	}
}
