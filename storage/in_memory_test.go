package storage

import (
	"reflect"
	"testing"
)

func Test_inMemoryStorage_Write(t *testing.T) {
	type args struct {
		events []Event
		query  *Query
		result *ResultSet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"One attribute written and retrieved", args{
			[]Event{{
				Attributes: map[string]string{"a": "a"},
				Timestamp:  1,
			}},
			&Query{
				Attributes:     map[string]string{"a": "a"},
				StartTimestamp: 1,
				EndTimestamp:   1,
			},
			&ResultSet{
				Attributes: map[string]string{"a": "a"},
				Value:      1,
			},
		}, false},
		{"Composite attributes written and retrieved", args{
			[]Event{{
				Attributes: map[string]string{"a": "a", "b": "b"},
				Timestamp:  1,
			}},
			&Query{
				Attributes:     map[string]string{"a": "a", "b": "b"},
				StartTimestamp: 1,
				EndTimestamp:   1,
			},
			&ResultSet{
				Attributes: map[string]string{"a": "a", "b": "b"},
				Value:      1,
			},
		}, false},
		{"Query attributes not in db, returns 0", args{
			[]Event{{
				Attributes: map[string]string{"a": "a"},
				Timestamp:  1,
			}},
			&Query{
				Attributes:     map[string]string{"b": "b"},
				StartTimestamp: 1,
				EndTimestamp:   1,
			},
			&ResultSet{
				Attributes: map[string]string{"b": "b"},
				Value:      0,
			},
		}, false},
		{"No data for provided time range, returns 0", args{
			[]Event{{
				Attributes: map[string]string{"a": "a"},
				Timestamp:  1,
			}},
			&Query{
				Attributes:     map[string]string{"a": "a"},
				StartTimestamp: 2,
				EndTimestamp:   4,
			},
			&ResultSet{
				Attributes: map[string]string{"a": "a"},
				Value:      0,
			},
		}, false},
		{"Multiple events, same attributes, same timestamp", args{
			[]Event{
				{
					Attributes: map[string]string{"a": "a"},
					Timestamp:  1,
				},
				{
					Attributes: map[string]string{"a": "a"},
					Timestamp:  1,
				},
			},
			&Query{
				Attributes:     map[string]string{"a": "a"},
				StartTimestamp: 1,
				EndTimestamp:   1,
			},
			&ResultSet{
				Attributes: map[string]string{"a": "a"},
				Value:      2,
			},
		}, false},
		{"Multiple events, same attributes, over 3 timestamps not in order", args{
			[]Event{
				{
					Attributes: map[string]string{"a": "a"},
					Timestamp:  3,
				},
				{
					Attributes: map[string]string{"a": "a"},
					Timestamp:  1,
				},
				{
					Attributes: map[string]string{"a": "a"},
					Timestamp:  2,
				},
			},
			&Query{
				Attributes:     map[string]string{"a": "a"},
				StartTimestamp: 1,
				EndTimestamp:   3,
			},
			&ResultSet{
				Attributes: map[string]string{"a": "a"},
				Value:      3,
			},
		}, false},
		{"Event with multiple attributes, query by one attribute", args{
			[]Event{
				{
					Attributes: map[string]string{"a": "a", "b": "b"},
					Timestamp:  1,
				},
			},
			&Query{
				Attributes:     map[string]string{"b": "b"},
				StartTimestamp: 1,
				EndTimestamp:   3,
			},
			&ResultSet{
				Attributes: map[string]string{"b": "b"},
				Value:      1,
			},
		}, false},
		{"Exclude even out of query's range", args{
			[]Event{
				{
					Attributes: map[string]string{"a": "a", "b": "b"},
					Timestamp:  1,
				},
				{
					Attributes: map[string]string{"a": "a", "b": "b"},
					Timestamp:  3,
				},
			},
			&Query{
				Attributes:     map[string]string{"b": "b"},
				StartTimestamp: 1,
				EndTimestamp:   2,
			},
			&ResultSet{
				Attributes: map[string]string{"b": "b"},
				Value:      1,
			},
		}, false},
		{"Four attributes, query by two, not consecutive", args{
			[]Event{
				{
					Attributes: map[string]string{"a": "a", "b": "b", "c": "c", "d": "d"},
					Timestamp:  1,
				},
			},
			&Query{
				Attributes:     map[string]string{"b": "b", "d": "d"},
				StartTimestamp: 1,
				EndTimestamp:   2,
			},
			&ResultSet{
				Attributes: map[string]string{"b": "b", "d": "d"},
				Value:      1,
			},
		}, false},
		{"Events with different attribute values", args{
			[]Event{
				{
					Attributes: map[string]string{"a": "a1"},
					Timestamp:  1,
				},
				{
					Attributes: map[string]string{"a": "a2"},
					Timestamp:  1,
				},
			},
			&Query{
				Attributes:     map[string]string{"a": "a1"},
				StartTimestamp: 1,
				EndTimestamp:   2,
			},
			&ResultSet{
				Attributes: map[string]string{"a": "a1"},
				Value:      1,
			},
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &inMemoryStorage{
				tree: newTree(),
			}
			for _, e := range tt.args.events {
				if err := s.Write(&e); (err != nil) != tt.wantErr {
					t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			result, err := s.Query(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(result, tt.args.result) {
				t.Errorf("ResultSet = %v, want %v", result, tt.args.result)
			}
		})
	}
}
