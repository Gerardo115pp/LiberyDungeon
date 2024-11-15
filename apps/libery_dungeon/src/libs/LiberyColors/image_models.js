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

        this.#width = image.width;
        this.#canvas.width = this.#width;

        this.#height = image.height;
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
     * @param {number} [quality=10]
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

        for (let h = 0; h < pixel_count; h+=quality) {
            const offset = h * 4;

            const red = this.#image_data.data[offset];
            const green = this.#image_data.data[offset + 1];
            const blue = this.#image_data.data[offset + 2];
            const alpha = this.#image_data.data[offset + 3];

            const color_channels_average = (red + green + blue) / 3;

            if (alpha != null && alpha >= opaque_definition && color_channels_average >= white_definition) {
                white_pixels_found++;
            }
        }

        return (white_pixels_found / (pixel_count/quality)) * 100;
    }
};
