


/**
 * This function implements a normal distribution using the Box-Muller transform, not amazing for performance because sqrt, log and cos are expensive operations
 * @date 5/8/2023 - 11:46:50 PM
 *
 * @param {*} mean
 * @param {*} stdev
 * @returns {*}
 */
export const gaussRandom = (mean, stdev) => {
    let u = 1 - Math.random();
    let v = Math.random();
    let z = Math.sqrt(-2.0 * Math.log(u)) * Math.cos(2.0 * Math.PI * v);

    return z * stdev + mean;
}

const getDistance = (x1, y1, x2, y2) => {
    return Math.sqrt( Math.pow(x1-x2, 2) + Math.pow(y1-y2, 2) );
}


/**
 * This particle is meant to be animated using translateX and translateY
 * 
 * @export
 * @class StaticParticle
 * @typedef {StaticParticle}
 * @property {number} x
 * @property {number} y
 * @property {number} radius
 * @property {string} color
 */
export class StaticParticle {

    
    /**
     * Creates an instance of StaticParticle.
     * @date 5/8/2023 - 11:39:08 PM
     *
     * @constructor
     * @param {number} x
     * @param {number} y
     * @param {number} radius
     * @param {string} color
     * @param {number} repulsion_factor
     */
    constructor(x, y, radius, color, repulsion_factor=1) {
        this.x = x;
        this.y = y;
        this.radius = radius;
        this.color  = color;
        this.repulsion_factor = repulsion_factor;
    }

    /**
     * checks if the particle is in a radius of repulsion_factor * radius of the particle
     * @param {StaticParticle} particle
     * @returns {boolean}
     */
    IsToClose = (particle) => {
        const distance = getDistance(this.x, this.y, particle.x, particle.y);
        return distance < (this.radius * this.repulsion_factor);
    }


}