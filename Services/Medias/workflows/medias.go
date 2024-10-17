package workflows

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	dungeon_helpers "libery-dungeon-libs/helpers"
	gif_parsing_models "libery-dungeon-libs/libs/gif_parsing/models"
	gif_parsing_workflows "libery-dungeon-libs/libs/gif_parsing/workflows"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_medias_service/Config"
	service_helpers "libery_medias_service/helpers"
	service_models "libery_medias_service/models"
	"libery_medias_service/repository"
	workflow_errors "libery_medias_service/workflows/errors"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"

	"github.com/Gerardo115pp/thumbnailer"
	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
)

// creates a thumbnail for an image file of mime type jpeg, png or webp. If keep_format is false, the thumbnail will be encoded as jpeg.
// *UNDER NO CIRCUMSTANCES* this function will upscale an image. Meaning if the original image width is 1000px and the thumbnail_width is 2000px, the resulting thumbnail will be 1000px.
func GetImageThumbnail(f *os.File, mime_type string, thumbnail_width int, keep_format bool) (*service_models.ThumbnailResponse, *dungeon_models.LabeledError) {
	var new_thumbnail_response *service_models.ThumbnailResponse = new(service_models.ThumbnailResponse)

	new_thumbnail_response.MimeType = mime_type
	new_thumbnail_response.Filename = filepath.Base(f.Name())

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error decoding image", dungeon_models.ErrProcessError)
	}

	image_rect := img.Bounds()
	image_width := image_rect.Dx()
	image_height := image_rect.Dy()
	height_ratio := float64(image_height) / float64(image_width)

	if thumbnail_width == 0 {
		thumbnail_width = app_config.THUMBNAIL_WIDTH
	}

	if thumbnail_width >= image_width {
		thumbnail_width = image_width
	}

	thumbnail_height := int(float64(thumbnail_width) * height_ratio)

	thumbnail := image.NewRGBA(image.Rect(0, 0, thumbnail_width, thumbnail_height))

	draw.ApproxBiLinear.Scale(thumbnail, thumbnail.Bounds(), img, img.Bounds(), draw.Over, nil)

	var buf bytes.Buffer

	if mime_type == "image/jpeg" || !keep_format {
		err = jpeg.Encode(&buf, thumbnail, nil)
	} else if mime_type == "image/png" {
		err = png.Encode(&buf, thumbnail)
	} else if mime_type == "image/webp" {
		err = webp.Encode(&buf, thumbnail, &webp.Options{Lossless: true, Quality: 75})
	} else {
		return nil, dungeon_models.NewLabeledError(fmt.Errorf("Unsupported format: %s", mime_type), "While encoding thumbnail", dungeon_models.ErrProcessError)
	}

	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error encoding thumbnail", dungeon_models.ErrProcessError)
	}

	new_thumbnail_response.MediaStream = &buf
	new_thumbnail_response.MediaLength = int64(buf.Len())
	new_thumbnail_response.Resized = true
	new_thumbnail_response.Size = &service_models.MediaSize{
		Width:  thumbnail_width,
		Height: thumbnail_height,
	}
	new_thumbnail_response.OrignialSize = &service_models.MediaSize{
		Width:  image_width,
		Height: image_height,
	}

	return new_thumbnail_response, nil
}

// creates a thumbnail for a gif file. The thumbnail will be the first frame of the gif, resized to the specified width. If the width is 0, the default thumbnail width will be used.
// *UNDER NO CIRCUMSTANCES* this function will upscale an image. Meaning if the original image width is 1000px and the thumbnail_width is 2000px, the resulting thumbnail will be 1000px.
func GetGifThumbnail(f *os.File, thumbnail_width int) (*service_models.ThumbnailResponse, *dungeon_models.LabeledError) {
	if is_gif, err := service_helpers.IsGIF(f); err != nil || !is_gif {
		return nil, dungeon_models.NewLabeledError(err, "File is not a gif", workflow_errors.ErrMediaNotSupported)
	}

	var gif_file *gif_parsing_models.ParsedGif

	gif_file, err := gif_parsing_workflows.ReadGifFile(f)
	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error reading gif file", dungeon_models.ErrProcessError)
	}

	if thumbnail_width == 0 {
		thumbnail_width = app_config.THUMBNAIL_WIDTH
	}

	var available_frames int = gif_file.GetFrameCount()
	if available_frames == 0 {
		return nil, dungeon_models.NewLabeledError(fmt.Errorf("No frames found in gif"), "While getting gif frames", dungeon_models.ErrProcessError)
	}

	var used_graphic_rendering_block *gif_parsing_models.GifGraphicRenderingBlock = &gif_file.GraphicRenderingBlocks[0]

	gif_frame, err := gif_file.RenderResizeFrame(0, thumbnail_width)
	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error rendering gif frame", dungeon_models.ErrProcessError)
	}

	var thumbnail_buffer *bytes.Buffer = new(bytes.Buffer)

	err = jpeg.Encode(thumbnail_buffer, gif_frame, nil)
	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error encoding gif thumbnail as jpeg", dungeon_models.ErrProcessError)
	}

	var new_thumbnail_response *service_models.ThumbnailResponse = new(service_models.ThumbnailResponse)

	new_thumbnail_response.MimeType = "image/jpeg"
	new_thumbnail_response.MediaStream = thumbnail_buffer
	new_thumbnail_response.MediaLength = int64(thumbnail_buffer.Len())
	new_thumbnail_response.Filename = filepath.Base(gif_file.Filename)
	new_thumbnail_response.Resized = true
	new_thumbnail_response.Size = &service_models.MediaSize{
		Width:  gif_frame.Bounds().Dx(),
		Height: gif_frame.Bounds().Dy(),
	}
	new_thumbnail_response.OrignialSize = &service_models.MediaSize{
		Width:  int(used_graphic_rendering_block.ImageDescriptor.ImageWidth),
		Height: int(used_graphic_rendering_block.ImageDescriptor.ImageHeight),
	}

	return new_thumbnail_response, nil
}

func GetFileThumbnail(f *os.File, thumbnail_width int) (*service_models.ThumbnailResponse, *dungeon_models.LabeledError) {
	var thumbnail_response *service_models.ThumbnailResponse
	var labeled_err *dungeon_models.LabeledError

	mime_type, err := service_helpers.GetMimeType(f)
	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error getting mime type", dungeon_models.ErrProcessError)
	}

	switch mime_type {
	case "image/jpeg", "image/png", "image/webp":
		thumbnail_response, labeled_err = GetImageThumbnail(f, mime_type, thumbnail_width, false)
	case "video/mp4", "video/webm":
		thumbnail_response, labeled_err = GetVideoThumbnail(f, thumbnail_width)
	case "image/gif":
		thumbnail_response, labeled_err = GetGifThumbnail(f, thumbnail_width)
	default:
		labeled_err = dungeon_models.NewLabeledError(fmt.Errorf("Unsupported mime type: %s", mime_type), "While getting file thumbnail", workflow_errors.ErrMediaNotSupported)
	}

	return thumbnail_response, labeled_err
}

// Ensures that a media identity is unique and creates a media record on the database.
// If the media identity name needs to be changed, the change will be made in place(meaning i will changed the object passed as an argument)
// Performs no file operations.
func InsertUniqueMediaIdentity(media_identity *dungeon_models.MediaIdentity) error {
	var err error
	category_path := filepath.Join(media_identity.ClusterPath, media_identity.CategoryPath)

	would_be_path := filepath.Join(category_path, media_identity.Media.Name)

	if dungeon_helpers.FileExists(would_be_path) {
		err = SetUniqueMediaName(media_identity)
		if err != nil {
			return fmt.Errorf("Error setting unique media name: %s", err.Error())
		}

		would_be_path = filepath.Join(category_path, media_identity.Media.Name)
	}

	err = repository.MediasRepo.InsertMedia(context.Background(), media_identity.Media)
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In CreateMediaFromChunkedUpload, while inserting media", dungeon_models.ErrDB_CouldNotConnectToDB)
		return labeled_err
	}

	return nil
}

func ResizeImage(file_data []byte, media_path string, width int) (*bytes.Buffer, error) {
	img, _, err := image.Decode(bytes.NewReader(file_data))
	if err != nil {
		return nil, fmt.Errorf("Error decoding image `%s`: %s", media_path, err.Error())
	}

	image_rect := img.Bounds()
	image_width := image_rect.Dx()
	image_height := image_rect.Dy()
	height_ratio := float64(image_height) / float64(image_width)

	var thumbnail *image.RGBA

	thumbnail, ok := img.(*image.RGBA)
	if !ok {
		thumbnail = image.NewRGBA(image.Rect(0, 0, image_width, image_height))
		draw.Draw(thumbnail, thumbnail.Bounds(), img, image.Point{0, 0}, draw.Src)
	}

	// only resize if the real width is bigger than the requested width else return the original image
	if width < image_width {
		resized_width := width
		resize_height := int(float64(resized_width) * height_ratio)

		thumbnail = image.NewRGBA(image.Rect(0, 0, resized_width, resize_height))

		draw.ApproxBiLinear.Scale(thumbnail, thumbnail.Bounds(), img, img.Bounds(), draw.Over, nil)
	}

	format := filepath.Ext(media_path)

	var buf bytes.Buffer

	if format == ".jpg" || format == ".jpeg" {
		err = jpeg.Encode(&buf, thumbnail, nil)
	} else if format == ".png" {
		err = png.Encode(&buf, thumbnail)
	} else {
		return nil, fmt.Errorf("Unsupported format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("Error encoding thumbnail: %s", err.Error())
	}

	return &buf, nil
}

// creates a thumbnail for a video file. The thumbnail will be the first frame of the video, resized to the specified width. If the width is 0, the default thumbnail width will be used.
// *UNDER NO CIRCUMSTANCES* this function will upscale an image. Meaning if the original image width is 1000px and the thumbnail_width is 2000px, the resulting thumbnail will be 1000px.
func GetVideoThumbnail(f *os.File, thumbnail_width int) (*service_models.ThumbnailResponse, *dungeon_models.LabeledError) {
	var new_thumbnail_response *service_models.ThumbnailResponse = new(service_models.ThumbnailResponse)

	new_thumbnail_response.MimeType = "image/jpeg"
	new_thumbnail_response.Filename = filepath.Base(f.Name())
	new_thumbnail_response.Resized = true

	thumbnailer_ctx, err := thumbnailer.NewFFContext(f)
	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error creating Thumbnailer FFContext", dungeon_models.ErrProcessError)
	}
	defer thumbnailer_ctx.Close()

	media_dimensions, err := thumbnailer_ctx.Dims()
	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error getting Thumbnailer media dimensions", dungeon_models.ErrProcessError)
	}

	height_ratio := float64(media_dimensions.Height) / float64(media_dimensions.Width)

	if thumbnail_width == 0 {
		thumbnail_width = app_config.THUMBNAIL_WIDTH
	}

	if thumbnail_width >= int(media_dimensions.Width) {
		thumbnail_width = int(media_dimensions.Width)
	}

	thumbnail_height := int(float64(thumbnail_width) * height_ratio)

	dimensions := thumbnailer.Dims{
		Width:  uint(thumbnail_width),
		Height: uint(thumbnail_height),
	}

	thumbnail, err := thumbnailer_ctx.Thumbnail(dimensions)
	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error creating video thumbnail", dungeon_models.ErrProcessError)
	}

	var buf *bytes.Buffer = new(bytes.Buffer)

	err = jpeg.Encode(buf, thumbnail, nil)
	if err != nil {
		return nil, dungeon_models.NewLabeledError(err, "Error encoding video thumbnail", dungeon_models.ErrProcessError)
	}

	new_thumbnail_response.MediaStream = buf
	new_thumbnail_response.MediaLength = int64(buf.Len())
	new_thumbnail_response.Size = &service_models.MediaSize{
		Width:  thumbnail_width,
		Height: thumbnail_height,
	}
	new_thumbnail_response.OrignialSize = &service_models.MediaSize{
		Width:  int(media_dimensions.Width),
		Height: int(media_dimensions.Height),
	}

	return new_thumbnail_response, nil
}

func SaveMediaFile(media_identity *dungeon_models.MediaIdentity, file *multipart.File) error {
	SetUniqueMediaName(media_identity)

	media_path := path.Join(media_identity.ClusterPath, media_identity.CategoryPath, media_identity.Media.Name)

	media_file, err := os.Create(media_path)
	if err != nil {
		return fmt.Errorf("Error creating media file `%s`: %s", media_path, err.Error())
	}

	defer media_file.Close()

	_, err = io.Copy(media_file, *file)
	if err != nil {
		return fmt.Errorf("Error writing media file `%s`: %s", media_path, err.Error())
	}

	return nil
}

// Takes a media file and assumes the file is not yet on it's corresponding cluster path. It will check if there is a file with the same name as the would be media file and if so,
// it will find a unique version of the media name. doesnt make any file operations.
func SetUniqueMediaName(media_identity *dungeon_models.MediaIdentity) error {
	unique_media_name, err := GetUniquieMediaName(media_identity)
	if err != nil {
		return fmt.Errorf("Error getting unique media name: %s", err.Error())
	}

	media_identity.Media.Name = unique_media_name

	return nil
}

func GetUniquieMediaName(media_identity *dungeon_models.MediaIdentity) (string, error) {
	category_path := path.Join(media_identity.ClusterPath, media_identity.CategoryPath)
	sibling_names := make(map[string]bool)

	unique_media_name := media_identity.Media.Name

	files, err := os.ReadDir(category_path)
	if err != nil {
		return "", fmt.Errorf("Error reading category directory `%s`: %s", category_path, err.Error())
	}

	for _, file := range files {
		sibling_names[file.Name()] = true
	}

	_, name_used := sibling_names[media_identity.Media.Name]
	media_extension := filepath.Ext(media_identity.Media.Name)
	media_name := media_identity.Media.Name[:len(media_identity.Media.Name)-len(media_extension)]
	repetition := 1

	for name_used {
		unique_media_name = fmt.Sprintf("%s.v%d%s", media_name, repetition, media_extension)
		_, name_used = sibling_names[unique_media_name]
		repetition++
	}

	return unique_media_name, nil
}
