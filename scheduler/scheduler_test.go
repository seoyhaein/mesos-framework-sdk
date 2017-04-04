package scheduler

import (
	"mesos-framework-sdk/client"
	"mesos-framework-sdk/include/mesos"
	"net/http"
	"testing"
)

type mockClient struct{}

func (m *mockClient) Request(interface{}) (*http.Response, error) {
	return new(http.Response), nil
}

func (m *mockClient) StreamID() string {
	return "test"
}

func (m *mockClient) SetStreamID(string) client.Client {
	return m
}

type mockLogger struct{}

func (m *mockLogger) Emit(severity uint8, template string, args ...interface{}) {

}

var c = new(mockClient)
var i = &mesos_v1.FrameworkInfo{}
var l = new(mockLogger)

// Checks the internal state of a new scheduler.
func TestNewDefaultScheduler(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	if s.Client != c || s.logger != l || s.Info != i {
		t.Fatal("Scheduler does not have the right internal state")
	}
}

// Measures performance of creating a new scheduler.
func BenchmarkNewDefaultScheduler(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewDefaultScheduler(c, i, l)
	}
}

// Tests our accept call to Mesos.
func TestDefaultScheduler_Accept(t *testing.T) {
	s := NewDefaultScheduler(c, i, l)
	offerIds := []*mesos_v1.OfferID{}
	tasks := []*mesos_v1.Offer_Operation{}
	filters := &mesos_v1.Filters{}

	_, err := s.Accept(offerIds, tasks, filters)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our accept call to Mesos.
func BenchmarkDefaultScheduler_Accept(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	offerIds := []*mesos_v1.OfferID{}
	tasks := []*mesos_v1.Offer_Operation{}
	filters := &mesos_v1.Filters{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.Accept(offerIds, tasks, filters)
	}
}
