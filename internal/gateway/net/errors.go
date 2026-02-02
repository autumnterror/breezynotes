package net

import (
	"net/http"

	"github.com/autumnterror/breezynotes/internal/gateway/domain"
	"github.com/autumnterror/utils_go/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func bNErrors(op string, err error) (int, domain.Error) {
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "bad status.FromError", err)
			return http.StatusBadGateway, domain.Error{Error: "bad status. Check logs"}
		}

		switch st.Code() {
		case codes.Unauthenticated:
			return http.StatusUnauthorized, domain.Error{Error: "u dont have permission to this"}
		case codes.NotFound:
			return http.StatusNotFound, domain.Error{Error: "not found"}
		case codes.FailedPrecondition:
			return http.StatusFailedDependency, domain.Error{Error: "type do not register"}
		case codes.PermissionDenied:
			return http.StatusLocked, domain.Error{Error: "block already in use"}
		case codes.InvalidArgument:
			return http.StatusBadRequest, domain.Error{Error: st.Message()}
		case codes.Internal:
			return http.StatusBadGateway, domain.Error{Error: "check logs on service"}
		case codes.DeadlineExceeded:
			return http.StatusGatewayTimeout, domain.Error{Error: "response get too much time"}
		default:
			return http.StatusBadGateway, domain.Error{Error: "check logs on service"}
		}
	}
	return http.StatusOK, domain.Error{}
}

func authErrors(op string, err error) (int, domain.Error) {
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error(op, "bad status.FromError", err)
			return http.StatusBadGateway, domain.Error{Error: "bad status. Check logs"}
		}

		switch st.Code() {
		case codes.Unauthenticated:
			return http.StatusUnauthorized, domain.Error{Error: "token or id bad"}
		case codes.NotFound:
			return http.StatusNotFound, domain.Error{Error: "not found"}
		case codes.AlreadyExists:
			return http.StatusFound, domain.Error{Error: "already exist"}
		case codes.FailedPrecondition:
			return http.StatusFailedDependency, domain.Error{Error: "bad foreign key"}
		case codes.InvalidArgument:
			return http.StatusBadRequest, domain.Error{Error: st.Message()}
		case codes.ResourceExhausted:
			return http.StatusUnauthorized, domain.Error{Error: "terminate auth. Refresh token is expired"}
		case codes.Internal:
			return http.StatusBadGateway, domain.Error{Error: "check logs on service"}
		case codes.DeadlineExceeded:
			return http.StatusGatewayTimeout, domain.Error{Error: "response get too much time"}
		default:
			return http.StatusBadGateway, domain.Error{Error: "check logs on service"}
		}
	}
	return http.StatusOK, domain.Error{}
}
