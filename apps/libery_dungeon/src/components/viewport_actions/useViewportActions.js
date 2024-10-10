/**
 * @type {IntersectionObserver}
 */
let intersection_observer;

function ensureIntersectionObserver() {
    if (intersection_observer != null) {
        return;
    }

    /**
     * @type {IntersectionObserverInit}
     */
    const options = {
        threshold: [0, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1],
    };

    intersection_observer = new IntersectionObserver(entries => {
        entries.forEach(entry => {
            let event_name;

            const entry_height_threshold = entry.target._height_threshold ?? 0;
            // entered the viewport
            if (entry.isIntersecting && entry.intersectionRect.height >= entry_height_threshold) {
                entry.target._isIntersecting = true;
                event_name = 'viewportEnter';
            } 
            // left the viewport
            else if (!entry.isIntersecting && entry.target._isIntersecting) {
                entry.target._isIntersecting = false;
                event_name = 'viewportLeave';
            }

            if (event_name) {
                const event = new CustomEvent(event_name, { detail: {
                    ratio: entry.intersectionRatio,
                    intersection_rect: entry.intersectionRect,
                    client_rect: entry.boundingClientRect,
                }});
                entry.target.dispatchEvent(event);
            }
        });
    }, options)
}

/**
 * Removes the intersection observer
 */
export const cleanViewportObserver = () => {
    if (intersection_observer != null) {
        intersection_observer.disconnect();
    }

    intersection_observer = null;
}

/**
 * @typedef {Object} ViewportEventDetail
 * @property {number} ratio - the ratio of the intersection
 * @property {DOMRectReadOnly} intersection_rect - the intersection rect
 * @property {DOMRectReadOnly} client_rect - the client rect
 */

/**
 * @typedef {Object} ViewportActionParams
 * @property {number} height_offset - e.g 0.5 on an element with height 100px will trigger the event viewportEnter when the element center collides with the viewport top.
 */

/**
 * a svelte action that triggers the 'viewportEnter' event when the node enters the viewport and 'viewportLeave' when it leaves
 * @param {HTMLElement} node 
 * @param {ViewportActionParams} params
 * @returns 
 */
const viewport = (node, params) => {
    ensureIntersectionObserver();
    node._isIntersecting = false;

    const height_offset = params?.height_offset ?? 0;

    const node_rect = node.getBoundingClientRect();

    node._height_threshold = node_rect.height * height_offset;

    intersection_observer.observe(node);
    
    return {
        destroy: () => {
            if (intersection_observer != null) {
                intersection_observer.unobserve(node);
            }
        }
    };
}

export default viewport;