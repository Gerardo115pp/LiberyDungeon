
/*----------  Size  ----------*/

/**
 * Returns a human readable string of the size of a give list of MediaFile
 * @param {import('./models').MediaFile[]} files 
 * @returns {string}
 */
export const mediaFilesReadableSize = files => {
    let total_size = 0;

    for (let f of files) {
        total_size += f.size;
    }

    let size_str = humanReadableSize(total_size);

    return size_str;
}

/**
 * Returns a human readable string of a given byte size
 * @param {number} byte_size
 * @returns {string}
 */
export const humanReadableSize = byte_size => {
    let size_str = '';

    if (byte_size < 1024) {
        size_str = `${byte_size} bytes`;
    } else if (byte_size < 1024 * 1024) {
        size_str = `${(byte_size / 1024).toFixed(2)} KB`;
    } else if (byte_size < 1024 * 1024 * 1024) {
        size_str = `${(byte_size / 1024 / 1024).toFixed(2)} MB`;
    } else if (byte_size < 1024 * 1024 * 1024 * 1024) {
        size_str = `${(byte_size / 1024 / 1024 / 1024).toFixed(2)} GB`;
    } else {
        size_str = `${(byte_size / 1024 / 1024 / 1024 / 1024).toFixed(2)} TB`; // this would be ridiculous. im tempted to throw an error here.
    }

    return size_str;
}


/*----------  File type validators  ----------*/

/**
 * Checks if the given file is a video.
 * @param {File} file
 * @returns {boolean}
 */
export const isVideoFile = file => {
    return file.type.startsWith('video/');
}

/**
 * Checks if the given file is an image.
 * @param {File} file
 * @returns {boolean}
 */
export const isImageFile = file => {
    return file.type.startsWith('image/');
}

/**
 * Returns whether the file is a supported media file.
 * @param {File} file
 * @returns {boolean}
 */
export const isMediaFileSupported = file => {
    return isVideoFile(file) || isImageFile(file);
}

/**
 * Returns whether the given file is a file and not a directory.
 * @param {File} file
 * @returns {Promise<boolean>} 
 */
export const isFile = async file => {
    let is_file = true;
    let is_possibly_a_directory = file.type === '' || file.size % 4096 === 0;

    if (is_possibly_a_directory) {
        try {
            await file.slice(0, 1).arrayBuffer(); // If reading the first byte fails, it's not a file.
        } catch (e) {
            is_file = false;
        }
    }

    return is_file;
}