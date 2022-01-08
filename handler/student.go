package handler

import (
	"fmt"
	"myapp/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	Area struct {
		ID        int64  `gorm:"column:id;primary_key"`
		AreaValue int64  `gorm:"column:area_value"`
		AreaType  string `gorm:"column:type"`
	}
)

func (h *Handler) Create(c echo.Context) (err error) {
	// Bind
	u := &model.Student{}
	if err = c.Bind(u); err != nil {
		return
	}

	errr := h.InsertArea(20, 20, "persegi")
	if errr != nil {
		fmt.Println(errr.Error())
	}

	h.GormDB.Save(u)
	return c.JSON(http.StatusCreated, u)
}

// type => types
// []string => string
// param1, param2 dan area disamakan tipe datanya
// Var => var
// declare var area di switch slah
// 'persegi panjang' => "persegi panjang"

func (_r *Handler) InsertArea(param1 int64, param2 int64, types string) (err error) {

	ar := &Area{}
	// var ar *Area
	// inst := _r.GormDB.Model(ar)

	// fmt.Println(inst.Commit())

	// var area int64
	// area = 0
	// fmt.Println(area)
	switch types {
	case "persegi panjang":
		var area int64 = param1 * param2
		ar.AreaValue = area
		ar.AreaType = "persegi panjang"
		_r.GormDB.Save(ar)

	case "persegi":
		var area int64 = param1 * param2
		ar.AreaValue = area
		ar.AreaType = "persegi"
		err = _r.GormDB.Create(&ar).Error
		if err != nil {
			return err
		}

	case "segitiga":
		var area int64 = (param1 * param2) / 2
		ar.AreaValue = area
		ar.AreaType = "segitiga"
		_r.GormDB.Save(ar)

	default:
		ar.AreaValue = 0
		ar.AreaType = "undefined data"
		err = _r.GormDB.Create(&ar).Error
		if err != nil {
			return err
		}
	}
	return nil
}
