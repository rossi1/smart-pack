package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rossi1/smart-pack/app/command"
	"github.com/rossi1/smart-pack/app/query"
	"github.com/rossi1/smart-pack/domain"
	"github.com/rossi1/smart-pack/pkg/server/dto"
	"github.com/rossi1/smart-pack/pkg/server/httperr"
	"github.com/rossi1/smart-pack/ports"
	"github.com/sirupsen/logrus"
)

func (h *SetPackSizesRequestValidation) valid() *httperr.ErrorMessageBody {
	validate := validator.New()
	validation := validate.Struct(h)
	if validation != nil {
		validationErrors := validator.ValidationErrors{}
		_ = errors.As(validation, &validationErrors)
		return httperr.NewErrorMessageBodyWithMessages(
			httperr.NewErrorMessage(
				domain.ErrorInvalidRequestBodyParameter,
				validationErrors[0].Field(),
			))
	}
	return nil
}
func (s *HTTPServer) CalculatePacks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ports.CalculateRequest
	if err := dto.Read(r, &req); err != nil {
		logrus.WithError(err).Error("Failed to read request body")
		httperr.UnprocessableEntity(domain.ErrorUnprocessableEntityLabel, "Invalid request body", err, w, r)
		return
	}

	if req.ItemsOrdered <= 0 {
		httperr.BadRequest(domain.ErrorBadRequestLabel, "items_ordered must be positive", nil, w, r)
		return
	}

	packSizes, err := s.app.Queries.GetPackSizes.Handle(ctx, &query.GetPackSizesQuery{})

	if v, ok := domain.IsHTTPCustomError(err); ok {
		httperr.WithStatus(v.Label(), "", err, w, r, v.Status())
		return
	}

	if err != nil {
		logrus.WithError(err).Error("Failed to get pack sizes")
		httperr.InternalError(domain.ErrorInternalServerErrorLabel, "", err, w, r)
		return
	}

	packSizesInt := make([]int, len(packSizes))
	for i, size := range packSizes {
		packSizesInt[i] = size.Size
	}

	result, err := s.app.PackCalculator.Calculate(req.ItemsOrdered, packSizesInt)

	if err != nil {
		logrus.WithError(err).Error("Failed to calculate packs")
		httperr.InternalError(domain.ErrorInternalServerErrorLabel, "", err, w, r)
		return
	}

	resp := mapDomainToPortsPackSolution(result)
	dto.Write(w, r, resp)
}

func (s *HTTPServer) GetPackSizes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sizes, err := s.app.Queries.GetPackSizes.Handle(ctx, &query.GetPackSizesQuery{})
	if v, ok := domain.IsHTTPCustomError(err); ok {
		logrus.WithError(err).Error("Failed to get pack sizes")
		httperr.WithStatus(v.Label(), "", err, w, r, v.Status())
		return
	}

	if err != nil {
		logrus.WithError(err).Error("Failed to get pack sizes")
		httperr.InternalError(domain.ErrorInternalServerErrorLabel, "", err, w, r)
		return
	}

	resp := mapSmartPackToPackSizesResponse(sizes)

	dto.Write(w, r, resp)
}

func (s *HTTPServer) SetPackSizes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req SetPackSizesRequestValidation
	if err := dto.Read(r, &req); err != nil {
		httperr.UnprocessableEntity(domain.ErrorUnprocessableEntityLabel, "", err, w, r)
		return
	}
	if validationErr := req.valid(); validationErr != nil {
		logrus.WithContext(ctx).Error(validationErr)
		httperr.BadRequest(validationErr.Messages[0].Label, validationErr.Messages[0].FormProperty, nil, w, r)
		return
	}

	cmd := command.SetPackSizesCommand{
		Sizes: mapToSmartPack(req),
	}
	if err := s.app.Commands.SetPackSizes.Handle(ctx, &cmd); err != nil {
		logrus.WithError(err).Error("Failed to set pack sizes")
		httperr.InternalError(domain.ErrorInternalServerErrorLabel, "", err, w, r)
		return
	}

	dto.Write(w, r, http.StatusOK)
}

func mapToSmartPack(req SetPackSizesRequestValidation) []domain.SmartPack {
	sizes := make([]domain.SmartPack, 0, len(req.PackSizes))
	for _, size := range req.PackSizes {
		sizes = append(sizes, domain.SmartPack{Size: size})
	}
	return sizes
}

func mapSmartPackToPackSizesResponse(smartPack []domain.SmartPack) ports.PackSizesResponse {
	sizes := make([]int, len(smartPack))
	for i, size := range smartPack {
		sizes[i] = size.Size
	}
	return ports.PackSizesResponse{
		PackSizes: sizes,
	}
}

func mapDomainToPortsPackSolution(result *domain.PackSolution) ports.PackSolution {
	resp := ports.PackSolution{
		ItemsOrdered: result.ItemsOrdered,
		TotalItems:   result.TotalItems,
		TotalPacks:   result.TotalPacks,
		Packs:        make(map[string]int),
		PackDetails:  make([]ports.PackDetail, 0, len(result.PackDetails)),
	}

	for size, qty := range result.Packs {
		resp.Packs[fmt.Sprintf("%d", size)] = qty
	}

	for _, d := range result.PackDetails {
		resp.PackDetails = append(resp.PackDetails, ports.PackDetail{
			Size:     d.Size,
			Quantity: d.Quantity,
		})
	}

	return resp
}
