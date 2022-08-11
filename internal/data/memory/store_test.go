package memory

import (
	"reflect"
	"testing"
	"time"

	gtfs "github.com/linhmtran168/511transit/api/protos/gtfs-realtime"
	"github.com/linhmtran168/511transit/internal/models"
	cmap "github.com/orcaman/concurrent-map/v2"
)

func Test_operatorsData_getCache(t *testing.T) {
	type fields struct {
		lastUpdated time.Time
		operators   []*models.Operator
	}
	type args struct {
		isExtended bool
	}

	testOperator := models.Operator{ID: "test", Name: "Test"}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*models.Operator
		want1  bool
	}{
		{
			name:   "Cache data does not exists",
			fields: fields{operators: []*models.Operator{}},
			args:   args{isExtended: false},
			want:   nil,
			want1:  false,
		},
		{
			name:   "Cache data does exists but stalled",
			fields: fields{operators: []*models.Operator{&testOperator}, lastUpdated: time.Now().Add(-OPERATORS_CACHE_TIMEOUT)},
			args:   args{isExtended: false},
			want:   nil,
			want1:  false,
		},
		{
			name:   "Cache data does exists and fresh",
			fields: fields{operators: []*models.Operator{&testOperator}, lastUpdated: time.Now()},
			args:   args{isExtended: false},
			want:   []*models.Operator{&testOperator},
			want1:  true,
		},
		{
			name:   "Cache data stalled but extended param is true",
			fields: fields{operators: []*models.Operator{&testOperator}, lastUpdated: time.Now().Add(-OPERATORS_CACHE_TIMEOUT)},
			args:   args{isExtended: true},
			want:   []*models.Operator{&testOperator},
			want1:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &operatorsData{
				lastUpdated: tt.fields.lastUpdated,
				operators:   tt.fields.operators,
			}
			got, got1 := store.GetCache(tt.args.isExtended)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("operatorsData.getCache() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("operatorsData.getCache() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_tripUpdatesData_getCache(t *testing.T) {
	type fields struct {
		tripMap cmap.ConcurrentMap[*tripEntry]
	}
	type args struct {
		operatorID string
		isExtended bool
	}

	tripId := "test"
	testTrips := []*gtfs.FeedEntity{{Id: &tripId}}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*gtfs.FeedEntity
		want1  bool
	}{
		{
			name:   "Cache data does not exists",
			fields: fields{tripMap: cmap.New[*tripEntry]()},
			args:   args{isExtended: false},
			want:   nil,
			want1:  false,
		},
		{
			name: "Cache data does exists but stalled",
			fields: func() fields {
				tripMap := cmap.New[*tripEntry]()
				tripMap.Set("test", &tripEntry{lastUpdated: time.Now().Add(-TRIP_UPDATE_CACHE_TIMEOUT), trips: testTrips})
				return fields{tripMap: tripMap}
			}(),
			args:  args{operatorID: "test", isExtended: false},
			want:  nil,
			want1: false,
		},
		{
			name: "Cache data does exists and fresh",
			fields: func() fields {
				tripMap := cmap.New[*tripEntry]()
				tripMap.Set("test", &tripEntry{lastUpdated: time.Now(), trips: testTrips})
				return fields{tripMap: tripMap}
			}(),
			args:  args{operatorID: "test", isExtended: false},
			want:  testTrips,
			want1: true,
		},
		{
			name: "Cache data stalled, but extended param is true",
			fields: func() fields {
				tripMap := cmap.New[*tripEntry]()
				tripMap.Set("test", &tripEntry{lastUpdated: time.Now().Add(-TRIP_UPDATE_CACHE_TIMEOUT), trips: testTrips})
				return fields{tripMap: tripMap}
			}(),
			args:  args{operatorID: "test", isExtended: true},
			want:  testTrips,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &tripUpdatesData{
				tripMap: tt.fields.tripMap,
			}
			got, got1 := store.GetCache(tt.args.operatorID, tt.args.isExtended)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tripUpdatesData.getCache() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("tripUpdatesData.getCache() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
