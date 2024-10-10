package models

import (
	"fmt"
	"io"
	gif_helpers "libery-dungeon-libs/libs/gif_parsing/helpers"
)

const (
	APPLICATION_EXTENSION     string = "Application Extension"
	GRAPHIC_CONTROL_EXTENSION string = "Graphic Control Extension"
	COMMENT_EXTENSION         string = "Comment Extension"
	PLAIN_TEXT_EXTENSION      string = "Plain Text Extension"
)

type BlockScope int

const (
	NO_SCOPE BlockScope = iota
	GRAPHIC_RENDERING_BLOCK_SCOPE
	DATA_STREAM_SCOPE
)

type GifExtensionBlock interface {
	ExtensionName() string
	BlockScope() BlockScope
	String() string
}

// Reads the first byte which is assumed to be the byte size of the sub-blocks. Then reads the sub-blocks using the byte size.
// as delimiter.
func ReadSubDataBlocks(rs io.ReadSeeker) ([]byte, error) {
	var data_blocks []byte
	var data_block []byte
	var block_size uint8

	sub_block_size := make([]byte, 1)
	_, err := rs.Read(sub_block_size)
	if err != nil {
		return nil, fmt.Errorf("Error while reading sub-block size: %s", err.Error())
	}

	block_size = sub_block_size[0]

	if block_size == 0 {
		return nil, nil
	}

	data_blocks = make([]byte, 0)

	for block_size != 0 {
		data_block = make([]byte, block_size)

		_, err = rs.Read(data_block)
		if err != nil {
			return nil, fmt.Errorf("Error while reading sub-block: %s", err.Error())
		}

		data_blocks = append(data_blocks, data_block...)

		_, err = rs.Read(sub_block_size)
		if err != nil {
			return nil, fmt.Errorf("Error while reading sub-block size: %s", err.Error())
		}

		// Redundant
		if sub_block_size[0] == 0 {
			break
		}
		block_size = sub_block_size[0]
	}

	return data_blocks, err
}

/*
- ===========================================================
-                      GIF APPLICATION EXTENSION
- ===========================================================
From the GIF89a specification:
26. Application Extension.

	     a. Description. The Application Extension contains application-specific
	     information; it conforms with the extension block syntax, as described
	     below, and its block label is 0xFF.

	     b. Required Version.  89a.

	     c. Syntax.

	     7 6 5 4 3 2 1 0        Field Name                    Type
	    +---------------+
	 0  |               |       Extension Introducer          Byte
	    +---------------+
	 1  |               |       Extension Label               Byte
	    +---------------+

	    +---------------+
	 0  |               |       Block Size                    Byte
	    +---------------+
	 1  |               |
	    +-             -+
	 2  |               |
	    +-             -+
	 3  |               |       Application Identifier        8 Bytes
	    +-             -+
	 4  |               |
	    +-             -+
	 5  |               |
	    +-             -+
	 6  |               |
	    +-             -+
	 7  |               |
	    +-             -+
	 8  |               |
	    +---------------+
	 9  |               |
	    +-             -+
	10  |               |       Appl. Authentication Code     3 Bytes
	    +-             -+
	11  |               |
	    +---------------+

	    +===============+
	    |               |
	    |               |       Application Data              Data Sub-blocks
	    |               |
	    |               |
	    +===============+

	    +---------------+
	 0  |               |       Block Terminator              Byte
	    +---------------+

	           i) Extension Introducer - Defines this block as an extension. This
	           field contains the fixed value 0x21.

	           ii) Application Extension Label - Identifies the block as an
	           Application Extension. This field contains the fixed value 0xFF.

	           iii) Block Size - Number of bytes in this extension block,
	           following the Block Size field, up to but not including the
	           beginning of the Application Data. This field contains the fixed
	           value 11.

	           iv) Application Identifier - Sequence of eight printable ASCII
	           characters used to identify the application owning the Application
	           Extension.

	           v) Application Authentication Code - Sequence of three bytes used
	           to authenticate the Application Identifier. An Application program
	           may use an algorithm to compute a binary code that uniquely
	           identifies it as the application owning the Application Extension.


	     d. Extensions and Scope. This block does not have scope. This block
	     cannot be modified by any extension.

	     e. Recommendation. None.
*/
type GifApplicationExtensionBlock struct {
	extensionIntroducer uint8
	extensionLabel      uint8
	BlockSize           uint8  `json:"block_size"`
	ApplicationID       string `json:"application_id"`
	AuthenticationCode  string `json:"authentication_code"`
	ApplicationData     []byte `json:"application_data"`
	FileFinalPosition   int64  `json:"file_final_position"`
	fileStartPosition   int64
}

func (g GifApplicationExtensionBlock) ExtensionName() string {
	return APPLICATION_EXTENSION
}

func (g GifApplicationExtensionBlock) BlockScope() BlockScope {
	return NO_SCOPE
}

func (g GifApplicationExtensionBlock) String() string {
	var application_extension_string string = ""

	application_extension_string += fmt.Sprintf("Extension Introducer: %#x\n", g.extensionIntroducer)
	application_extension_string += fmt.Sprintf("Extension Label: %#x\n", g.extensionLabel)
	application_extension_string += fmt.Sprintf("Block Size: %d\n", g.BlockSize)
	application_extension_string += fmt.Sprintf("Application ID: %s\n", g.ApplicationID)
	application_extension_string += fmt.Sprintf("Authentication Code: %s\n", g.AuthenticationCode)

	application_extension_string += fmt.Sprintf("Application Data Size: %d\n", len(g.ApplicationData))

	application_extension_string += fmt.Sprintf("\n\nApplication Extension Block: %#x - %#x\n\n", g.fileStartPosition, g.FileFinalPosition)

	return application_extension_string
}

func (g GifApplicationExtensionBlock) GetSize() int {
	var extension_section_size int = 2      // Extension Introducer(1 Byte) + Extension Label(1 Byte)
	var application_extension_size int = 12 // Block Size(1 Byte) + Application Identifier(8 Bytes) + Application Authentication Code(3 Bytes)
	var application_data_size int = len(g.ApplicationData)

	// There is at least one block size byte(which could be the terminator) but if the size of application data is larger than 255
	// there will be multiple block size bytes. This is because every data sub block has at most 255 bytes preceded by a byte
	// indicating the size of the block.
	var application_data_block_sizes int = 1

	if application_data_size > 255 {
		application_data_block_sizes = 1 + (application_data_size / 255)
	}

	return extension_section_size + application_extension_size + application_data_size + application_data_block_sizes
}

// Create a new GifApplicationExtensionBlock from a reader-seeker. The offset should be at the beginning of the block.
// Denoted by the bytes 0x21 0xFF. The position of the reader-seeker will set to the next byte after the Block Terminator.
// If there is an error durning the creation, the reader-seeker will be at the same position as before the function call.
// Unless the error ocurred while trying to 'Rollback Seek' , in that case the position will be left as is.
func NewGifApplicationExtensionBlock(rs io.ReadSeeker) (*GifApplicationExtensionBlock, error) {
	offset_before_error, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current offset: %s", err.Error())
	}

	is_application_extension, err := gif_helpers.IsApplicationExtensionBlock(rs) // Preserves the offset
	if err != nil {
		return nil, err
	} else if !is_application_extension {
		return nil, fmt.Errorf("Try to read an application extension block but signature does not match")
	}
	var application_extension_block *GifApplicationExtensionBlock = new(GifApplicationExtensionBlock)
	application_extension_block.fileStartPosition = offset_before_error

	var application_extension_block_descriptor []byte = make([]byte, 14) // Extension Fields(2 bytes) + Application Extension Fields(12 bytes)

	_, err = rs.Read(application_extension_block_descriptor)
	if err != nil {
		_, rollback_err := rs.Seek(offset_before_error, io.SeekStart)
		if rollback_err != nil {
			return nil, fmt.Errorf("Error while reading application extension block: %s. Error while rolling back seek: %s", err.Error(), rollback_err.Error())
		}

		return nil, fmt.Errorf("Error while reading application extension block: %s", err.Error())
	}

	application_extension_block.extensionIntroducer = application_extension_block_descriptor[0]
	application_extension_block.extensionLabel = application_extension_block_descriptor[1]
	application_extension_block.BlockSize = application_extension_block_descriptor[2]

	// Value is fixed and must be 11
	if application_extension_block.BlockSize != 11 {
		_, rollback_err := rs.Seek(offset_before_error, io.SeekStart)
		if rollback_err != nil {
			return nil, fmt.Errorf("Error while reading application extension block: Block size does not match. Error while rolling back seek: %s", rollback_err.Error())
		}

		return nil, fmt.Errorf("Error while reading application extension block: Block size does not match")
	}

	application_extension_block.ApplicationID = string(application_extension_block_descriptor[3:11])
	application_extension_block.AuthenticationCode = string(application_extension_block_descriptor[11:14])

	data_blocks, err := ReadSubDataBlocks(rs) // Does not preserve the offset
	if err != nil {
		_, rollback_err := rs.Seek(offset_before_error, io.SeekStart)
		if rollback_err != nil {
			return nil, fmt.Errorf("Error while reading application extension block: %s. Error while rolling back seek: %s", err.Error(), rollback_err.Error())
		}

		return nil, fmt.Errorf("Error while reading application extension block: %s", err.Error())
	}

	application_extension_block.ApplicationData = data_blocks

	current_position, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current position: %s", err.Error())
	}

	application_extension_block.FileFinalPosition = current_position - 1

	return application_extension_block, nil
}

/*
- ===========================================================
-                      GIF GRAPHIC CONTROL EXTENSION
- ===========================================================
From the GIF89a specification:

 23. Graphic Control Extension.
    a. Description. The Graphic Control Extension contains parameters used
    when processing a graphic rendering block. The scope of this extension is
    the first graphic rendering block to follow. The extension contains only
    one data sub-block.

    This block is OPTIONAL; at most one Graphic Control Extension may precede
    a graphic rendering block. This is the only limit to the number of
    Graphic Control Extensions that may be contained in a Data Stream.

    b. Required Version.  89a.

    c. Syntax.

    7 6 5 4 3 2 1 0        Field Name                    Type
    +---------------+
    0  |               |       Extension Introducer          Byte
    +---------------+
    1  |               |       Graphic Control Label         Byte
    +---------------+

    +---------------+
    0  |               |       Block Size                    Byte
    +---------------+
    1  |     |     | | |       <Packed Fields>               See below
    +---------------+
    2  |               |       Delay Time                    Unsigned
    +-             -+
    3  |               |
    +---------------+
    4  |               |       Transparent Color Index       Byte
    +---------------+

    +---------------+
    0  |               |       Block Terminator              Byte
    +---------------+

    <Packed Fields>  =     Reserved                      3 Bits
    Disposal Method               3 Bits
    User Input Flag               1 Bit
    Transparent Color Flag        1 Bit

    i) Extension Introducer - Identifies the beginning of an extension
    block. This field contains the fixed value 0x21.

    ii) Graphic Control Label - Identifies the current block as a
    Graphic Control Extension. This field contains the fixed value
    0xF9.

    iii) Block Size - Number of bytes in the block, after the Block
    Size field and up to but not including the Block Terminator.  This
    field contains the fixed value 4.

    iv) Disposal Method - Indicates the way in which the graphic is to
    be treated after being displayed.

    Values :    0 -   No disposal specified. The decoder is
    not required to take any action.
    1 -   Do not dispose. The graphic is to be left
    in place.
    2 -   Restore to background color. The area used by the
    graphic must be restored to the background color.
    3 -   Restore to previous. The decoder is required to
    restore the area overwritten by the graphic with
    what was there prior to rendering the graphic.
    4-7 -    To be defined.

    v) User Input Flag - Indicates whether or not user input is
    expected before continuing. If the flag is set, processing will
    continue when user input is entered. The nature of the User input
    is determined by the application (Carriage Return, Mouse Button
    Click, etc.).

    Values :    0 -   User input is not expected.
    1 -   User input is expected.

    When a Delay Time is used and the User Input Flag is set,
    processing will continue when user input is received or when the
    delay time expires, whichever occurs first.

    vi) Transparency Flag - Indicates whether a transparency index is
    given in the Transparent Index field. (This field is the least
    significant bit of the byte.)

    Values :    0 -   Transparent Index is not given.
    1 -   Transparent Index is given.

    vii) Delay Time - If not 0, this field specifies the number of
    hundredths (1/100) of a second to wait before continuing with the
    processing of the Data Stream. The clock starts ticking immediately
    after the graphic is rendered. This field may be used in
    conjunction with the User Input Flag field.

    viii) Transparency Index - The Transparency Index is such that when
    encountered, the corresponding pixel of the display device is not
    modified and processing goes on to the next pixel. The index is
    present if and only if the Transparency Flag is set to 1.

    ix) Block Terminator - This zero-length data block marks the end of
    the Graphic Control Extension.

    d. Extensions and Scope. The scope of this Extension is the graphic
    rendering block that follows it; it is possible for other extensions to
    be present between this block and its target. This block can modify the
    Image Descriptor Block and the Plain Text Extension.

    e. Recommendations.

    i) Disposal Method - The mode Restore To Previous is intended to be
    used in small sections of the graphic; the use of this mode imposes
    severe demands on the decoder to store the section of the graphic
    that needs to be saved. For this reason, this mode should be used
    sparingly.  This mode is not intended to save an entire graphic or
    large areas of a graphic; when this is the case, the encoder should
    make every attempt to make the sections of the graphic to be
    restored be separate graphics in the data stream. In the case where
    a decoder is not capable of saving an area of a graphic marked as
    Restore To Previous, it is recommended that a decoder restore to
    the background color.

    ii) User Input Flag - When the flag is set, indicating that user
    input is expected, the decoder may sound the bell (0x07) to alert
    the user that input is being expected.  In the absence of a
    specified Delay Time, the decoder should wait for user input
    indefinitely.  It is recommended that the encoder not set the User
    Input Flag without a Delay Time specified.
*/
type GifGraphicControlExtensionBlock struct {
	extensionIntroducer uint8
	extensionLabel      uint8
	BlockSize           uint8                    `json:"block_size"`
	DisposalMethod      GifGraphicDisposalMethod `json:"disposal_method"`
	HaltForUserInput    bool                     `json:"halt_for_user_input"`
	HasTransparency     bool                     `json:"has_transparency"`
	DelayMs             int                      `json:"delay_ms"`
	TransparentIndex    uint8                    `json:"transparent_index"`
	FileFinalPosition   int64                    `json:"file_final_position"`
	fileStartPosition   int64
}

func (ggceb GifGraphicControlExtensionBlock) ExtensionName() string {
	return GRAPHIC_CONTROL_EXTENSION
}

func (ggceb GifGraphicControlExtensionBlock) BlockScope() BlockScope {
	return GRAPHIC_RENDERING_BLOCK_SCOPE
}

func (ggceb GifGraphicControlExtensionBlock) String() string {
	var graphic_control_extension_string string = "Graphic Control Extension\n"
	// return fmt.Sprintf("Extension Introducer: %d\nExtension Label: %d\nBlock Size: %d\nDisposal Method: %d\nHalt For User Input: %t\nHas Transparency: %t\nDelay: %d", g.extensionIntroducer, g.extensionLabel, g.BlockSize, g.DisposalMethod, g.HaltForUserInput, g.HasTransparency, g.DelayMs)

	graphic_control_extension_string += fmt.Sprintf("Extension Introducer: %#x\n", ggceb.extensionIntroducer)
	graphic_control_extension_string += fmt.Sprintf("Extension Label: %#x\n", ggceb.extensionLabel)
	graphic_control_extension_string += fmt.Sprintf("Block Size: %d\n", ggceb.BlockSize)
	graphic_control_extension_string += fmt.Sprintf("Disposal Method: %d\n", ggceb.DisposalMethod)
	graphic_control_extension_string += fmt.Sprintf("Halt For User Input: %t\n", ggceb.HaltForUserInput)
	graphic_control_extension_string += fmt.Sprintf("Has Transparency: %t\n", ggceb.HasTransparency)
	graphic_control_extension_string += fmt.Sprintf("Delay MS: %d\n", ggceb.DelayMs)
	graphic_control_extension_string += fmt.Sprintf("Transparent Index: %d\n", ggceb.TransparentIndex)

	graphic_control_extension_string += fmt.Sprintf("\n\nGraphic Control Extension Block: %#x - %#x\n\n", ggceb.fileStartPosition, ggceb.FileFinalPosition)

	return graphic_control_extension_string
}

type GifGraphicDisposalMethod uint8

const (
	NO_DISPOSAL_SPECIFIED       GifGraphicDisposalMethod = 0
	DO_NOT_DISPOSE              GifGraphicDisposalMethod = 1
	RESTORE_TO_BACKGROUND_COLOR GifGraphicDisposalMethod = 2
	RESTORE_TO_PREVIOUS         GifGraphicDisposalMethod = 3
)

// Create a new GifGraphicControlExtensionBlock from a reader-seeker. The offset should be at the beginning of the block
// Marked by the bytes 0x21 0xF9. The position of the reader-seeker will set to the next byte after the Block Terminator.
// If there is an error durning the creation, the reader-seeker will be at the same position as before the function call.
// Unless the error ocurred while trying to 'Rollback Seek' , in that case the position will be left as is.
func NewGifGraphicControlExtensionBlock(rs io.ReadSeeker) (*GifGraphicControlExtensionBlock, error) {
	offset_before_error, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current offset: %s", err.Error())
	}

	is_graphic_control_extension, err := gif_helpers.IsGraphicControlExtensionBlock(rs) // Preserves the offset
	if err != nil {
		return nil, err
	} else if !is_graphic_control_extension {
		return nil, fmt.Errorf("Try to read a graphic control extension block but signature does not match")
	}

	var graphic_control_extension_block *GifGraphicControlExtensionBlock = new(GifGraphicControlExtensionBlock)
	graphic_control_extension_block.fileStartPosition = offset_before_error

	var graphic_control_extension_block_descriptor []byte = make([]byte, 7) // Extension Fields(2 bytes) + Graphic Control Extension Fields(5 bytes)

	_, err = rs.Read(graphic_control_extension_block_descriptor)
	if err != nil {
		_, rollback_err := rs.Seek(offset_before_error, io.SeekStart)
		if rollback_err != nil {
			return nil, fmt.Errorf("Error while reading graphic control extension block: %s. Error while rolling back seek: %s", err.Error(), rollback_err.Error())
		}

		return nil, fmt.Errorf("Error while reading graphic control extension block: %s", err.Error())
	}

	graphic_control_extension_block.extensionIntroducer = graphic_control_extension_block_descriptor[0]
	graphic_control_extension_block.extensionLabel = graphic_control_extension_block_descriptor[1]
	graphic_control_extension_block.BlockSize = graphic_control_extension_block_descriptor[2]

	// Value is fixed and must be 4
	if graphic_control_extension_block.BlockSize != 4 {
		_, rollback_err := rs.Seek(offset_before_error, io.SeekStart)
		if rollback_err != nil {
			return nil, fmt.Errorf("Error while reading graphic control extension block: Block size does not match. Error while rolling back seek: %s", rollback_err.Error())
		}

		return nil, fmt.Errorf("Error while reading graphic control extension block: Block size does not match")
	}

	packed_fields := graphic_control_extension_block_descriptor[3]
	graphic_control_extension_block.DisposalMethod = GifGraphicDisposalMethod((packed_fields >> 2) & 0x07)
	graphic_control_extension_block.HaltForUserInput = (packed_fields & 0x02) != 0
	graphic_control_extension_block.HasTransparency = (packed_fields & 0x01) != 0

	var delay_in_hundredths_of_seconds uint16 = uint16(graphic_control_extension_block_descriptor[4]) | (uint16(graphic_control_extension_block_descriptor[5]) << 8)
	graphic_control_extension_block.DelayMs = int(delay_in_hundredths_of_seconds) * 10

	transparent_color_index := graphic_control_extension_block_descriptor[6]

	if graphic_control_extension_block.HasTransparency && transparent_color_index == 0 {
		fmt.Println("Warning: Transparent color index is 0 but transparency flag is set. This is against the GIF89a specification.")
	}

	graphic_control_extension_block.TransparentIndex = transparent_color_index

	seek_position, err := rs.Seek(1, io.SeekCurrent) // Skip the block terminator
	if err != nil {
		return nil, fmt.Errorf("Error while skipping block terminator: %s", err.Error())
	}

	graphic_control_extension_block.FileFinalPosition = seek_position - 1

	return graphic_control_extension_block, err
}

/*
- ===========================================================
-                      GIF COMMENT EXTENSION
- ===========================================================
From the GIF89a specification:
24. Comment Extension.

	    a. Description. The Comment Extension contains textual information which
	    is not part of the actual graphics in the GIF Data Stream. It is suitable
	    for including comments about the graphics, credits, descriptions or any
	    other type of non-control and non-graphic data.  The Comment Extension
	    may be ignored by the decoder, or it may be saved for later processing;
	    under no circumstances should a Comment Extension disrupt or interfere
	    with the processing of the Data Stream.

	    This block is OPTIONAL; any number of them may appear in the Data Stream.

	    b. Required Version.  89a.






















	                                                                      18


	    c. Syntax.

	    7 6 5 4 3 2 1 0        Field Name                    Type
	   +---------------+
	0  |               |       Extension Introducer          Byte
	   +---------------+
	1  |               |       Comment Label                 Byte
	   +---------------+

	   +===============+
	   |               |
	N  |               |       Comment Data                  Data Sub-blocks
	   |               |
	   +===============+

	   +---------------+
	0  |               |       Block Terminator              Byte
	   +---------------+

	          i) Extension Introducer - Identifies the beginning of an extension
	          block. This field contains the fixed value 0x21.

	          ii) Comment Label - Identifies the block as a Comment Extension.
	          This field contains the fixed value 0xFE.

	          iii) Comment Data - Sequence of sub-blocks, each of size at most
	          255 bytes and at least 1 byte, with the size in a byte preceding
	          the data.  The end of the sequence is marked by the Block
	          Terminator.

	          iv) Block Terminator - This zero-length data block marks the end of
	          the Comment Extension.

	    d. Extensions and Scope. This block does not have scope. This block
	    cannot be modified by any extension.

	    e. Recommendations.

	          i) Data - This block is intended for humans.  It should contain
	          text using the 7-bit ASCII character set. This block should
	          not be used to store control information for custom processing.

	          ii) Position - This block may appear at any point in the Data
	          Stream at which a block can begin; however, it is recommended that
	          Comment Extensions do not interfere with Control or Data blocks;
	          they should be located at the beginning or at the end of the Data
	          Stream to the extent possible.
*/
type GifCommentExtensionBlock struct {
	extensionIntroducer uint8
	extensionLabel      uint8
	CommentData         []byte `json:"comment_data"`
	FileFinalPosition   int64  `json:"file_final_position"`
	FileStartPosition   int64  `json:"file_start_position"`
}

func (g GifCommentExtensionBlock) String() string {
	var comment_extension_string string = "Comment Extension:\n"
	// return fmt.Sprintf("Extension Introducer: %d\nExtension Label: %d\nComment Data: %s", g.extensionIntroducer, g.extensionLabel, g.CommentData)
	comment_extension_string += fmt.Sprintf("Extension Introducer: %#x\n", g.extensionIntroducer)
	comment_extension_string += fmt.Sprintf("Extension Label: %#x\n", g.extensionLabel)
	comment_extension_string += fmt.Sprintf("Comment: %s\n", string(g.CommentData))

	comment_extension_string += fmt.Sprintf("\n\nComment Extension Block: %#x - %#x\n\n", g.FileStartPosition, g.FileFinalPosition)

	return comment_extension_string
}

func (g GifCommentExtensionBlock) ExtensionName() string {
	return COMMENT_EXTENSION
}

func (g GifCommentExtensionBlock) BlockScope() BlockScope {
	return NO_SCOPE
}

// Create a new GifCommentExtensionBlock from a reader-seeker. The offset should be at the beginning of the block
// marked by the bytes 0x21 0xFE. The position of the reader-seeker will set to the next byte after the Block Terminator.
// if there is an error durning the creation, the reader-seeker will be at the same position as before the function call.
// Unless the error ocurred while trying to 'Rollback Seek' , in that case the position will be left as is.
func NewGifCommentExtensionBlock(rs io.ReadSeeker) (*GifCommentExtensionBlock, error) {
	offset_before_error, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current offset: %s", err.Error())
	}

	is_comment_extension, err := gif_helpers.IsCommentExtensionBlock(rs) // Preserves the offset
	if err != nil {
		return nil, err
	} else if !is_comment_extension {
		return nil, fmt.Errorf("Try to read a comment extension block but signature does not match")
	}

	var comment_extension_block *GifCommentExtensionBlock = new(GifCommentExtensionBlock)
	comment_extension_block.FileStartPosition = offset_before_error
	comment_extension_block.extensionIntroducer = 0x21
	comment_extension_block.extensionLabel = 0xFE

	// We don't need to read the extension introducer or label. but we need to set the offset to the beginning data sub-blocks
	_, err = rs.Seek(2, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while setting offset to the beginning of the comment data: %s", err.Error())
	}

	comment_data, err := ReadSubDataBlocks(rs) // Does not preserve the offset
	if err != nil {
		_, rollback_err := rs.Seek(offset_before_error, io.SeekStart)
		if rollback_err != nil {
			return nil, fmt.Errorf("Error while reading comment extension block: %s. Error while rolling back seek: %s", err.Error(), rollback_err.Error())
		}

		return nil, fmt.Errorf("Error while reading comment extension block: %s", err.Error())
	}

	comment_extension_block.CommentData = comment_data

	current_position, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current position: %s", err.Error())
	}

	comment_extension_block.FileFinalPosition = current_position - 1

	// The block terminator is the same as the terminator of the data sub-blocks for the comment extension.

	return comment_extension_block, nil
}

/*
- ===========================================================
-                      GIF PLAIN TEXT EXTENSION
- ===========================================================
From the GIF89a specification:

 25. Plain Text Extension.
    a. Description. The Plain Text Extension contains textual data and the
    parameters necessary to render that data as a graphic, in a simple form.
    The textual data will be encoded with the 7-bit printable ASCII
    characters.  Text data are rendered using a grid of character cells

    19

    defined by the parameters in the block fields. Each character is rendered
    in an individual cell. The textual data in this block is to be rendered
    as mono-spaced characters, one character per cell, with a best fitting
    font and size. For further information, see the section on
    Recommendations below. The data characters are taken sequentially from
    the data portion of the block and rendered within a cell, starting with
    the upper left cell in the grid and proceeding from left to right and
    from top to bottom. Text data is rendered until the end of data is
    reached or the character grid is filled.  The Character Grid contains an
    integral number of cells; in the case that the cell dimensions do not
    allow for an integral number, fractional cells must be discarded; an
    encoder must be careful to specify the grid dimensions accurately so that
    this does not happen. This block requires a Global Color Table to be
    available; the colors used by this block reference the Global Color Table
    in the Stream if there is one, or the Global Color Table from a previous
    Stream, if one was saved. This block is a graphic rendering block,
    therefore it may be modified by a Graphic Control Extension.  This block
    is OPTIONAL; any number of them may appear in the Data Stream.

    b. Required Version.  89a.

    20

    c. Syntax.

    7 6 5 4 3 2 1 0        Field Name                    Type
    +---------------+
    0  |               |       Extension Introducer          Byte
    +---------------+
    1  |               |       Plain Text Label              Byte
    +---------------+

    +---------------+
    0  |               |       Block Size                    Byte
    +---------------+
    1  |               |       Text Grid Left Position       Unsigned
    +-             -+
    2  |               |
    +---------------+
    3  |               |       Text Grid Top Position        Unsigned
    +-             -+
    4  |               |
    +---------------+
    5  |               |       Text Grid Width               Unsigned
    +-             -+
    6  |               |
    +---------------+
    7  |               |       Text Grid Height              Unsigned
    +-             -+
    8  |               |
    +---------------+
    9  |               |       Character Cell Width          Byte
    +---------------+
    10  |               |       Character Cell Height         Byte
    +---------------+
    11  |               |       Text Foreground Color Index   Byte
    +---------------+
    12  |               |       Text Background Color Index   Byte
    +---------------+

    +===============+
    |               |
    N  |               |       Plain Text Data               Data Sub-blocks
    |               |
    +===============+

    +---------------+
    0  |               |       Block Terminator              Byte
    +---------------+

    i) Extension Introducer - Identifies the beginning of an extension
    block. This field contains the fixed value 0x21.

    ii) Plain Text Label - Identifies the current block as a Plain Text
    Extension. This field contains the fixed value 0x01.

    iii) Block Size - Number of bytes in the extension, after the Block
    Size field and up to but not including the beginning of the data
    portion. This field contains the fixed value 12.

    21

    iv) Text Grid Left Position - Column number, in pixels, of the left
    edge of the text grid, with respect to the left edge of the Logical
    Screen.

    v) Text Grid Top Position - Row number, in pixels, of the top edge
    of the text grid, with respect to the top edge of the Logical
    Screen.

    vi) Image Grid Width - Width of the text grid in pixels.

    vii) Image Grid Height - Height of the text grid in pixels.

    viii) Character Cell Width - Width, in pixels, of each cell in the
    grid.

    ix) Character Cell Height - Height, in pixels, of each cell in the
    grid.

    x) Text Foreground Color Index - Index into the Global Color Table
    to be used to render the text foreground.

    xi) Text Background Color Index - Index into the Global Color Table
    to be used to render the text background.

    xii) Plain Text Data - Sequence of sub-blocks, each of size at most
    255 bytes and at least 1 byte, with the size in a byte preceding
    the data.  The end of the sequence is marked by the Block
    Terminator.

    xiii) Block Terminator - This zero-length data block marks the end
    of the Plain Text Data Blocks.

    d. Extensions and Scope. The scope of this block is the Plain Text Data
    Block contained in it. This block may be modified by the Graphic Control
    Extension.

    e. Recommendations. The data in the Plain Text Extension is assumed to be
    preformatted. The selection of font and size is left to the discretion of
    the decoder.  If characters less than 0x20 or greater than 0xf7 are
    encountered, it is recommended that the decoder display a Space character
    (0x20). The encoder should use grid and cell dimensions such that an
    integral number of cells fit in the grid both horizontally as well as
    vertically.  For broadest compatibility, character cell dimensions should
    be around 8x8 or 8x16 (width x height); consider an image for unusual
    sized text.
*/
type GifPlainTextExtensionBlock struct {
	extensionIntroducer uint8
	extensionLabel      uint8
	BlockSize           uint8  `json:"block_size"`
	TextGridLeftPositon uint16 `json:"text_grid_left_position"`
	TextGridTopPosition uint16 `json:"text_grid_top_position"`
	TextGridWidth       uint16 `json:"text_grid_width"`
	TextGridHeight      uint16 `json:"text_grid_height"`
	CharacterCellWidth  uint8  `json:"character_cell_width"`
	CharacterCellHeight uint8  `json:"character_cell_height"`
	TextForegroundColor uint8  `json:"text_foreground_color"`
	TextBackgroundColor uint8  `json:"text_background_color"`
	PlainTextData       []byte `json:"plain_text_data"`
}

func (g GifPlainTextExtensionBlock) ExtensionName() string {
	return PLAIN_TEXT_EXTENSION
}

func (g GifPlainTextExtensionBlock) BlockScope() BlockScope {
	return GRAPHIC_RENDERING_BLOCK_SCOPE
}

func (g GifPlainTextExtensionBlock) String() string {
	return fmt.Sprintf("Extension Introducer: %d\nExtension Label: %d\nBlock Size: %d\nText Grid Left Position: %d\nText Grid Top Position: %d\nText Grid Width: %d\nText Grid Height: %d\nCharacter Cell Width: %d\nCharacter Cell Height: %d\nText Foreground Color: %d\nText Background Color: %d\nPlain Text Data: %s", g.extensionIntroducer, g.extensionLabel, g.BlockSize, g.TextGridLeftPositon, g.TextGridTopPosition, g.TextGridWidth, g.TextGridHeight, g.CharacterCellWidth, g.CharacterCellHeight, g.TextForegroundColor, g.TextBackgroundColor, g.PlainTextData)
}

// Create a new GifPlainTextExtensionBlock from a reader-seeker. The offset should be at the beginning of the block
// marked by the bytes 0x21 0x01. The position of the reader-seeker will set to the next byte after the Block Terminator.
// if there is an error durning the creation, the reader-seeker will be at the same position as before the function call.
// Unless the error ocurred while trying to 'Rollback Seek' , in that case the position will be left as is.
func NewGifPlainTextExtensionBlock(rs io.ReadSeeker) (*GifPlainTextExtensionBlock, error) {
	offset_before_error, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current offset: %s", err.Error())
	}

	is_plain_text_extension, err := gif_helpers.IsPlainTextExtensionBlock(rs) // Preserves the offset
	if err != nil {
		return nil, err
	} else if !is_plain_text_extension {
		return nil, fmt.Errorf("Try to read a plain text extension block but signature does not match")
	}

	var plain_text_extension_block *GifPlainTextExtensionBlock = new(GifPlainTextExtensionBlock)
	var plain_text_extension_block_descriptor []byte = make([]byte, 15) // Extension Fields(2 bytes) + Plain Text Extension Fields(13 bytes)

	_, err = rs.Read(plain_text_extension_block_descriptor)
	if err != nil {
		_, rollback_err := rs.Seek(offset_before_error, io.SeekStart)
		if rollback_err != nil {
			return nil, fmt.Errorf("Error while reading plain text extension block: %s. Error while rolling back seek: %s", err.Error(), rollback_err.Error())
		}

		return nil, fmt.Errorf("Error while reading plain text extension block: %s", err.Error())
	}

	plain_text_extension_block.extensionIntroducer = plain_text_extension_block_descriptor[0]
	plain_text_extension_block.extensionLabel = plain_text_extension_block_descriptor[1]
	plain_text_extension_block.BlockSize = plain_text_extension_block_descriptor[2]
	plain_text_extension_block.TextGridLeftPositon = uint16(plain_text_extension_block_descriptor[3]) | (uint16(plain_text_extension_block_descriptor[4]) << 8)
	plain_text_extension_block.TextGridTopPosition = uint16(plain_text_extension_block_descriptor[5]) | (uint16(plain_text_extension_block_descriptor[6]) << 8)
	plain_text_extension_block.TextGridWidth = uint16(plain_text_extension_block_descriptor[7]) | (uint16(plain_text_extension_block_descriptor[8]) << 8)
	plain_text_extension_block.TextGridHeight = uint16(plain_text_extension_block_descriptor[9]) | (uint16(plain_text_extension_block_descriptor[10]) << 8)
	plain_text_extension_block.CharacterCellWidth = plain_text_extension_block_descriptor[11]
	plain_text_extension_block.CharacterCellHeight = plain_text_extension_block_descriptor[12]
	plain_text_extension_block.TextForegroundColor = plain_text_extension_block_descriptor[13]
	plain_text_extension_block.TextBackgroundColor = plain_text_extension_block_descriptor[14]

	var plain_text_data []byte

	plain_text_data, err = ReadSubDataBlocks(rs) // Does not preserve the offset
	if err != nil {
		_, rollback_err := rs.Seek(offset_before_error, io.SeekStart)
		if rollback_err != nil {
			return nil, fmt.Errorf("Error while reading plain text extension block: %s. Error while rolling back seek: %s", err.Error(), rollback_err.Error())
		}

		return nil, fmt.Errorf("Error while reading plain text extension block: %s", err.Error())
	}

	plain_text_extension_block.PlainTextData = plain_text_data

	rs.Seek(1, io.SeekCurrent) // Skip the block terminator

	return plain_text_extension_block, nil
}
