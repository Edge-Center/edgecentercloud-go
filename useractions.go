package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	userActionsBasePathV1 = "/v1/user_actions"
)

const (
	listLogSubscriptions  = "subscriptions_list"
	listAMQPSubscriptions = "amqp_subscriptions_list"
	subscribeLog          = "subscribe"
	unsubscribeLog        = "unsubscribe"
	subscribeAMQP         = "subscribe_amqp"
	unsubscribeAMQP       = "unsubscribe_amqp"
)

// UsersService is an interface for provides access to user action logs and client subscription lists with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.online/cloud#tag/User-actions
type UserActionsService interface {
	ListLogSubscriptions(context.Context) (*LogSubscriptions, *Response, error)
	ListAMQPSubscriptions(context.Context) (*AMQPSubscriptions, *Response, error)
	SubscribeLog(context.Context, *LogSubscriptionCreateRequest) (*Response, error)
	UnsubscribeLog(context.Context) (*Response, error)
	SubscribeAMQP(context.Context, *AMQPSubscriptionCreateRequest) (*Response, error)
	UnsubscribeAMQP(context.Context) (*Response, error)
}

// UserActionsServiceOp handles communication with User Actions methods of the EdgecenterCloud API.
type UserActionsServiceOp struct {
	client *Client
}

var _ UserActionsService = &UserActionsServiceOp{}

// LogSubscriptions represents an EdgecenterCloud user action  log subscriptions.
type LogSubscriptions struct {
	Count   int               `json:"count"`
	Results []LogSubscription `json:"results"`
}

// LogSubscription represents an EdgecenterCloud user action log subscription result.
type LogSubscription struct {
	ID              int    `json:"id"`
	URL             string `json:"url"`
	AuthHeaderName  string `json:"auth_header_name"`
	AuthHeaderValue string `json:"auth_header_value"`
}

// LogSubscriptionCreateRequest represents a request to subscribe for a user action  log subscription.
type LogSubscriptionCreateRequest struct {
	URL             string `json:"url"`
	AuthHeaderName  string `json:"auth_header_name"`
	AuthHeaderValue string `json:"auth_header_value"`
}

// AMQPSubscriptions represents an EdgecenterCloud AMQP subscriptions.
type AMQPSubscriptions struct {
	Count   int                `json:"count"`
	Results []AMQPSubscription `json:"results"`
}

// AMQPSubscription represents an EdgecenterCloud AMQP subscription result.
type AMQPSubscription struct {
	ID                       int    `json:"id"`
	ConnectionString         string `json:"connection_string"`
	ReceiveChildClientEvents bool   `json:"receive_child_client_events"`
	RoutingKey               string `json:"routing_key"`
	Exchange                 string `json:"exchange"`
}

// AMQPSubscriptionCreateRequest represents a request to subscribe for a AMQP subscription.
type AMQPSubscriptionCreateRequest struct {
	ConnectionString         string  `json:"connection_string"`
	ReceiveChildClientEvents bool    `json:"receive_child_client_events"`
	RoutingKey               *string `json:"routing_key"`
	Exchange                 *string `json:"exchange"`
}

// ListLogSubscriptions get a list of user action log subscriptions.
func (s *UserActionsServiceOp) ListLogSubscriptions(ctx context.Context) (*LogSubscriptions, *Response, error) {
	pathReq := fmt.Sprintf("%s/%s", userActionsBasePathV1, listLogSubscriptions)

	req, err := s.client.NewRequest(ctx, http.MethodGet, pathReq, nil)
	if err != nil {
		return nil, nil, err
	}

	logSubs := new(LogSubscriptions)
	resp, err := s.client.Do(ctx, req, logSubs)
	if err != nil {
		return nil, resp, err
	}

	return logSubs, resp, nil
}

// ListAMQPSubscriptions get a list of AMQP user subscriptions.
func (s *UserActionsServiceOp) ListAMQPSubscriptions(ctx context.Context) (*AMQPSubscriptions, *Response, error) {
	pathReq := fmt.Sprintf("%s/%s", userActionsBasePathV1, listAMQPSubscriptions)

	req, err := s.client.NewRequest(ctx, http.MethodGet, pathReq, nil)
	if err != nil {
		return nil, nil, err
	}

	amqpSubs := new(AMQPSubscriptions)
	resp, err := s.client.Do(ctx, req, amqpSubs)
	if err != nil {
		return nil, resp, err
	}

	return amqpSubs, resp, nil
}

// SubscribeLog subscribe to the user action log. Subscription is created for the current client.
func (s *UserActionsServiceOp) SubscribeLog(ctx context.Context, reqBody *LogSubscriptionCreateRequest) (*Response, error) {
	if reqBody == nil {
		return nil, NewArgError("reqBody", "cannot be nil")
	}

	pathReq := fmt.Sprintf("%s/%s", userActionsBasePathV1, subscribeLog)

	req, err := s.client.NewRequest(ctx, http.MethodPost, pathReq, reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UnsubscribeLog unsubscribe from a user action log.
func (s *UserActionsServiceOp) UnsubscribeLog(ctx context.Context) (*Response, error) {
	pathReq := fmt.Sprintf("%s/%s", userActionsBasePathV1, unsubscribeLog)

	req, err := s.client.NewRequest(ctx, http.MethodPost, pathReq, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SubscribeAMQP subscribe to the user action log over AMQP. Subscription is created for the current client.
func (s *UserActionsServiceOp) SubscribeAMQP(ctx context.Context, reqBody *AMQPSubscriptionCreateRequest) (*Response, error) {
	if reqBody == nil {
		return nil, NewArgError("reqBody", "cannot be nil")
	}

	pathReq := fmt.Sprintf("%s/%s", userActionsBasePathV1, subscribeAMQP)

	req, err := s.client.NewRequest(ctx, http.MethodPost, pathReq, reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UnsubscribeAMQP unsubscribe from the user action log over AMQP.
func (s *UserActionsServiceOp) UnsubscribeAMQP(ctx context.Context) (*Response, error) {
	pathReq := fmt.Sprintf("%s/%s", userActionsBasePathV1, unsubscribeAMQP)

	req, err := s.client.NewRequest(ctx, http.MethodPost, pathReq, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
