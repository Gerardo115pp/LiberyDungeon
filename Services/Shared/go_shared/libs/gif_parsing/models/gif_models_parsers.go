package models

import (
	"fmt"
	"io"
	gif_helpers "libery-dungeon-libs/libs/gif_parsing/helpers"
)

type GifBlockDetector struct {
	IsImageDescriptorBlock bool
	IsFrameBlock           bool
	IsGraphicControlBlock  bool
	IsApplicationBlock     bool
	IsCommentBlock         bool
}

func FillBlockDetector(data_stream io.ReadSeeker) (*GifBlockDetector, error) {
	var block_detector *GifBlockDetector = new(GifBlockDetector)
	var err error

	block_detector.IsImageDescriptorBlock, err = gif_helpers.IsImageDescriptorBlock(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while checking if block is an image descriptor block: %s", err.Error())
	}
	block_detector.IsFrameBlock, err = gif_helpers.IsFrameBlock(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while checking if block is a graphic control extension block: %s", err.Error())
	}
	block_detector.IsGraphicControlBlock, err = gif_helpers.IsGraphicControlExtensionBlock(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while checking if block is a graphic control extension block: %s", err.Error())
	}
	block_detector.IsApplicationBlock, err = gif_helpers.IsApplicationExtensionBlock(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while checking if block is an application extension block: %s", err.Error())
	}
	block_detector.IsCommentBlock, err = gif_helpers.IsCommentExtensionBlock(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while checking if block is a comment block: %s", err.Error())
	}

	return block_detector, nil
}

func GetGifSignature(data_stream io.Reader) (*GifHeader, error) {
	var gif_file_signature []byte = make([]byte, 6)

	_, err := data_stream.Read(gif_file_signature)
	if err != nil {
		return nil, fmt.Errorf("Error while reading gif file signature: %s", err.Error())
	}

	gif_header, err := NewGifHeader(gif_file_signature)
	if err != nil {
		return nil, fmt.Errorf("Error while creating gif header: %s", err.Error())
	}

	gif_header.FileFinalPosition = 0x06

	return gif_header, nil
}

func GetGifLogicalScreenDescriptor(data_stream io.ReadSeeker) (*GifLogicalScreenDescriptor, error) {
	var gif_logical_screen_descriptor []byte = make([]byte, 7)

	_, err := data_stream.Read(gif_logical_screen_descriptor)
	if err != nil {
		return nil, fmt.Errorf("Error while reading gif logical screen descriptor: %s", err.Error())
	}

	gif_lsd, err := NewGifLogicalScreenDescriptor(gif_logical_screen_descriptor)
	if err != nil {
		return nil, fmt.Errorf("Error while creating gif logical screen descriptor: %s", err.Error())
	}

	current_position, err := data_stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current position: %s", err.Error())
	}

	gif_lsd.FileFinalPosition = current_position

	return gif_lsd, nil
}

func ParseGlobalColorTable(data_stream io.ReadSeeker, gif_data *ParsedGif) (*GifColorTable, error) {
	var global_color_table *GifColorTable
	var global_color_table_byte_count int = gif_data.LogicalScreenDescriptor.GetGlobalColorTableByteCount()

	global_color_table, err := ParseColorTable(data_stream, global_color_table_byte_count)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing global color table: %s", err.Error())
	}

	return global_color_table, nil
}

func ParseColorTable(data_stream io.ReadSeeker, color_table_byte_count int) (*GifColorTable, error) {
	var current_seek_position int64

	current_seek_position, err := data_stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current seek position: %s", err.Error())
	}

	var color_table_chunk []byte = make([]byte, color_table_byte_count)
	_, err = data_stream.Read(color_table_chunk)
	if err != nil {
		return nil, fmt.Errorf("Error while reading color table: %s", err.Error())
	}

	color_table, err := NewGifColorTable(color_table_chunk)
	if err != nil {
		return nil, fmt.Errorf("Error while creating color table: %s", err.Error())
	}

	color_table.FileStartPosition = current_seek_position

	current_seek_position, err = data_stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current seek position: %s", err.Error())
	}

	color_table.FileFinalPosition = current_seek_position

	return color_table, nil
}

func ParseExtensionBlock(data_stream io.ReadSeeker) (GifExtensionBlock, error) {
	is_extension_block, err := gif_helpers.IsExtensionBlock(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while checking if block is an extension block: %s", err.Error())
	}
	if !is_extension_block {
		return nil, fmt.Errorf("Block is not an extension block")
	}

	var extension_block GifExtensionBlock

	if is_application_extension_block, err := gif_helpers.IsApplicationExtensionBlock(data_stream); err == nil && is_application_extension_block {
		var application_extension_block *GifApplicationExtensionBlock
		application_extension_block, err = NewGifApplicationExtensionBlock(data_stream)
		if err != nil {
			fmt.Printf("Error while creating application extension block: ")
			gif_helpers.PrintFileOffset(data_stream)
			return nil, fmt.Errorf("Error while creating application extension block: %s", err.Error())
		}

		extension_block = *application_extension_block
	} else if is_graphic_control_extension_block, err := gif_helpers.IsGraphicControlExtensionBlock(data_stream); err == nil && is_graphic_control_extension_block {
		var graphic_control_extension_block *GifGraphicControlExtensionBlock
		graphic_control_extension_block, err = NewGifGraphicControlExtensionBlock(data_stream)
		if err != nil {
			fmt.Printf("Error while creating graphic control extension block: ")
			gif_helpers.PrintFileOffset(data_stream)
			return nil, fmt.Errorf("Error while creating graphic control extension block: %s", err.Error())
		}

		extension_block = *graphic_control_extension_block
	} else if is_comment_block, err := gif_helpers.IsCommentExtensionBlock(data_stream); err == nil && is_comment_block {
		var comment_block *GifCommentExtensionBlock
		comment_block, err = NewGifCommentExtensionBlock(data_stream)
		if err != nil {
			fmt.Printf("Error while creating comment extension block: ")
			gif_helpers.PrintFileOffset(data_stream)
			return nil, fmt.Errorf("Error while creating comment extension block: %s", err.Error())
		}

		extension_block = *comment_block
	} else {
		return nil, fmt.Errorf("Block is not an application extension block")
	}

	return extension_block, nil
}

func ParseImageDescriptorBlock(data_stream io.ReadSeeker) (*GifImageDescriptor, error) {
	if next_is_image_descriptor_block, err := gif_helpers.IsImageDescriptorBlock(data_stream); err != nil || !next_is_image_descriptor_block {
		return nil, fmt.Errorf("Block has no image descriptor data")
	}

	var current_seek_position int64

	current_seek_position, err := data_stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current seek position: %s", err.Error())
	}

	var image_descriptor_chunk []byte = make([]byte, 10)
	_, err = data_stream.Read(image_descriptor_chunk)
	if err != nil {
		return nil, fmt.Errorf("Error while reading image descriptor: %s", err.Error())
	}

	image_descriptor, err := NewGifImageDescriptor(image_descriptor_chunk)
	if err != nil {
		return nil, fmt.Errorf("Error while creating image descriptor: %s", err.Error())
	}

	image_descriptor.SetFileStartPosition(current_seek_position)
	image_descriptor.FileFinalPosition = current_seek_position + 9

	return image_descriptor, nil
}

func ParseImageGraphicRenderingBlock(data_stream io.ReadSeeker) (*GifGraphicRenderingBlock, error) {
	var current_seek_position int64
	var err error

	current_seek_position, err = data_stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current seek position: %s", err.Error())
	}

	var new_graphic_rendering_block *GifGraphicRenderingBlock = new(GifGraphicRenderingBlock)
	new_graphic_rendering_block.SetFileStartPosition(current_seek_position)

	image_descriptor, err := ParseImageDescriptorBlock(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing image descriptor: %s", err.Error())
	}

	new_graphic_rendering_block.ImageDescriptor = image_descriptor

	if new_graphic_rendering_block.ImageDescriptor.HasLocalColorTable {
		var local_color_table *GifColorTable
		local_color_table, err = ParseColorTable(data_stream, new_graphic_rendering_block.ImageDescriptor.GetLocalColorTableByteCount())

		new_graphic_rendering_block.LocalColorTable = local_color_table
	}

	var table_based_image_data *GifTableBasedImageData

	table_based_image_data, err = NewGifTableBasedImageData(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while creating table based image data: %s", err.Error())
	}

	new_graphic_rendering_block.ImageData = table_based_image_data

	current_seek_position, err = data_stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current seek position: %s", err.Error())
	}

	new_graphic_rendering_block.FileFinalPosition = current_seek_position - 1

	return new_graphic_rendering_block, nil
}

/**
 * Parses The Header, Logical Screen Descriptor, and Global Color Table if one is present. A Reader at position 0 is expected.
 */
func ParseGifGlobalData(data_stream io.ReadSeeker) (*ParsedGif, error) {
	var gif_data *ParsedGif = NewParsedGif()
	var gif_file_signature *GifHeader
	var gif_logical_screen_descriptor *GifLogicalScreenDescriptor

	gif_file_signature, err := GetGifSignature(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while getting gif signature: %s", err.Error())
	}

	gif_logical_screen_descriptor, err = GetGifLogicalScreenDescriptor(data_stream)
	if err != nil {
		return nil, fmt.Errorf("Error while getting gif logical screen descriptor: %s", err.Error())
	}

	gif_data.Header = *gif_file_signature
	gif_data.LogicalScreenDescriptor = *gif_logical_screen_descriptor

	if gif_data.LogicalScreenDescriptor.HasGlobalColorTable {
		global_color_table, err := ParseGlobalColorTable(data_stream, gif_data)
		if err != nil {
			return nil, fmt.Errorf("Error while parsing global color table: %s", err.Error())
		}

		gif_data.GlobalColorTable = global_color_table
	}

	return gif_data, nil
}

func ParseGifBlocks(data_stream io.ReadSeeker, gif_data *ParsedGif) error {
	var found_trailer bool
	var block_detector *GifBlockDetector
	var graphic_control_extension_block *GifGraphicControlExtensionBlock // This will be added to the next graphic rendering block found

	// Debug Variables
	var round_number int = 0

	var err error

	found_trailer, err = gif_helpers.IsTrailer(data_stream)
	if err != nil {
		return fmt.Errorf("Error while checking if block is a trailer: %s", err.Error())
	}

	for !found_trailer {
		block_detector, err = FillBlockDetector(data_stream)
		if err != nil {
			return fmt.Errorf("Error while filling block detector: %s", err.Error())
		}

		switch true {
		case block_detector.IsGraphicControlBlock:
			graphic_control_extension_block, err = NewGifGraphicControlExtensionBlock(data_stream)
			if err != nil {
				return fmt.Errorf("Error while creating graphic control extension block: %s", err.Error())
			}
		case block_detector.IsImageDescriptorBlock:
			graphic_rendering_block, err := ParseImageGraphicRenderingBlock(data_stream)
			if err != nil {
				return fmt.Errorf("Error while parsing graphic rendering block: %s", err.Error())
			}

			// If we found a graphic control extension block before this image descriptor block, add it to the graphic rendering block
			if graphic_control_extension_block != nil {
				graphic_rendering_block.GraphicControlExtension = graphic_control_extension_block
				graphic_control_extension_block = nil
			}

			gif_data.GraphicRenderingBlocks = append(gif_data.GraphicRenderingBlocks, *graphic_rendering_block)
		case block_detector.IsApplicationBlock:
			extension_block, err := ParseExtensionBlock(data_stream)
			if err != nil {
				return fmt.Errorf("Error while parsing extension block: %s", err.Error())
			}

			gif_data.NoScopeExtensions = append(gif_data.NoScopeExtensions, extension_block)
		case block_detector.IsCommentBlock:
			extension_block, err := ParseExtensionBlock(data_stream)
			if err != nil {
				return fmt.Errorf("Error while parsing extension block: %s", err.Error())
			}

			gif_data.NoScopeExtensions = append(gif_data.NoScopeExtensions, extension_block)
		default:
			fmt.Printf("Halted at round %d\n", round_number)
			gif_helpers.PrintFileOffset(data_stream)

			if round_number != 0 {
				fmt.Println(gif_data)
			}

			block_label, err := gif_helpers.ReadPreservingOffset(data_stream, 2)
			if err != nil {
				return fmt.Errorf("Error while reading block label: %s", err.Error())
			}

			return fmt.Errorf("Untrusted block found. Halting so i can test it.\nblock label: %#x %#x", block_label[0], block_label[1])
		}

		round_number++

		found_trailer, err = gif_helpers.IsTrailer(data_stream)
		if err != nil {
			return fmt.Errorf("Error while checking if block is a trailer: %s", err.Error())
		}
	}

	return nil
}
