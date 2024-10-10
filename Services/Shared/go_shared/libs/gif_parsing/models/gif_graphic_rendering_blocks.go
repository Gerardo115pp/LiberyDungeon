package models

import (
	"fmt"
	"image"
	"image/color"
)

type GifGraphicRenderingBlock struct {
	GraphicControlExtension *GifGraphicControlExtensionBlock `json:"graphic_control_extension"`
	ImageDescriptor         *GifImageDescriptor              `json:"image_descriptor"`
	LocalColorTable         *GifColorTable                   `json:"local_color_table"`
	ImageData               *GifTableBasedImageData          `json:"image_data"`
	FileFinalPosition       int64                            `json:"file_final_position"`
	fileStartPosition       int64
}

func (ggrb *GifGraphicRenderingBlock) SetFileStartPosition(start_offset int64) {
	if ggrb.fileStartPosition == 0 {
		ggrb.fileStartPosition = start_offset
	}
}

func (ggrb *GifGraphicRenderingBlock) ToImage(gif_color_table *GifColorTable) (*image.Paletted, error) {
	var rgba_color_palette []color.Color = gif_color_table.ToColorPalette()
	var image_decompressed_data []byte
	var err error

	image_decompressed_data, err = ggrb.ImageData.Decompress()
	if err != nil {
		fmt.Printf("Error while decompressing image data: %s\n", err.Error())
		return nil, err
	}

	var int_width int = int(ggrb.ImageDescriptor.ImageWidth)
	var int_height int = int(ggrb.ImageDescriptor.ImageHeight)

	var palletted_image *image.Paletted
	palletted_image = image.NewPaletted(image.Rect(0, 0, int_width, int_height), rgba_color_palette)

	for h, idx := range image_decompressed_data {
		x := h % int_width
		y := h / int_width

		palletted_image.SetColorIndex(x, y, idx)
	}

	return palletted_image, nil
}

func (ggrb GifGraphicRenderingBlock) String() string {
	var graphic_rendering_block_string string = ""

	if ggrb.GraphicControlExtension != nil {
		graphic_rendering_block_string += fmt.Sprintf("-> %s\n", ggrb.GraphicControlExtension)
	}

	graphic_rendering_block_string += fmt.Sprintf("-> %s\n", ggrb.ImageDescriptor)

	if ggrb.LocalColorTable != nil {
		graphic_rendering_block_string += fmt.Sprintf("Local color table:\n%s\n", ggrb.LocalColorTable)
	}

	if ggrb.ImageData == nil {
		graphic_rendering_block_string += "Image data: There is likely a problem. No image data found."
	} else {
		graphic_rendering_block_string += fmt.Sprintf("-> %s\n", ggrb.ImageData)
	}

	graphic_rendering_block_string += fmt.Sprintf("\n\n\nGraphic rendering block: %#x - %#x\n\n======================================\n", ggrb.fileStartPosition, ggrb.FileFinalPosition)

	return graphic_rendering_block_string
}
