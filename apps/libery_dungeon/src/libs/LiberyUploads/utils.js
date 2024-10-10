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