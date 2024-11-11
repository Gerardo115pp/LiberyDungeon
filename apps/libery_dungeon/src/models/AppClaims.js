import { GetCategoriesClusterSignAccessRequest } from "@libs/DungeonsCommunication/services_requests/categories_cluster_requests";

class ClaimsGrantResponse {
    /**
     * @param {{redirect_url: string, granted: boolean}} param0 
     */
    constructor({redirect_url, granted}) {
        this.redirect_url = redirect_url;
        this.granted = granted;
    }
}

/**
 * Sends a request to get an http only jwt cookie which allows the user to access the cluster.
 * @param {string} cluster_id 
 * @returns {Promise<ClaimsGrantResponse | null>}
 */
export const getCategoriesClusterSignAccess = async cluster_id => {
    const access_request = new GetCategoriesClusterSignAccessRequest(cluster_id);

    const access_response = await access_request.do();

    /**
     * @type {ClaimsGrantResponse | null}
     */
    let grant_response = null;

    if (access_response.status === 200) {
        try {
            grant_response = new ClaimsGrantResponse(access_response.data);
        } catch (error) {
            console.error("Error parsing response data", error, access_response.data);
        }
    }

    return grant_response;
}