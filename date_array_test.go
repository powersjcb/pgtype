package pgtype_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
)

func TestDateArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "date[]", []interface{}{
		&pgtype.DateArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     pgtype.Present,
		},
		&pgtype.DateArray{
			Elements: []pgtype.Date{
				{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
				{Status: pgtype.Null},
			},
			Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     pgtype.Present,
		},
		&pgtype.DateArray{Status: pgtype.Null},
		&pgtype.DateArray{
			Elements: []pgtype.Date{
				{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
				{Time: time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
				{Time: time.Date(2017, 2, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
				{Time: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
				{Status: pgtype.Null},
				{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
			},
			Dimensions: []pgtype.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     pgtype.Present,
		},
		&pgtype.DateArray{
			Elements: []pgtype.Date{
				{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
				{Time: time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
				{Time: time.Date(2015, 2, 3, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
				{Time: time.Date(2015, 2, 4, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
			},
			Dimensions: []pgtype.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: pgtype.Present,
		},
	})
}

func TestDateArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result pgtype.DateArray
	}{
		{
			source: []time.Time{time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)},
			result: pgtype.DateArray{
				Elements:   []pgtype.Date{{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     pgtype.Present},
		},
		{
			source: (([]time.Time)(nil)),
			result: pgtype.DateArray{Status: pgtype.Null},
		},
	}

	for i, tt := range successfulTests {
		var r pgtype.DateArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestDateArrayAssignTo(t *testing.T) {
	var timeSlice []time.Time

	simpleTests := []struct {
		src      pgtype.DateArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: pgtype.DateArray{
				Elements:   []pgtype.Date{{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     pgtype.Present,
			},
			dst:      &timeSlice,
			expected: []time.Time{time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			src:      pgtype.DateArray{Status: pgtype.Null},
			dst:      &timeSlice,
			expected: (([]time.Time)(nil)),
		},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src pgtype.DateArray
		dst interface{}
	}{
		{
			src: pgtype.DateArray{
				Elements:   []pgtype.Date{{Status: pgtype.Null}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     pgtype.Present,
			},
			dst: &timeSlice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}

}

const dateFormat = "2006-01-02"

func TestDateArrayMarshalJSON(t *testing.T) {
	t1s := "2020-06-13"
	t1, err := time.Parse(dateFormat, t1s)
	if err != nil {
		t.Fatal(err)
	}
	t2s := "2020-06-13"
	t2, err := time.Parse(dateFormat, t2s)
	if err != nil {
		t.Fatal(err)
	}
	successfulTests := []struct {
		source pgtype.DateArray
		result string
	}{
		{source: pgtype.DateArray{Elements: []pgtype.Date{}, Status: pgtype.Null}, result: "null"},
		{source: pgtype.DateArray{Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 0}}, Status: pgtype.Present}, result: "[]"},
		{
			source: pgtype.DateArray{
				Elements:   []pgtype.Date{{Time: t1, Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{Length: 1, LowerBound: 1}},
				Status:     pgtype.Present,
			},
			result: fmt.Sprintf(`["%s"]`, t1s),
		},
		{
			source: pgtype.DateArray{
				Elements: []pgtype.Date{
					{Time: t1, Status: pgtype.Present},
					{Status: pgtype.Null},
					{Time: t2, Status: pgtype.Present},
				},
				Dimensions: []pgtype.ArrayDimension{{Length: 3, LowerBound: 1}},
				Status:     pgtype.Present,
			},
			result: fmt.Sprintf(`["%s",null,"%s"]`, t1s, t2s),
		},
	}
	for i, tt := range successfulTests {
		r, err := tt.source.MarshalJSON()
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if string(r) != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, string(r))
		}
	}
}
