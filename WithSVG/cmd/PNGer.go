package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shomali11/gridder"
	"image/color"
	"io"
	"log"
	"net/http"
	"os"
)

func hellopng (c* gin.Context){
	imageConfig := gridder.ImageConfig{
		Width:  1000,
		Height: 1000,
		Name:   "example4.png",
	}
	gridConfig := gridder.GridConfig{
		Rows:              100,
		Columns:           100,
		LineStrokeWidth:   0.1,
		BorderStrokeWidth: 0.1,
	}

	grid, err := gridder.New(imageConfig, gridConfig)
	if err != nil {
		log.Fatal(err)
	}

	grid.DrawRectangle(0, 0, gridder.RectangleConfig{Width: 8, Height: 8, Color: color.Black, Stroke: true, Rotate: 45})
	grid.DrawRectangle(3, 0, gridder.RectangleConfig{Width: 8, Height: 8, Color: color.Black, Stroke: true, Rotate: 45, Dashes: 10})
	grid.DrawRectangle(0, 3, gridder.RectangleConfig{Width: 8, Height: 8, Color: color.Black, Stroke: true, StrokeWidth: 25})
	grid.DrawRectangle(2, 1, gridder.RectangleConfig{Width: 12, Height: 12, Color: color.RGBA{R: 255 / 2, A: 255 / 2}})
	grid.DrawRectangle(3, 3, gridder.RectangleConfig{Width: 8, Height: 8, Color: color.Black, Stroke: false})
	grid.SavePNG()

	if rdr, err:=os.Open(imageConfig.Name); err==nil {
		c.Status(http.StatusOK)
		io.Copy(c.Writer, rdr)
	}	else{
		svgError(c)
	}
}
