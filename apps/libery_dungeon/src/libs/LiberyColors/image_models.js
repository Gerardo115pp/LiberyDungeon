export class CanvasImage {

    /**
     * The image canvas.
     * @type {HTMLCanvasElement}
     */
    #canvas;

    /**
     * The canvas rendering context 2D.
     * @type {CanvasRenderingContext2D}
     */
    #context;

    /**
     * The image width.
     * @type {number}
     */
    #width;

    /**
     * The image height
     * @type {number}
     */
    #height;

    /**
     * The image data.
     * @type {ImageData}
     */
    #image_data;

    /**
     * @param {HTMLVideoElement | HTMLImageElement} image 
     */
    constructor(image) {
        this.#canvas = document.createElement('canvas');
        const rendering_context = this.#canvas.getContext('2d');

        if (rendering_context === null) {
            throw new Error("In CanvasImage.Constructor: Could not get a 2d context, running in IE?");
        }

        this.#context = rendering_context;

        this.#width = image instanceof HTMLVideoElement ? image.videoWidth : image.width;
        this.#canvas.width = this.#width;

        this.#height = image instanceof HTMLVideoElement ? image.videoHeight : image.height;
        this.#canvas.height = this.#height;

        
        this.#context.drawImage(image, 0, 0, this.#width, this.#height);
        this.#image_data = this.#context.getImageData(0, 0, this.#width, this.#height);
    }

    /**
     * The height of the image.
     * @type {number}
     */
    get Height() {
        return this.#height;
    }

    /**
     * Returns the one-dimensional array containing the data in RGBA order, with integer values between 0 and 255 (inclusive).
     * @type {Uint8ClampedArray}
     */
    get PixelData() {
        return this.#image_data.data;
    }

    /**
     * The width of the image.
     * @type {number}
     */
    get Width() {
        return this.#width;
    }

    /**
     * returns fuzzy logic function to determine the white percentage of an image is, that is, out of all the pixels of an image in an img tag, how many of them are whitish.
     * @param {number} [quality=10] - larger values will make the function run faster but be less accurate. minimum value is 1 and max is the pixel count of the image.
     * @returns {number}
     */
    whitePercentage = (quality) => {
        if (quality == null) {
            quality =  10;
        }

        const pixel_count = this.#width * this.#height;

        const white_definition = 200;
        const opaque_definition = 127;

        let white_pixels_found = 0;
        let pixels_checked = 0;

        quality = Math.min(4 * quality, pixel_count);

        for (let h = 0; h < pixel_count; h+=quality) {
            const offset = h * 4;

            const pixel_whitish = this.isPixelWhite(offset);

            if (pixel_whitish) {
                white_pixels_found++;
            }
        }

        return (white_pixels_found / (pixel_count/quality)) * 100;
    }

    /**
     * fuzzy logic function to determine whether a given pixel is white or not.
     * @param {number} pixel_index
     * @returns {boolean}
     */
    isPixelWhite = (pixel_index) => {
        const offset = pixel_index * 4;

        const red = this.#image_data.data[offset];
        const green = this.#image_data.data[offset + 1];
        const blue = this.#image_data.data[offset + 2];
        const alpha = this.#image_data.data[offset + 3];

        if (alpha == null || alpha < 127) return false;

        const pixel_lstar = perceivedLuminance(red, green, blue);
        
        return pixel_lstar > 20;
    }
};

/**
 * Return the luminance of a pixel. based on the sRGB color space formula. Fast but not very accurate.
 * @param {number} red
 * @param {number} green
 * @param {number} blue
 * @returns {number}
 */
const stdRgbLuminance = (red, green, blue) => {
    return (0.299 * red) + (0.587 * green) + (0.114 * blue);
}

/**
 * Return the luminance of a pixel.
 * @param {number} red - 0 to 255
 * @param {number} green - 0 to 255
 * @param {number} blue - 0 to 255
 * @returns {number}
 */
const sRGBLuminance = (red, green, blue) => {
    const vR = red / 255;
    const vG = green / 255;
    const vB = blue / 255;

    const Y = (0.2126 * sRGBToLinear(vR)) + (0.7152 * sRGBToLinear(vG)) + (0.0722 * sRGBToLinear(vB));

    return Y;
}

/**
 * Send this function a decimal sRGB gamma encoded color value
 * between 0.0 and 1.0, and it returns a linearized value.
 * @param {number}  color_channel
 * @returns {number}
 */
const sRGBToLinear = (color_channel) => {
    if (color_channel <= 0.04045) {
        return color_channel / 12.92;
    }

    return Math.pow((color_channel + 0.055) / 1.055, 2.4);
}

/**
 * Returns the perceived luminance of a pixel.
 * @param {number} red - 0 to 255
 * @param {number} green - 0 to 255`
 * @param {number} blue - 0 to 255
 * @returns {number}
 */
const perceivedLuminance = (red, green, blue) => {
    const Y = sRGBLuminance(red, green, blue);

    if (Y <= 0.008856451679) {
        return Y * 903.2962963;
    }

    return Math.pow(Y, 1/3) * 116 - 16;
}
