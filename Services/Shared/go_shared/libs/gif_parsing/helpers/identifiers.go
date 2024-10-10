package helpers

import (
	"io"
)

// Check if the block is a control block. It expects the reader to be at position where a control block is expected to be read.
// It will maintain the reader's position.
func IsControlBlock(rs io.ReadSeeker) (bool, error) {
	block_label, err := ReadPreservingOffset(rs, 1)
	if err != nil {
		return false, err
	}

	return block_label[0] >= 0x80 && block_label[0] <= 0xF9, nil
}

// Check if the block is a special purpose block. It expects the reader to be at position where a special purpose
// block is expected to be read. It will maintain the reader's position.
func IsSpecialPurposeBlock(rs io.ReadSeeker) (bool, error) {
	block_label, err := ReadPreservingOffset(rs, 1)
	if err != nil {
		return false, err
	}

	return block_label[0] >= 0xFA && block_label[0] <= 0xFF, nil
}

// Check if the block is a graphic rendering block. It expects the reader to be at position where a graphic rendering
// block is expected to be read. It will maintain the reader's position.
func IsGraphicRenderingBlock(rs io.ReadSeeker) (bool, error) {
	block_label, err := ReadPreservingOffset(rs, 1)
	if err != nil {
		return false, err
	}

	return block_label[0] >= 0x00 && block_label[0] <= 0x7F && block_label[0] != 0x3B, nil
}

// Frame blocks are not a thing in the GIF specification. But in this program context, they refer to a set of blocks that hold
// and image or are related only to it. meaning Graphic Control Extension, Image Descriptor. Also Local Color Table and Image Data.
// but these last two are not labeled and instead are expected to appear right after the Image Descriptor(Local Color Table, Image Data).
func IsFrameBlock(rs io.ReadSeeker) (bool, error) {
	var is_frame_block bool = false

	is_frame_block, err := IsGraphicControlExtensionBlock(rs)
	if err != nil {
		return false, err
	}

	if !is_frame_block {
		is_frame_block, err = IsImageDescriptorBlock(rs)
		if err != nil {
			return false, err
		}
	}

	return is_frame_block, nil
}

// Check if the block is an image descriptor block. It expects the reader to be at position where an image descriptor
// block is expected to be read. It will maintain the reader's position.
func IsImageDescriptorBlock(rs io.ReadSeeker) (bool, error) {
	block_label, err := ReadPreservingOffset(rs, 1)
	if err != nil {
		return false, err
	}

	return block_label[0] == 0x2C, nil
}

func IsExtensionBlock(r io.ReadSeeker) (bool, error) {
	block_label, err := ReadPreservingOffset(r, 1)
	if err != nil {
		return false, err
	}

	return block_label[0] == 0x21, nil
}

// Check if the block is an application extension block. It expects the reader to be at position where an application extension
// block is expected to be read. It will maintain the reader's position.
func IsApplicationExtensionBlock(r io.ReadSeeker) (bool, error) {
	extension_block_label, err := ReadPreservingOffset(r, 2)
	if err != nil {
		return false, err
	}

	// fmt.Printf("Extension block label: %#x %#x\n", extension_block_label[0], extension_block_label[1])

	return extension_block_label[0] == 0x21 && extension_block_label[1] == 0xFF, nil
}

// Check if the block is a graphic control extension block. It expects the reader to be at position where a graphic control
// extension block is expected to be read. It will maintain the reader's position.
func IsGraphicControlExtensionBlock(r io.ReadSeeker) (bool, error) {
	extension_block_label, err := ReadPreservingOffset(r, 2)
	if err != nil {
		return false, err
	}

	return extension_block_label[0] == 0x21 && extension_block_label[1] == 0xF9, nil
}

// Check if the block is a comment extension block. It expects the reader to be at position where a comment extension
// block is expected to be read. It will maintain the reader's position.
func IsCommentExtensionBlock(r io.ReadSeeker) (bool, error) {
	extension_block_label, err := ReadPreservingOffset(r, 2)
	if err != nil {
		return false, err
	}

	return extension_block_label[0] == 0x21 && extension_block_label[1] == 0xFE, nil
}

// Check if the block is a plain text extension block. It expects the reader to be at position where a plain text extension
// block is expected to be read. It will maintain the reader's position.
func IsPlainTextExtensionBlock(r io.ReadSeeker) (bool, error) {
	extension_block_label, err := ReadPreservingOffset(r, 2)
	if err != nil {
		return false, err
	}

	return extension_block_label[0] == 0x21 && extension_block_label[1] == 0x01, nil
}

// Trailer is a single field block indicating the end of the GIF data stream. Check if the block is a trailer block.
// It expects the reader to be at position where a trailer block is expected to be read. It will maintain the reader's position.
func IsTrailer(rs io.ReadSeeker) (bool, error) {
	block_label, err := ReadPreservingOffset(rs, 1)
	if err != nil {
		return false, err
	}

	return block_label[0] == 0x3B, nil
}
