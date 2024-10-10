import { GetCategoryNameAvailabilityRequest } from "./HttpRequests";

/**
 * Returns true if there is no other category with the same name on the same parent category. The verification is case insensitive. So
 * "Category" and "category" are considered the same.
 * @param {string} name - The candidate name
 * @param {string} parent_id - The category where availability will be checked
 * @returns {boolean} True if the name is available, false otherwise
 * @throws {Error} If the server response is not 200 or 409
 */
export const isCategoryNameAvailable = async (name, parent_id) => {
    const request = new GetCategoryNameAvailabilityRequest(name, parent_id);
    const response = await request.do();

    if (response.data === null) {
        throw new Error(`Unexpected server response: ${response.status}`);
    }

    return response.data;
};