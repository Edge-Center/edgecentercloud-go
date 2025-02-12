package edgecloud

import (
	"encoding/json"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserActions_ListLogSubscriptions(t *testing.T) {
	setup()
	defer teardown()

	ls := LogSubscription{
		AuthHeaderName:  "Authorization",
		AuthHeaderValue: "Bearer YTH_s67j7xyWlFLy093RxReT5PmitnawLr25Jh7Ix14",
		URL:             "https://your-url.com/receive-user-action-messages",
		ID:              17,
	}

	expectedResp := LogSubscriptions{
		Count:   1,
		Results: []LogSubscription{ls},
	}

	URL := path.Join(userActionsBasePathV1, "/", listLogSubscriptions)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, err = w.Write(resp)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	respActual, resp, err := client.UserActions.ListLogSubscriptions(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestUserActions_ListAMQPSubscriptions(t *testing.T) {
	setup()
	defer teardown()

	amqps := AMQPSubscription{
		ConnectionString:         "amqps://guest:guest@192.168.123.20:5671/user_action_events",
		ReceiveChildClientEvents: false,
		RoutingKey:               "foo",
		ID:                       17,
	}

	expectedResp := AMQPSubscriptions{
		Count:   1,
		Results: []AMQPSubscription{amqps},
	}

	URL := path.Join(userActionsBasePathV1, "/", listAMQPSubscriptions)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, err = w.Write(resp)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	respActual, resp, err := client.UserActions.ListAMQPSubscriptions(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestUserActions_SubscribeLog(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(userActionsBasePathV1, "/", subscribeLog)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(204)
	})

	resp, err := client.UserActions.SubscribeLog(ctx, &LogSubscriptionCreateRequest{})
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 204)
}

func TestUserActions_UnsubscribeLog(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(userActionsBasePathV1, "/", unsubscribeLog)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(204)
	})

	resp, err := client.UserActions.UnsubscribeLog(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 204)
}

func TestUserActions_SubscribeAMQP(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(userActionsBasePathV1, "/", subscribeAMQP)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(204)
	})

	resp, err := client.UserActions.SubscribeAMQP(ctx, &AMQPSubscriptionCreateRequest{})
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 204)
}

func TestUserActions_UnsubscribeAMQP(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(userActionsBasePathV1, "/", unsubscribeAMQP)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(204)
	})

	resp, err := client.UserActions.UnsubscribeAMQP(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 204)
}
