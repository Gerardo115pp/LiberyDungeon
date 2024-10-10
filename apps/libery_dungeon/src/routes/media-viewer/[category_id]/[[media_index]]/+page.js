import { error } from "@sveltejs/kit";

/**
 * 
 * @param {MediaViewerQueryParams} param0 
 * @typedef {Object} MediaViewerQueryParams
 * @property {string} category_id
 * @property {string} media_index
 * @returns {MediaViewerQueryParams}
 */
export function load({ params }) {
    return {
        ...params
    }
}