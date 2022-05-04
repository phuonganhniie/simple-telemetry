package handler

import (
	"fmt"
	"net/http"

	pb_checkout "github.com/phuonganhniie/simpleTelemetry/proto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func CheckOutHandler(c pb_checkout.CheckoutClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow only POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Create a tracer span
		tr := otel.Tracer("http")
		ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		defer span.End()

		// Make the GRPC call to checkout-service
		_, err := c.DoCheckOut(ctx, &pb_checkout.CheckOutRequest{
			ItemsID: []int32{1, 2, 3, 4},
		})

		// Check for errors
		rStatus := status.Convert(err)
		if rStatus != nil {
			span.SetStatus(codes.Error, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
