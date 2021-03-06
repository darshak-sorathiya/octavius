// Package metadata implements metadata repository related functions
package metadata

import (
	"context"
	"errors"
	"octavius/internal/pkg/constant"
	"octavius/internal/pkg/db/etcd"
	"octavius/internal/pkg/log"
	"octavius/internal/pkg/protofiles"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	log.Init("info", "", false, 1)
}

func Test_metadataRepository_SaveMetadata(t *testing.T) {
	mockClient := new(etcd.ClientMock)
	metadataVal := &protofiles.Metadata{
		Author:      "littlestar642",
		ImageName:   "demo image",
		Name:        "test data",
		Description: "sample test metadata",
	}
	val, err := proto.Marshal(metadataVal)
	if err != nil {
		t.Error("error in marshalling metadata")
	}
	mockClient.On("PutValue", "metadata/test data", string(val)).Return(nil)
	mockClient.On("GetValue", "metadata/test data").Return("", nil)

	testMetadataRepo := NewMetadataRepository(mockClient)
	ctx := context.Background()
	sr, err := testMetadataRepo.SaveMetadata(ctx, "test data", metadataVal)

	if err != nil {
		t.Error(err, "saving metadata failed")
	}
	if sr.Name != "test data" {
		t.Errorf("expected %s, got %s", "test data", sr.Name)
	}

	mockClient.AssertExpectations(t)
}

func Test_metadataRepository_SaveMetadata_KeyAlreadyPresent(t *testing.T) {
	mockClient := new(etcd.ClientMock)
	metadataVal := &protofiles.Metadata{
		Author:      "littlestar642",
		ImageName:   "demo image",
		Name:        "test data",
		Description: "sample test metadata",
	}
	val, err := proto.Marshal(metadataVal)
	if err != nil {
		t.Error("error in marshalling metadata")
	}
	mockClient.On("PutValue", "metadata/test data", string(val)).Return(nil)
	mockClient.On("GetValue", "metadata/test data").Return("some key", nil)

	testMetadataRepo := NewMetadataRepository(mockClient)
	ctx := context.Background()
	_, err = testMetadataRepo.SaveMetadata(ctx, "test data", metadataVal)

	if err.Error() != status.Error(codes.AlreadyExists, constant.Etcd+constant.KeyAlreadyPresent).Error() {
		t.Error("key already present error expected")
	}
}

func Test_metadataRepository_SaveMetadata_GetValueError(t *testing.T) {
	mockClient := new(etcd.ClientMock)
	metadataVal := &protofiles.Metadata{
		Author:      "littlestar642",
		ImageName:   "demo image",
		Name:        "test data",
		Description: "sample test metadata",
	}
	val, err := proto.Marshal(metadataVal)
	if err != nil {
		t.Error("error in marshalling metadata")
	}
	mockClient.On("PutValue", "metadata/test data", string(val)).Return(nil)
	mockClient.On("GetValue", "metadata/test data").Return("", errors.New("some error"))

	testMetadataRepo := NewMetadataRepository(mockClient)
	ctx := context.Background()
	_, err = testMetadataRepo.SaveMetadata(ctx, "test data", metadataVal)

	if err.Error() != status.Error(codes.Internal, "some error").Error() {
		t.Error("get value error expected")
	}
}

func Test_metadataRepository_SaveMetadata_PutValueError(t *testing.T) {
	mockClient := new(etcd.ClientMock)
	metadataVal := &protofiles.Metadata{
		Author:      "littlestar642",
		ImageName:   "demo image",
		Name:        "test data",
		Description: "sample test metadata",
	}
	val, err := proto.Marshal(metadataVal)
	if err != nil {
		t.Error("error in marshalling metadata")
	}
	mockClient.On("PutValue", "metadata/test data", string(val)).Return(errors.New("some error"))
	mockClient.On("GetValue", "metadata/test data").Return("", nil)

	testMetadataRepo := NewMetadataRepository(mockClient)
	ctx := context.Background()
	_, err = testMetadataRepo.SaveMetadata(ctx, "test data", metadataVal)

	if err.Error() != status.Error(codes.Internal, "some error").Error() {
		t.Error("put value error expected")
	}
}

func TestGetMetadata(t *testing.T) {
	mockClient := new(etcd.ClientMock)
	testMetadataRepo := NewMetadataRepository(mockClient)
	jobName := "testJobName"
	key := "metadata/" + jobName

	var testMetadata = &protofiles.Metadata{
		Name:        "testJobName",
		Description: "This is a test image",
		ImageName:   "images/test-image",
	}

	str, err := proto.Marshal(testMetadata)
	if err != nil {
		log.Error(err, "error in test data marshalling")
	}
	mockClient.On("GetValue", key).Return(string(str), nil)
	resultMetadata, getValueErr := testMetadataRepo.GetMetadata(context.Background(), jobName)
	assert.Equal(t, resultMetadata.Name, testMetadata.Name)
	assert.Equal(t, resultMetadata.ImageName, testMetadata.ImageName)
	assert.Equal(t, resultMetadata.Description, testMetadata.Description)
	assert.Nil(t, getValueErr)
	mockClient.AssertExpectations(t)
}

func Test_metadataRepository_GetAvailableJobs(t *testing.T) {
	mockClient := new(etcd.ClientMock)
	var jobList []string
	jobList = append(jobList, "demo-image-name")
	jobList = append(jobList, "demo-image-name-1")

	testResponse := &protofiles.JobList{
		Jobs: jobList,
	}

	var keys []string
	keys = append(keys, "metadata/demo-image-name")
	keys = append(keys, "metadata/demo-image-name-1")

	var values []string

	mockClient.On("GetAllKeyAndValues", "metadata/").Return(keys, values, nil)

	testMetadataRepo := NewMetadataRepository(mockClient)
	ctx := context.Background()
	res, err := testMetadataRepo.GetAvailableJobs(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, res, testResponse)
	mockClient.AssertExpectations(t)
}

func Test_metadataRepository_GetAvailableJobs_ForEtcdClientFailure(t *testing.T) {
	mockClient := new(etcd.ClientMock)

	var keys []string
	var values []string

	mockClient.On("GetAllKeyAndValues", "metadata/").Return(keys, values, errors.New("error in etcd"))

	testMetadataRepo := NewMetadataRepository(mockClient)
	ctx := context.Background()
	_, err := testMetadataRepo.GetAvailableJobs(ctx)
	assert.Equal(t, status.Error(codes.Internal, "error in etcd").Error(), err.Error())

	mockClient.AssertExpectations(t)
}
