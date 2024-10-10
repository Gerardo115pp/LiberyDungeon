package models

import (
	"bytes"
	"compress/lzw"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"

	"golang.org/x/image/draw"
)

type GifBlockLabel byte

const (
	GIF_TRAILER_LABEL          GifBlockLabel = 0x3B
	GIF_IMAGE_DESCRIPTOR_LABEL GifBlockLabel = 0x2C
)

/*
Represent an entire parsed gif file.
*/
type ParsedGif struct {
	Filename                string
	Header                  GifHeader
	LogicalScreenDescriptor GifLogicalScreenDescriptor
	GlobalColorTable        *GifColorTable
	NoScopeExtensions       []GifExtensionBlock
	GraphicRenderingBlocks  []GifGraphicRenderingBlock
	TrailerFilePosition     int64
}

func (pg ParsedGif) String() string {
	var gif_string string = fmt.Sprintf("Filename: %s\n\n", pg.Filename)

	gif_string += fmt.Sprintf("Gif file signature: %s\n", pg.Header)
	gif_string += pg.LogicalScreenDescriptor.String()

	if pg.LogicalScreenDescriptor.HasGlobalColorTable {
		gif_string += fmt.Sprintf("\nGlobal color table color count: %d\n", pg.LogicalScreenDescriptor.GetGlobalColorTableByteCount()/3)
		gif_string += fmt.Sprintf("Global color table: %s\n", pg.GlobalColorTable)
	}

	for h, extension_block := range pg.NoScopeExtensions {
		gif_string += fmt.Sprintf("Extension index %d:\n%s\n", h, extension_block)
	}

	for h, graphic_rendering_block := range pg.GraphicRenderingBlocks {
		gif_string += fmt.Sprintf("Graphic rendering block index %d:\n%s\n", h, graphic_rendering_block)
	}

	gif_string += fmt.Sprintf("\n\n\n\nTrailer position: %#x\n", pg.TrailerFilePosition)

	return gif_string
}

func (pg ParsedGif) GetFrameCount() int {
	return len(pg.GraphicRenderingBlocks)
}

func (pg ParsedGif) GetCommentBlocks() []GifCommentExtensionBlock {
	var comment_blocks []GifCommentExtensionBlock = make([]GifCommentExtensionBlock, 0)

	for _, extension_block := range pg.NoScopeExtensions {
		if comment_block, is_comment_block := extension_block.(GifCommentExtensionBlock); is_comment_block {
			comment_blocks = append(comment_blocks, comment_block)
		}
	}

	return comment_blocks
}

func (pg ParsedGif) RenderFrame(frame_index int) (image.Image, error) {
	var selected_frame GifGraphicRenderingBlock = pg.GraphicRenderingBlocks[frame_index]
	var authoritative_color_table *GifColorTable = pg.GlobalColorTable
	var err error

	if selected_frame.ImageDescriptor.HasLocalColorTable {
		authoritative_color_table = selected_frame.LocalColorTable
	}

	var frame_image *image.Paletted
	frame_image, err = selected_frame.ToImage(authoritative_color_table)
	if err != nil {
		return nil, fmt.Errorf("Error while rendering frame: %s", err.Error())
	}

	return frame_image, nil
}

// Render a frame with a specific width. This function under no circumstances will upscale the image or change the aspect ratio.
func (pg ParsedGif) RenderResizeFrame(frame_index int, width int) (image.Image, error) {
	var selected_frame GifGraphicRenderingBlock = pg.GraphicRenderingBlocks[frame_index]
	var frame image.Image
	var err error

	var frame_width float64 = float64(selected_frame.ImageDescriptor.ImageWidth)
	var target_width float64 = float64(width)

	if selected_frame.ImageDescriptor.ImageWidth > uint16(width) {
		var scale_down_factor float64 = target_width / frame_width

		frame, err = pg.renderScaledDownFrame(frame_index, scale_down_factor)
		if err != nil {
			return nil, fmt.Errorf("Error while rendering frame: %s", err.Error())
		}
	} else {
		frame, err = pg.RenderFrame(frame_index)
		if err != nil {
			return nil, fmt.Errorf("Error while rendering frame: %s", err.Error())
		}
	}

	return frame, nil
}

func (pg ParsedGif) renderScaledDownFrame(frame_index int, size_percentage float64) (image.Image, error) {

	var frame_image, err = pg.RenderFrame(frame_index)
	if err != nil {
		return nil, fmt.Errorf("Error while rendering frame: %s", err.Error())
	}

	var new_width int = int(float64(frame_image.Bounds().Dx()) * size_percentage)
	var new_height int = int(float64(frame_image.Bounds().Dy()) * size_percentage)

	scaled_down_image := image.NewRGBA(image.Rect(0, 0, new_width, new_height))
	draw.CatmullRom.Scale(scaled_down_image, scaled_down_image.Bounds(), frame_image, frame_image.Bounds(), draw.Over, nil)

	return scaled_down_image, nil
}

func NewParsedGif() *ParsedGif {
	var new_parsed_gif *ParsedGif = new(ParsedGif)

	new_parsed_gif.NoScopeExtensions = make([]GifExtensionBlock, 0)

	return new_parsed_gif
}

/*
 * ===========================================================
 *                      GIF HEADER
 * ===========================================================
 */

type GifHeader struct {
	Signature         string `json:"signature"`
	Version           string `json:"version"`
	FileFinalPosition int64  `json:"file_position"`
}

func (gh GifHeader) String() string {
	return fmt.Sprintf("%s%s\n\nHeader range: 0x0 - %#x\n\n", gh.Signature, gh.Version, gh.GetSize()-1)
}

func (gh GifHeader) GetSize() int {
	return 6
}

func NewGifHeader(header_fragment []byte) (*GifHeader, error) {
	if len(header_fragment) < 6 {
		return nil, fmt.Errorf("Invalid header fragment")
	}

	var new_gif_header *GifHeader = new(GifHeader)

	var supported_gif_signatures = [][]byte{
		{0x47, 0x49, 0x46, 0x38, 0x39, 0x61},
		{0x47, 0x49, 0x46, 0x38, 0x37, 0x61},
	}

	if !bytes.Equal(header_fragment, supported_gif_signatures[0]) && !bytes.Equal(header_fragment, supported_gif_signatures[1]) {
		return nil, fmt.Errorf("Unsupported format: '%s'", string(header_fragment))
	}

	new_gif_header.Signature = string(header_fragment[:3])
	new_gif_header.Version = string(header_fragment[3:6])

	return new_gif_header, nil
}

/*
* ===========================================================
*                      LOGICAL SCREEN DESCRIPTOR
* ===========================================================
 From the GIF89a specification:
	18. Logical Screen Descriptor.

		a. Description.  The Logical Screen Descriptor contains the parameters
		necessary to define the area of the display device within which the
		images will be rendered.  The coordinates in this block are given with
		respect to the top-left corner of the virtual screen; they do not
		necessarily refer to absolute coordinates on the display device.  This
		implies that they could refer to window coordinates in a window-based
		environment or printer coordinates when a printer is used.

		This block is REQUIRED; exactly one Logical Screen Descriptor must be
		present per Data Stream.

		b. Required Version.  Not applicable. This block is not subject to a
		version number. This block must appear immediately after the Header.

		c. Syntax.

		7 6 5 4 3 2 1 0        Field Name                    Type
		+---------------+
	0  |               |       Logical Screen Width          Unsigned
		+-             -+
	1  |               |
		+---------------+
	2  |               |       Logical Screen Height         Unsigned
		+-             -+
	3  |               |
		+---------------+
	4  | |     | |     |       <Packed Fields>               See below
		+---------------+
	5  |               |       Background Color Index        Byte
		+---------------+
	6  |               |       Pixel Aspect Ratio            Byte
		+---------------+

		<Packed Fields>  =      Global Color Table Flag       1 Bit
								Color Resolution              3 Bits
								Sort Flag                     1 Bit
								Size of Global Color Table    3 Bits

				i) Logical Screen Width - Width, in pixels, of the Logical Screen
				where the images will be rendered in the displaying device.

				ii) Logical Screen Height - Height, in pixels, of the Logical
				Screen where the images will be rendered in the displaying device.

				iii) Global Color Table Flag - Flag indicating the presence of a
				Global Color Table; if the flag is set, the Global Color Table will
				immediately follow the Logical Screen Descriptor. This flag also
				selects the interpretation of the Background Color Index; if the
				flag is set, the value of the Background Color Index field should
				be used as the table index of the background color. (This field is
				the most significant bit of the byte.)

				Values :    0 -   No Global Color Table follows, the Background
								Color Index field is meaningless.
							1 -   A Global Color Table will immediately follow, the
								Background Color Index field is meaningful.

				iv) Color Resolution - Number of bits per primary color available
				to the original image, minus 1. This value represents the size of
				the entire palette from which the colors in the graphic were
				selected, not the number of colors actually used in the graphic.
				For example, if the value in this field is 3, then the palette of
				the original image had 4 bits per primary color available to create
				the image.  This value should be set to indicate the richness of
				the original palette, even if not every color from the whole
				palette is available on the source machine.

				v) Sort Flag - Indicates whether the Global Color Table is sorted.
				If the flag is set, the Global Color Table is sorted, in order of
				decreasing importance. Typically, the order would be decreasing
				frequency, with most frequent color first. This assists a decoder,
				with fewer available colors, in choosing the best subset of colors;
				the decoder may use an initial segment of the table to render the
				graphic.

				Values :    0 -   Not ordered.
							1 -   Ordered by decreasing importance, most
								important color first.

				vi) Size of Global Color Table - If the Global Color Table Flag is
				set to 1, the value in this field is used to calculate the number
				of bytes contained in the Global Color Table. To determine that
				actual size of the color table, raise 2 to [the value of the field
				+ 1].  Even if there is no Global Color Table specified, set this
				field according to the above formula so that decoders can choose
				the best graphics mode to display the stream in.  (This field is
				made up of the 3 least significant bits of the byte.)

				vii) Background Color Index - Index into the Global Color Table for

				the Background Color. The Background Color is the color used for
				those pixels on the screen that are not covered by an image. If the
				Global Color Table Flag is set to (zero), this field should be zero
				and should be ignored.

				viii) Pixel Aspect Ratio - Factor used to compute an approximation
				of the aspect ratio of the pixel in the original image.  If the
				value of the field is not 0, this approximation of the aspect ratio
				is computed based on the formula:

				Aspect Ratio = (Pixel Aspect Ratio + 15) / 64

				The Pixel Aspect Ratio is defined to be the quotient of the pixel's
				width over its height.  The value range in this field allows
				specification of the widest pixel of 4:1 to the tallest pixel of
				1:4 in increments of 1/64th.

				Values :        0 -   No aspect ratio information is given.
						1..255 -   Value used in the computation.

		d. Extensions and Scope. The scope of this block is the entire Data
		Stream. This block cannot be modified by any extension.

		e. Recommendations. None.

	for the whole thing, see: https://www.w3.org/Graphics/GIF/spec-gif89a.txt
*/

type GifLogicalScreenDescriptor struct {
	LogicalScreenWidth     uint16 `json:"logical_screen_width"`
	LogicalScreenHeight    uint16 `json:"logical_screen_height"`
	HasGlobalColorTable    bool   `json:"has_global_color_table"`
	ColorResolution        uint8  `json:"color_resolution"`
	GlobalColorTableSorted bool   `json:"global_color_table_sorted"`
	GlobalColorTableSize   uint8  `json:"global_color_table_size"`
	BackgroundColorIndex   uint8  `json:"background_color_index"`
	PixelAspectRatio       uint8  `json:"pixel_aspect_ratio"`
	FileFinalPosition      int64  `json:"file_position"`
}

func (lsd GifLogicalScreenDescriptor) String() string {
	var start_position int64 = lsd.FileFinalPosition - int64(lsd.GetSize())

	return fmt.Sprintf("Width: %d\nHeight: %d\nHas Global Color Table: %t\nColor Resolution: %d\nGlobal Color Table Sorted: %t\nGlobal Color Table Size: %d\nBackground Color Index: %d\nPixel Aspect Ratio: %d\n\nLogical Screen Descriptor: %#x - %#x\n\n", lsd.LogicalScreenWidth, lsd.LogicalScreenHeight, lsd.HasGlobalColorTable, lsd.ColorResolution, lsd.GlobalColorTableSorted, lsd.GlobalColorTableSize, lsd.BackgroundColorIndex, lsd.PixelAspectRatio, start_position, lsd.FileFinalPosition-1)
}

func (lsd GifLogicalScreenDescriptor) GetSize() int {
	return 7
}

func (lsd GifLogicalScreenDescriptor) GetAspectRatio() float64 {
	if lsd.PixelAspectRatio == 0 {
		return 1.0
	}

	return float64(lsd.PixelAspectRatio+15) / 64.0
}

func (lsd GifLogicalScreenDescriptor) GetGlobalColorTableByteCount() int {
	if !lsd.HasGlobalColorTable {
		return 0
	}
	var byte_count int = 3 * int(math.Pow(2, float64(lsd.GlobalColorTableSize+1)))

	return byte_count
}

func NewGifLogicalScreenDescriptor(lsd_fragment []byte) (*GifLogicalScreenDescriptor, error) {
	if len(lsd_fragment) < 7 {
		return nil, fmt.Errorf("Invalid logical screen descriptor fragment")
	}

	var new_gif_logical_screen_descriptor *GifLogicalScreenDescriptor = new(GifLogicalScreenDescriptor)

	new_gif_logical_screen_descriptor.LogicalScreenWidth = uint16(lsd_fragment[0]) | (uint16(lsd_fragment[1]) << 8)
	new_gif_logical_screen_descriptor.LogicalScreenHeight = uint16(lsd_fragment[2]) | (uint16(lsd_fragment[3]) << 8)

	var packed_fields byte = lsd_fragment[4]

	new_gif_logical_screen_descriptor.HasGlobalColorTable = (packed_fields & 0x80) != 0
	new_gif_logical_screen_descriptor.ColorResolution = uint8(math.Pow(2, float64((((packed_fields & 0x70) >> 4) + 1))))
	new_gif_logical_screen_descriptor.GlobalColorTableSorted = (packed_fields & 0x08) != 0
	new_gif_logical_screen_descriptor.GlobalColorTableSize = packed_fields & 0x07

	new_gif_logical_screen_descriptor.BackgroundColorIndex = lsd_fragment[5]
	new_gif_logical_screen_descriptor.PixelAspectRatio = lsd_fragment[6]

	return new_gif_logical_screen_descriptor, nil
}

/*
- ===========================================================
- COLOR TABLE
- ===========================================================
From the GIF89a specification:

 19. Global Color Table.

    a. Description. This block contains a color table, which is a sequence of
    bytes representing red-green-blue color triplets. The Global Color Table
    is used by images without a Local Color Table and by Plain Text
    Extensions. Its presence is marked by the Global Color Table Flag being
    set to 1 in the Logical Screen Descriptor; if present, it immediately
    follows the Logical Screen Descriptor and contains a number of bytes
    equal to
    3 x 2^(Size of Global Color Table+1).

    This block is OPTIONAL; at most one Global Color Table may be present
    per Data Stream.

    b. Required Version.  87a

    11

    c. Syntax.

    7 6 5 4 3 2 1 0        Field Name                    Type
    +===============+
    0  |               |       Red 0                         Byte
    +-             -+
    1  |               |       Green 0                       Byte
    +-             -+
    2  |               |       Blue 0                        Byte
    +-             -+
    3  |               |       Red 1                         Byte
    +-             -+
    |               |       Green 1                       Byte
    +-             -+
    up  |               |
    +-   . . . .   -+       ...
    to  |               |
    +-             -+
    |               |       Green 255                     Byte
    +-             -+
    767  |               |       Blue 255                      Byte
    +===============+

    d. Extensions and Scope. The scope of this block is the entire Data
    Stream. This block cannot be modified by any extension.

    e. Recommendations. None.

    for the whole thing, see: https://www.w3.org/Graphics/GIF/spec-gif89a.txt
*/
type GifColorTable struct {
	Colors            []RGBColor `json:"colors"`
	FileFinalPosition int64      `json:"file_position"`
	FileStartPosition int64
}

func (gct GifColorTable) String() string {
	var colors_string string = ""

	for h, color := range gct.Colors {
		colors_string += color.String()

		if h != (len(gct.Colors) - 1) {
			colors_string += ", "
		} else {
			colors_string += "\n"
		}
	}

	var start_position int64 = gct.FileFinalPosition - int64(gct.GetSize())

	colors_string += fmt.Sprintf("\n\nColor Table: %#x - %#x\n\n", start_position, gct.FileFinalPosition-1)

	return colors_string
}

func (gct GifColorTable) GetSize() int {
	return len(gct.Colors) * 3 // 3 bytes per color
}

func (gct GifColorTable) ToColorPalette() []color.Color {
	var color_palette []color.Color

	for _, rgb_color := range gct.Colors {
		color_palette = append(color_palette, color.RGBA{rgb_color.Red, rgb_color.Green, rgb_color.Blue, 0xFF})
	}

	return color_palette
}

// Parses a sequence of bytes into a GifColorTable. It assumes that byte sequence only includes the color table data.
// Attempts to prevent errors by checking if the length of the byte sequence is a multiple of 3. But this is in no way
// a guarantee that the byte sequence is a valid color table. The caller is responsible for only passing the color table chunk.
func NewGifColorTable(color_table_fragment []byte) (*GifColorTable, error) {
	if len(color_table_fragment)%3 != 0 {
		return nil, fmt.Errorf("Invalid color table fragment")
	}

	var new_gif_color_table *GifColorTable = new(GifColorTable)
	new_gif_color_table.Colors = make([]RGBColor, 0)

	for h := 0; h < len(color_table_fragment); h += 3 {
		new_rgb_color, err := NewRGBColor(color_table_fragment[h : h+3])
		if err != nil {
			return nil, fmt.Errorf("Error while creating new RGB color: %s", err.Error())
		}

		new_gif_color_table.Colors = append(new_gif_color_table.Colors, *new_rgb_color)
	}

	return new_gif_color_table, nil
}

type RGBColor struct {
	Red   uint8 `json:"red"`
	Green uint8 `json:"green"`
	Blue  uint8 `json:"blue"`
}

func (rgb RGBColor) String() string {
	return fmt.Sprintf("(%d, %d, %d)", rgb.Red, rgb.Green, rgb.Blue)
}

func NewRGBColor(color_fragment []byte) (*RGBColor, error) {
	if len(color_fragment) != 3 {
		return nil, fmt.Errorf("Invalid color fragment")
	}

	var new_rgb_color *RGBColor = new(RGBColor)

	new_rgb_color.Red = color_fragment[0]
	new_rgb_color.Green = color_fragment[1]
	new_rgb_color.Blue = color_fragment[2]

	return new_rgb_color, nil

}

/*
- ===========================================================
-                      GIF IMAGE DESCRIPTOR
- ===========================================================
From the GIF89a specification:
20. Image Descriptor.

	      a. Description. Each image in the Data Stream is composed of an Image
	      Descriptor, an optional Local Color Table, and the image data.  Each
	      image must fit within the boundaries of the Logical Screen, as defined
	      in the Logical Screen Descriptor.

	      The Image Descriptor contains the parameters necessary to process a table
	      based image. The coordinates given in this block refer to coordinates
	      within the Logical Screen, and are given in pixels. This block is a
	      Graphic-Rendering Block, optionally preceded by one or more Control
	      blocks such as the Graphic Control Extension, and may be optionally
	      followed by a Local Color Table; the Image Descriptor is always followed
	      by the image data.

	      This block is REQUIRED for an image.  Exactly one Image Descriptor must
	      be present per image in the Data Stream.  An unlimited number of images
	      may be present per Data Stream.

	      b. Required Version.  87a.

	      c. Syntax.

	      7 6 5 4 3 2 1 0        Field Name                    Type
	     +---------------+
	  0  |               |       Image Separator               Byte
	     +---------------+
	  1  |               |       Image Left Position           Unsigned
	     +-             -+
	  2  |               |
	     +---------------+
	  3  |               |       Image Top Position            Unsigned
	     +-             -+
	  4  |               |
	     +---------------+
	  5  |               |       Image Width                   Unsigned
	     +-             -+
	  6  |               |
	     +---------------+
	  7  |               |       Image Height                  Unsigned
	     +-             -+
	  8  |               |
	     +---------------+
	  9  | | | |   |     |       <Packed Fields>               See below
	     +---------------+

	     <Packed Fields>  =      Local Color Table Flag        1 Bit
	                             Interlace Flag                1 Bit
	                             Sort Flag                     1 Bit
	                             Reserved                      2 Bits
	                             Size of Local Color Table     3 Bits

	           i) Image Separator - Identifies the beginning of an Image
	           Descriptor. This field contains the fixed value 0x2C.

	           ii) Image Left Position - Column number, in pixels, of the left edge
	           of the image, with respect to the left edge of the Logical Screen.
	           Leftmost column of the Logical Screen is 0.

	           iii) Image Top Position - Row number, in pixels, of the top edge of
	           the image with respect to the top edge of the Logical Screen. Top
	           row of the Logical Screen is 0.

	           iv) Image Width - Width of the image in pixels.

	           v) Image Height - Height of the image in pixels.

	           vi) Local Color Table Flag - Indicates the presence of a Local Color
	           Table immediately following this Image Descriptor. (This field is
	           the most significant bit of the byte.)


	           Values :    0 -   Local Color Table is not present. Use
	                             Global Color Table if available.
	                       1 -   Local Color Table present, and to follow
	                             immediately after this Image Descriptor.


	           vii) Interlace Flag - Indicates if the image is interlaced. An image
	           is interlaced in a four-pass interlace pattern; see Appendix E for
	           details.

	           Values :    0 - Image is not interlaced.
	                       1 - Image is interlaced.

	            viii) Sort Flag - Indicates whether the Local Color Table is
	            sorted.  If the flag is set, the Local Color Table is sorted, in
	            order of decreasing importance. Typically, the order would be
	            decreasing frequency, with most frequent color first. This assists
	            a decoder, with fewer available colors, in choosing the best subset
	            of colors; the decoder may use an initial segment of the table to
	            render the graphic.

	            Values :    0 -   Not ordered.
	                        1 -   Ordered by decreasing importance, most
	                              important color first.

	            ix) Size of Local Color Table - If the Local Color Table Flag is
	            set to 1, the value in this field is used to calculate the number
	            of bytes contained in the Local Color Table. To determine that
	            actual size of the color table, raise 2 to the value of the field
	            + 1. This value should be 0 if there is no Local Color Table
	            specified. (This field is made up of the 3 least significant bits
	            of the byte.)

	     d. Extensions and Scope. The scope of this block is the Table-based Image
	     Data Block that follows it. This block may be modified by the Graphic
	     Control Extension.

	     e. Recommendation. None.

		 for the whole thing, see: https://www.w3.org/Graphics/GIF/spec-gif89a.txt
*/
type GifImageDescriptor struct {
	ImageLeftPosition     uint16 `json:"image_left_position"`
	ImageTopPosition      uint16 `json:"image_top_position"`
	ImageWidth            uint16 `json:"image_width"`
	ImageHeight           uint16 `json:"image_height"`
	HasLocalColorTable    bool   `json:"has_local_color_table"`
	IsInterlaced          bool   `json:"is_interlaced"`
	LocalColorTableSorted bool   `json:"local_color_table_sorted"`
	LocalColorTableSize   uint8  `json:"local_color_table_size"`
	FileFinalPosition     int64  `json:"file_position"`
	fileStartPosition     int64
}

func (id GifImageDescriptor) String() string {
	var image_descriptor_string string = "Image Descriptor:\n"
	// return fmt.Sprintf("Image Descriptor:\nLeft Position: %d\nTop Position: %d\nWidth: %d\nHeight: %d\nHas Local Color Table: %t\nIs Interlaced: %t\nLocal Color Table Sorted: %t\nLocal Color Table Size: %d\n", id.ImageLeftPosition, id.ImageTopPosition, id.ImageWidth, id.ImageHeight, id.HasLocalColorTable, id.IsInterlaced, id.LocalColorTableSorted, id.LocalColorTableSize)

	image_descriptor_string += fmt.Sprintf("Left Position: %d\n", id.ImageLeftPosition)
	image_descriptor_string += fmt.Sprintf("Top Position: %d\n", id.ImageTopPosition)
	image_descriptor_string += fmt.Sprintf("Width: %d\n", id.ImageWidth)
	image_descriptor_string += fmt.Sprintf("Height: %d\n", id.ImageHeight)
	image_descriptor_string += fmt.Sprintf("Has Local Color Table: %t\n", id.HasLocalColorTable)
	image_descriptor_string += fmt.Sprintf("Is Interlaced: %t\n", id.IsInterlaced)
	image_descriptor_string += fmt.Sprintf("Local Color Table Sorted: %t\n", id.LocalColorTableSorted)
	image_descriptor_string += fmt.Sprintf("Local Color Table Size: %d\n", id.LocalColorTableSize)

	image_descriptor_string += fmt.Sprintf("\n\nImage Descriptor: %#x - %#x\n\n", id.fileStartPosition, id.FileFinalPosition)

	return image_descriptor_string
}

func (id GifImageDescriptor) GetSize() int {
	return 10
}

func (id GifImageDescriptor) GetLocalColorTableByteCount() int {
	if !id.HasLocalColorTable {
		return 0
	}

	var byte_count int = 3 * int(math.Pow(2, float64(id.LocalColorTableSize+1)))

	return byte_count
}

func (id *GifImageDescriptor) SetFileStartPosition(start_offset int64) {
	if id.fileStartPosition == 0 {
		id.fileStartPosition = start_offset
	}
}

// Parses a sequence of bytes into a GifImageDescriptor. It assumes that byte sequence only includes the image descriptor data. and
// Also assumes that the image separator byte is included in the byte sequence. Meaning it expects a chunk of 10 bytes.
func NewGifImageDescriptor(image_descriptor_fragment []byte) (*GifImageDescriptor, error) {
	if len(image_descriptor_fragment) != 10 {
		return nil, fmt.Errorf("Invalid image descriptor fragment")
	}

	var new_gif_image_descriptor *GifImageDescriptor = new(GifImageDescriptor)

	new_gif_image_descriptor.ImageLeftPosition = uint16(image_descriptor_fragment[1]) | (uint16(image_descriptor_fragment[2]) << 8)
	new_gif_image_descriptor.ImageTopPosition = uint16(image_descriptor_fragment[3]) | (uint16(image_descriptor_fragment[4]) << 8)

	new_gif_image_descriptor.ImageWidth = uint16(image_descriptor_fragment[5]) | (uint16(image_descriptor_fragment[6]) << 8)
	new_gif_image_descriptor.ImageHeight = uint16(image_descriptor_fragment[7]) | (uint16(image_descriptor_fragment[8]) << 8)

	var packed_fields byte = image_descriptor_fragment[9]

	new_gif_image_descriptor.HasLocalColorTable = (packed_fields & 0x80) != 0
	new_gif_image_descriptor.IsInterlaced = (packed_fields & 0x40) != 0
	new_gif_image_descriptor.LocalColorTableSorted = (packed_fields & 0x20) != 0
	new_gif_image_descriptor.LocalColorTableSize = packed_fields & 0x07

	return new_gif_image_descriptor, nil
}

/*
- ===========================================================
-                      GIF TABLE-BASED IMAGE DATA
- ===========================================================
From the GIF89a specification:
22. Table Based Image Data.

	a. Description. The image data for a table based image consists of a
	sequence of sub-blocks, of size at most 255 bytes each, containing an
	index into the active color table, for each pixel in the image.  Pixel
	indices are in order of left to right and from top to bottom.  Each index
	must be within the range of the size of the active color table, starting
	at 0. The sequence of indices is encoded using the LZW Algorithm with
	variable-length code, as described in Appendix F

	b. Required Version.  87a.

	c. Syntax. The image data format is as follows:

	 7 6 5 4 3 2 1 0        Field Name                    Type
	+---------------+
	|               |       LZW Minimum Code Size         Byte
	+---------------+

	+===============+
	|               |
	/               /       Image Data                    Data Sub-blocks
	|               |
	+===============+

	       i) LZW Minimum Code Size.  This byte determines the initial number
	       of bits used for LZW codes in the image data, as described in
	       Appendix F.

	d. Extensions and Scope. This block has no scope, it contains raster
	data. Extensions intended to modify a Table-based image must appear
	before the corresponding Image Descriptor.

	e. Recommendations. None.
*/
type GifTableBasedImageData struct {
	LZWMinimumCodeSize byte   `json:"lzw_minimum_code_size"`
	ImageData          []byte `json:"image_data"`
	FileFinalPosition  int64  `json:"file_position"`
	fileStartPosition  int64
}

func (tbid GifTableBasedImageData) String() string {
	var thrid_of_length int = len(tbid.ImageData) / 3
	var image_data_string string = "Table Based Image Data:\n"

	image_data_string += fmt.Sprintf("LZW Minimum Code Size: %d\n", tbid.LZWMinimumCodeSize)
	image_data_string += fmt.Sprintf("Image Data Size: %d bytes\n", len(tbid.ImageData))

	if thrid_of_length > 15 {
		image_data_string += fmt.Sprintf("Image data ends: %s...%s\n", tbid.getStrPrefix(7), tbid.getStrSuffix(7))
	} else {
		image_data_string += fmt.Sprintf("Image data: %s\n", tbid.getStrPrefix(thrid_of_length*4)) // full hex string
	}

	image_data_string += fmt.Sprintf("\n\nTable Based Image Data: %#x - %#x\n\n", tbid.fileStartPosition, tbid.FileFinalPosition)

	return image_data_string
}

func (tbid GifTableBasedImageData) GetPrefix(length int) []byte {
	if length > len(tbid.ImageData) {
		return tbid.ImageData
	}

	return tbid.ImageData[:length]
}

func (tbid GifTableBasedImageData) getStrPrefix(length int) string {
	var prefix []byte = tbid.GetPrefix(length)
	var string_prefix string = ""

	for h, byte_value := range prefix {
		string_prefix += fmt.Sprintf("%#x", byte_value)

		if h != (len(prefix) - 1) {
			string_prefix += " "
		}
	}

	return string_prefix
}

func (tbid GifTableBasedImageData) GetSuffix(length int) []byte {
	if length > len(tbid.ImageData) {
		return tbid.ImageData
	}

	return tbid.ImageData[len(tbid.ImageData)-length:]
}

func (tbid GifTableBasedImageData) getStrSuffix(length int) string {
	var suffix []byte = tbid.GetSuffix(length)
	var string_suffix string = ""

	for h, byte_value := range suffix {
		string_suffix += fmt.Sprintf("%#x", byte_value)

		if h != (len(suffix) - 1) {
			string_suffix += " "
		}
	}

	return string_suffix
}

func (tbid GifTableBasedImageData) GetSize() int {
	return int(tbid.FileFinalPosition - tbid.fileStartPosition)
}

func (tbid *GifTableBasedImageData) Decompress() ([]byte, error) {
	var compressed_data io.Reader = bytes.NewReader(tbid.ImageData)
	var err error
	lzw_reader := lzw.NewReader(compressed_data, lzw.LSB, int(tbid.LZWMinimumCodeSize))
	defer lzw_reader.Close()

	var decompressed_data []byte
	decompressed_data, err = io.ReadAll(lzw_reader)
	if err != nil {
		if decompressed_data != nil && len(decompressed_data) > 0 && err == io.ErrUnexpectedEOF {
			return decompressed_data, nil
		}
		return nil, fmt.Errorf("Error while decompressing image data: %s", err.Error())
	}

	return decompressed_data, nil
}

func NewGifTableBasedImageData(rs io.ReadSeeker) (*GifTableBasedImageData, error) {
	var current_seek_position int64
	var err error

	current_seek_position, err = rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current seek position: %s", err.Error())
	}

	var new_gif_table_based_image_data *GifTableBasedImageData = new(GifTableBasedImageData)
	new_gif_table_based_image_data.fileStartPosition = current_seek_position

	var lzw_minimum_code_size []byte = make([]byte, 1)
	_, err = rs.Read(lzw_minimum_code_size)
	if err != nil {
		return nil, fmt.Errorf("Error while reading LZW minimum code size: %s", err.Error())
	}

	new_gif_table_based_image_data.LZWMinimumCodeSize = lzw_minimum_code_size[0]

	image_data, err := ReadSubDataBlocks(rs)
	if err != nil {
		return nil, fmt.Errorf("Error while reading image data: %s", err.Error())
	}

	new_gif_table_based_image_data.ImageData = image_data

	current_seek_position, err = rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current seek position: %s", err.Error())
	}

	new_gif_table_based_image_data.FileFinalPosition = current_seek_position - 1

	return new_gif_table_based_image_data, nil
}
