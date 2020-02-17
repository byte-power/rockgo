package log

import (
	"testing"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/stretchr/testify/mock"
)

type fluentMock struct {
	mock.Mock
}

func (f *fluentMock) Post(tag string, message interface{}) error {
	f.Called(tag, message)
	return nil
}

func TestFluentLogFormat(t *testing.T) {
	logger := FluentOutput{level: LevelInfo, tag: "test-fluent", config: fluent.Config{
		FluentHost: "127.0.0.1",
		FluentPort: 24225,
		Async: false,
	}}
	testObj := new(fluentMock)
	logger.output = testObj
	testObj.On("Post", "test-fluent", map[string]interface{}{"message": "haha"}).Return(nil)
	logger.Log(LevelInfo, "haha", nil)
	testObj.AssertExpectations(t)
}
