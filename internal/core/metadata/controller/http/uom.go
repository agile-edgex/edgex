//
// Copyright (C) 2022-2023 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"

	bootstrapContainer "github.com/agile-edge/go-mod-bootstrap/v3/bootstrap/container"
	"github.com/agile-edge/go-mod-bootstrap/v3/di"
	"github.com/agile-edge/go-mod-core-contracts/v3/common"
	"github.com/agile-edge/go-mod-core-contracts/v3/dtos/responses"

	"github.com/agile-edge/edgex/internal/core/metadata/container"
	"github.com/agile-edge/edgex/internal/pkg"
	"github.com/agile-edge/edgex/internal/pkg/utils"

	"github.com/labstack/echo/v4"
)

type UnitOfMeasureController struct {
	dic *di.Container
}

func NewUnitOfMeasureController(dic *di.Container) *UnitOfMeasureController {
	return &UnitOfMeasureController{
		dic: dic,
	}
}

func (uc *UnitOfMeasureController) UnitsOfMeasure(c echo.Context) error {
	r := c.Request()
	w := c.Response()
	ctx := r.Context()
	u := container.UnitsOfMeasureFrom(uc.dic.Get)
	lc := bootstrapContainer.LoggingClientFrom(uc.dic.Get)

	response := responses.NewUnitsOfMeasureResponse("", "", http.StatusOK, u)

	utils.WriteHttpHeader(w, ctx, http.StatusOK)

	switch r.Header.Get(common.Accept) {
	case common.ContentTypeYAML:
		return pkg.EncodeAndWriteYamlResponse(u, w, lc)
	default:
		return pkg.EncodeAndWriteResponse(response, w, lc)
	}
}
