
/**
 * 
 * @param {PageData} param0 
 * @typedef {Object} PageData
 * @property {MediaExplorerParams} params
 * @typedef {Object} MediaExplorerParams
 * @property {string} category_id
 * @returns {MediaExplorerParams}
 */
export function load({ params }) {
    return {
        ...params
    }
}