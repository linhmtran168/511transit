package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	gtfs "github.com/linhmtran168/511transit/api/protos/gtfs-realtime"
	"github.com/linhmtran168/511transit/internal/mock"
	"github.com/linhmtran168/511transit/internal/models"
	transitapi "github.com/linhmtran168/511transit/internal/transit-api"
	"google.golang.org/protobuf/proto"
)

func createOperatorCache(ctrl *gomock.Controller) *mock.MockOperatorCache {
	mockCache := mock.NewMockOperatorCache(ctrl)
	mockCache.EXPECT().Lock()
	mockCache.EXPECT().Unlock()
	return mockCache
}
func TestMemoryRepository_GetOperators(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData1 := []*models.Operator{{ID: "1", Name: "test"}}
	testResponse, _ := json.Marshal(testData1)
	testData2 := []*models.Operator{{ID: "2", Name: "test"}}

	apiClient := mock.NewMockTransitAPI(ctrl)
	apiClient.EXPECT().GetOperators().
		Return(testResponse, nil).
		AnyTimes()

	type fields struct {
		apiClient     transitapi.TransitAPI
		operatorStore OperatorCache
		tripStore     TripUpdateCache
	}

	tests := []struct {
		name    string
		fields  fields
		want    []*models.Operator
		wantErr bool
	}{
		{
			"No data from cache, need to call api",
			fields{
				apiClient: apiClient,
				operatorStore: func() OperatorCache {
					store := createOperatorCache(ctrl)
					store.EXPECT().GetCache(false).Times(1).
						Return(nil, false)
					store.EXPECT().UpdateCache(testData1).Times(1)
					return store
				}(),
			},
			testData1,
			false,
		},
		{
			"Has data from cache no need to call api",
			fields{
				apiClient: apiClient,
				operatorStore: func() OperatorCache {
					store := createOperatorCache(ctrl)
					store.EXPECT().GetCache(false).Times(1).
						Return(testData2, true)
					return store
				}(),
			},
			testData2,
			false,
		},
		{
			"Error from api, use data from cache",
			fields{
				apiClient: func() transitapi.TransitAPI {
					apiClient := mock.NewMockTransitAPI(ctrl)
					apiClient.EXPECT().GetOperators().
						Return(nil, errors.New("API error"))
					return apiClient
				}(),
				operatorStore: func() OperatorCache {
					store := createOperatorCache(ctrl)
					store.EXPECT().GetCache(false).Times(1).
						Return(nil, false)
					store.EXPECT().GetCache(true).Times(1).
						Return(testData2, true)
					return store
				}(),
			},
			testData2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MemoryRepository{
				apiClient:     tt.fields.apiClient,
				operatorStore: tt.fields.operatorStore,
				tripStore:     tt.fields.tripStore,
			}
			got, err := p.GetOperators()
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryRepository.GetOperators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryRepository.GetOperators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryRepository_GetTripUpdates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testID1, testID2 := "1", "2"
	version := "1"
	entity1, entity2 := []*gtfs.FeedEntity{{Id: &testID1}}, []*gtfs.FeedEntity{{Id: &testID2}}
	feedData1, _ := proto.Marshal(&gtfs.FeedMessage{Entity: entity1, Header: &gtfs.FeedHeader{GtfsRealtimeVersion: &version}})

	apiClient := mock.NewMockTransitAPI(ctrl)
	operatorID := "test"
	apiClient.EXPECT().GetTripUpdates(operatorID).
		Return(feedData1, nil).
		AnyTimes()

	type fields struct {
		apiClient     transitapi.TransitAPI
		operatorStore OperatorCache
		tripStore     TripUpdateCache
	}
	type args struct {
		operatorID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*gtfs.FeedEntity
		wantErr bool
	}{
		{
			"No data from cache, need to call api",
			fields{
				apiClient: apiClient,
				tripStore: func() TripUpdateCache {
					store := mock.NewMockTripUpdateCache(ctrl)
					store.EXPECT().GetCache(operatorID, false).Times(1).
						Return(nil, false)
					store.EXPECT().UpdateCache(operatorID, gomock.AssignableToTypeOf(entity1)).Times(1)
					return store
				}(),
			},
			args{operatorID: operatorID},
			entity1,
			false,
		},
		{
			"Has data from cache, no need to call api",
			fields{
				apiClient: apiClient,
				tripStore: func() TripUpdateCache {
					store := mock.NewMockTripUpdateCache(ctrl)
					store.EXPECT().GetCache(operatorID, false).Times(1).
						Return(entity2, true)
					return store
				}(),
			},
			args{operatorID: operatorID},
			entity2,
			false,
		},
		{
			"Error from api, use data from cache",
			fields{
				apiClient: func() transitapi.TransitAPI {
					apiClient := mock.NewMockTransitAPI(ctrl)
					apiClient.EXPECT().GetTripUpdates(operatorID).
						Return(nil, errors.New("Api error"))
					return apiClient
				}(),
				tripStore: func() TripUpdateCache {
					store := mock.NewMockTripUpdateCache(ctrl)
					store.EXPECT().GetCache(operatorID, false).Times(1).
						Return(nil, false)
					store.EXPECT().GetCache(operatorID, true).Times(1).
						Return(entity2, true)
					return store
				}(),
			},
			args{operatorID: operatorID},
			entity2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MemoryRepository{
				apiClient:     tt.fields.apiClient,
				operatorStore: tt.fields.operatorStore,
				tripStore:     tt.fields.tripStore,
			}
			got, err := p.GetTripUpdates(tt.args.operatorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryRepository.GetTripUpdates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if fmt.Sprintf("%+v", got) != fmt.Sprintf("%+v", tt.want) {
				t.Errorf("MemoryRepository.GetTripUpdates() = %v, want %v", got, tt.want)
			}
		})
	}
}
