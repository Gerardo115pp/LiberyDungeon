/**
 * A base class for all services responses
 * @template T
 */
export class HttpResponse {
    /**
     * @type {number}
     */
    status;

    /**
     * @type {T}
     */
    data;

    /**
     * @type {Headers}
     */
    headers;

    /**
     * @param {Response} response
     * @param {T} data
     */
    constructor(response, data) {
        this.status = response.status;
        this.data = data;
        this.headers = response.headers;
    }

    /**
     * Returns true if the response status larger than or equal to 200 and less than 300
     * @returns {boolean}
     */
    get Ok() {
        return this.status >= 200 && this.status < 300;
    }
    
    /**
     * indicates that the request has succeeded and has led to the creation of a resource. The new resource, or a description and link to the new resource, is effectively created before
     * the response is sent back and the newly created items are returned in the body of the message, located at either the URL of the request, or at the URL in the value of the 
     * Location header.
     * @returns {boolean}
     */
    get Created() {
        return this.status === 201;
    }
}

/**
* @typedef {Object} ReasonedBooleanResponse
 * @property {boolean} response
 * @property {string} reason
 */
/**
* @typedef {Object} BooleanResponse
 * @property {boolean} response
 */
/**
* @typedef {Object} SingleStringResponse
 * @property {string} response
 */
/**
* @typedef {Object} SingleNumberResponse
 * @property {number} response
*/

export function attributesToJson() {
    const json_data = {};
    console.log("AttributestoJson:" + this);
    Object.entries(this).forEach(([key, value]) => {
        if (!(this[key] instanceof Function) && key[0] !== '_') {
            json_data[key] = value;
        }
    });
    return JSON.stringify(json_data);
}

export function attributesToJsonExclusive() {
    const json_data = {};
    Object.entries(this).forEach(([key, value]) => {
        if (!(this[key] instanceof Function) && key[0] !== '_' && value !== null) {
            json_data[key] = value;   
        }
    });

    return JSON.stringify(json_data);
}