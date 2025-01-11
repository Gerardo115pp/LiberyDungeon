const MOBILE_BREAKPOINT = 768;


export const isMobile = () => {
    let is_mobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    
    if (!is_mobile && window.innerWidth < MOBILE_BREAKPOINT) {
        is_mobile = true;
    }

    
    return is_mobile;
}

/**
 * Returns as a HSL string.
 * @param {number} saturation - the saturation of the color
 * @param {number} lightness - the lightness of the color
 */
export const genRandomColor = (saturation, lightness) => {
    let hue = Math.trunc(Math.random() * 360);
    return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
}

/**
 * @param {string} payload 
 * @returns 
 */
export const createUnsecureJWT = payload => {
    /* 
        Keep in mind that this method of creating a JWT is not secure, as the JWT is not signed and could be easily tampered with. It is only suitable for passing simple parameters that do not need to be secured.
    */

    const headers = {
        alg: "none",
        typ: "JWT"
    }

    const encoded_headers = window.btoa(JSON.stringify(headers)); // stupid vscode doesnt relize we are not working in node but in the browser

    const encoded_payload = window.btoa(JSON.stringify(payload));

    return `${encoded_headers}.${encoded_payload}.`;
}

/**
 * Parses a JWT token without verifying the signature. 
 * @param {string} jwt 
 */
export const parseUnsecureJWT = jwt => {
    let parts = jwt.split('.');
    let header = JSON.parse(window.atob(parts[0]));
    let payload = JSON.parse(window.atob(parts[1]));

    return {
        header,
        payload
    }
}

// TODO: Remove attributesToJson and attributesToJsonExclusive from this file.
/**
 * @this {any}
 */
export function attributesToJson() {
    /**
     * @type {Object.<string, any>}
     */
    const json_data = {};

    Object.entries(this).forEach(([key, value]) => {
        if (!(this[key] instanceof Function) && key[0] !== '_') {
            json_data[key] = value;
        }
    });
    return JSON.stringify(json_data);
}

/**
 * @this {any}
 */
export function attributesToJsonExclusive() {
    /** @type {Object<string, any>} */
    const json_data = {};

    Object.entries(this).forEach(([key, value]) => {
        if (!(this[key] instanceof Function) && key[0] !== '_' && value !== null) {
            json_data[key] = value;   
        }
    });

    return JSON.stringify(json_data);
}

export function enterFullScreen() {
    let elem = document.documentElement;

    if (elem.requestFullscreen) {
        elem.requestFullscreen();
    // @ts-ignore
    } else if (elem.webkitRequestFullscreen) { /* Safari */
        // @ts-ignore
        elem.webkitRequestFullscreen();
        // @ts-ignore
    } else if (elem.msRequestFullscreen) { /* IE11 */
        // @ts-ignore
        elem.msRequestFullscreen();
    }
}

/**
 * @param {string} key 
 * @returns 
 */
export const getUrlPARAM = key => {
    let url_string = window.location.href; 
    url_string = url_string.replace(/\/.{0,3}#/, ""); // remove #
    let url = new URL(url_string);
    return url.searchParams.get(key);
}

/**
 * @param {string} media_url 
 * @returns 
 */
export const isUrlVideo = media_url => {
    const video_extensions = ["mp4", "webm", "ogg"];

    /** @type {string | undefined} */
    let extension = media_url.split('.').pop();

    if (extension == null) {
        return false;
    }

    extension = extension.toLowerCase();


    return video_extensions.includes(extension);
}

/**
 * @param {string} media_url 
 * @returns 
 */
export const isUrlImage = media_url => {
    const image_extensions = ["jpg", "jpeg", "png", "gif", "webp"];

    /** @type {string | undefined} */
    let extension = media_url.split('.').pop();

    if (extension == null) {
        return false;
    }
    
    extension = extension.toLowerCase();

    return image_extensions.includes(extension);
}

/**
 * @param {string} media_url 
 * @returns 
 */
export const isUrlMediaFile = media_url => {
    return isUrlVideo(media_url) || isUrlImage(media_url);
}

/**
 * Whether a given HtMLElemnt is completly invisible in the scroll view. That is, nor the bottom or the top of the element is visible. requires
 * to be runing in a window context.
 * @param {HTMLElement} element
 * @returns {boolean}
 */
export const isOutTheViewport = element => {
    if (!hasWindowContext()) {
        throw Error("isInViewport requires a window context.");
    }

    const rect = element.getBoundingClientRect();


    return rect.top <= 0 && rect.bottom <= 0 || rect.top >= window.innerHeight && rect.bottom >= window.innerHeight;    
}

/**
 * Scrolls an element into view. if the element is completly out of view, it uses instant scroll, otherwise it uses smooth scroll.
 * @param {HTMLElement} element
 * @returns {void}
 */
export const ensureElementVisible = element => {
    if (!hasWindowContext()) return;

    const element_out_of_view = isOutTheViewport(element);

    element.scrollIntoView({
        behavior: element_out_of_view ? 'instant' : 'smooth',
        block: 'center',
        inline: 'center'
    });
}

/**
 * Returns the filename from a media url 
 * @param {string} media_url 
 * @returns 
 */
export const getMediaFilename = media_url => {
    return media_url.split('/').pop();
}

/**
 * Returns the file extension from base file name such as 'file.jpg' -> 'jpg'
 * @param {string} filename 
 * @returns {string}
 */
export const getFileExtension = filename => {
    const filename_fragments = filename.split('.');

    // @ts-ignore - if the length is greater than 0, it's impossible for the return type of .pop() to be undefined. 
    return filename_fragments.length >= 2 ? filename_fragments.pop() : "";
}

/**
 * Clears a string from problematic characters for filenames. changes spaces to underscores and optionaly
 * lowercases the string.
 * @param {string} new_filename 
 * @param {boolean} [lower_case]
 */
export const cleanFilenameString = (new_filename, lower_case=true) => {
    const invalidChars = /[\]\[<>:"/\\|?*\x00-\x1F]/g
    let cleaned_filename = new_filename.replace(/\s/g, '_');
    cleaned_filename = cleaned_filename.replace(invalidChars, '');

    if (lower_case) {
        cleaned_filename = cleaned_filename.toLowerCase();
    }

    return cleaned_filename;
}

/**
 * Returns a version of the passed string that is safe to use as an id selector.
 * @param {string} unsafe_string
 * @returns {string}
 */
export const cleanIdSelector = unsafe_string => {
    return unsafe_string.replace(/[^a-zA-Z0-9-_]/g, "-");
}

/**
 * Converts a duration in seconds to a string. if the video is longer than an hour, it will be formatted as HH:MM:SS, otherwise it will be MM:SS
 * It expects the duration to be in seconds but is_seconds can be set to false to treat the duration as milliseconds.
 * @param {number} duration - the duration of the video in seconds
 * @param {boolean} is_seconds - whether the duration is in seconds or milliseconds
 * @returns {string} - the duration as a string
 */
export const videoDurationToString = (duration, is_seconds = true) => {
    let duration_string = "";

    if (!is_seconds) {
        duration = duration / 1000; // convert to seconds
    }

    let hours = Math.floor(duration / 3600);
    let minutes = Math.floor((duration % 3600) / 60);
    let seconds = Math.floor(duration % 60);

    duration_string = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;

    if (hours > 0) {
        duration_string = `${hours.toString().padStart(2, '0')}:${duration_string}`;
    }

    return duration_string;
}

/**
 * Returns the base name of a give unix-like path.
 * @param {string} path 
 * @returns {string}
 */
export const getPathBasename = path => {
    const path_fragments = path.split('/');

    if (path_fragments.length === 0) {
        return "";
    }

    let basename = "";

    for (let h = (path_fragments.length - 1); h >= 0; h--) {
        basename = path_fragments[h];

        if (basename !== "") {
            break;
        }
    }

    return basename;
}

/*=============================================
=            Data structures            =
=============================================*/
    /**
     * @template T
     */
    class StackNode {
        /** @type {T} */
        #value;
        /** @type {StackNode<T> | null} */
        #next;
        
        /**
         * @param {T} value
         */
        constructor(value) {
            this.#value = value;
            this.#next = null;
        }

        get Value() {
            return this.#value;
        }

        set Value(value) {
            this.#value = value;
        }

        get Next() {
            return this.#next;
        }

        set Next(next) {
            this.#next = next;
        }
    }

    /**
     * @template T
     */
    export class DoublyLinkedStackNode {
        /** @type {T} */
        #value;
        /** @type {DoublyLinkedStackNode<T> | null} */
        #next;
        /** @type {DoublyLinkedStackNode<T> | null} */
        #prev;

        /**
         * @param {T} value 
         */
        constructor(value) {
            this.#value = value;
            this.#next = null;
            this.#prev = null;
        }

        get Value() {
            return this.#value;
        }

        set Value(value) {
            this.#value = value;
        }

        get Next() {
            return this.#next;
        }

        set Next(next) {
            this.#next = next;
        }

        get Prev() {
            return this.#prev;
        }

        set Prev(prev) {
            this.#prev = prev;
        }
    }

    /**
     * @template T
     */
    export class Stack {
        /** 
         * @type {StackNode<T> | null} 
         */
        #top;

        /**
         * The count of the elements in the stack
         * @type {number}
         */
        #size;

        constructor() {
            this.#top = null;
            this.#size = 0;
        }

        /**
         * @param {T} value - the value to add to the stack
         */
        Add(value) {
            let new_node = new StackNode(value);
            new_node.Next = this.#top;
            this.#top = new_node;
            this.#size++;
        }

        /**
         * Sets the top node to null, which effectively clears the stack
         * @returns {void}
        */
        Clear() {
            this.#top = null;
            this.#size = 0;
        }

        /**
         * returns the value at the top of the stack without removing it
         * @returns {T | null} the value of the top node or null if the stack is empty
         */
        Peek() {
            if (this.#top === null) {
                return null;
            }

            return this.#top.Value;
        }

        /**
         * removes the top node from the stack and returns its value
         * @returns {T | null} the value of the top node or null if the stack is empty
         */
        Pop() {
            let top_node = this.#top;

            this.#top = this.#top?.Next ?? null;

            if (top_node !== null) {
                this.#size--;
            }
            
            return top_node?.Value ?? null;
        }

        IsEmpty() {
            return this.#top === null;
        }

        /**
         * The count of the elements in the stack
         * @type {number}
         */
        get Size() {
            return this.#size;
        }

        /**
         * Returns a human readable string of the stack. Mainly intended for debugging.
         * @returns {string}
         */
        toString() {
            let current_node = this.#top;
            let stack_str = "";

            while (current_node !== null) {
                stack_str += `${current_node.Value} -> `;
                current_node = current_node.Next;
            }

            stack_str += "null";

            return stack_str;
        }

    }

    /**
     * A stack with a max size. After the stack reaches the max size, adding a new element will remove the oldest element.
     * @template T
     */
    export class StackBuffer {
        /**
         * The max size of the stack
         * @type {number}
         */
        #buffer_size;

        /**
         * The current size of the stack
         * @type {number}
         */
        #size;

        /**
         * The top of the stack
         * @type {DoublyLinkedStackNode<T> | null}
         */
        #top;

        /**
         * The bottom of the stack
         * @type {DoublyLinkedStackNode<T> | null}
         */
        #bottom;

        /**
         * A node for traversing the stack from the top
         * @type {DoublyLinkedStackNode<T> | null} 
         */
        #traversing_node;

        /**
         * @param {number} buffer_size - the max size of the stack
         */
        constructor(buffer_size) {
            if (buffer_size < 0) {
                throw Error(`passed buffer_size: ${buffer_size}\nCannot create a StackBuffer with size smaller then 0.`);
            }
            this.#buffer_size = buffer_size;
            this.#size = 0;
            this.#top = null;
            this.#bottom = null;
            this.#traversing_node = null;

        }

        /**
         * Adds a value to the stack.
         * @param {T} value 
         */
        Add(value) {
            /** @type {DoublyLinkedStackNode<T>} */
            let new_node = new DoublyLinkedStackNode(value);

            new_node.Next = this.#top;

            if (this.#top != null) {
                this.#top.Prev = new_node;
            }

            if (this.#bottom === null) {
                this.#bottom = new_node
            }

            this.#size++;
            this.#top = new_node;

            this.#mantaineSize();
            return;
        }

        /**
         * Returns the stack buffer size
         * @type {number}
         */
        get BufferSize() {
            return this.#buffer_size;
        }

        /**
         * Clears the stack state.
         */
        Clear() {
            this.#top = null;
            this.#bottom = null;
            this.#size = 0;
        }

        /**
         * checks if the stack has overflowed and in so, it deletes an item from the bottom. Call from Add only
         */
        #mantaineSize() {
            if (this.#size <= this.#buffer_size) return;

            let unlinked_node = this.#bottom;

            this.#bottom = unlinked_node?.Prev ?? null;

            if (this.#bottom != null) {
                this.#bottom.Next = null;
            }

            this.#size--;
        }

        /**
         * Returns the element on the top without removing it from the stack.
         * @returns {T | null}
         */
        Peek() {
            return this.#top?.Value ?? null;
        }

        /**
         * Peeks the Nth element zero indexed element, starting from top. Where top is 0. If the element is out of bounds, returns null.
         * Keep in mind, internally this is not an array, look up time is O(n). 
         * @param {number} index
         * @returns {T | null}
         */
        PeekN(index) {
            if (index < 0 || index >= this.#size) return null;

            let current_node = this.#top;

            for (let h = 0; h < index && current_node != null; h++) {
                current_node = current_node?.Next ?? null;
            }

            return current_node?.Value ?? null;
        }

        /**
         * Peeks the current traversing node value.
         * @returns {T}
         */
        PeekTraversing() {
            if (this.#traversing_node === null && this.#top === null) {
                throw Error("Cannot peek traversing node of an empty stack.");
            }

            /**
             * @type {T}
             */
            let current_value;

            if (this.#traversing_node === null) {
                current_value = /** @type {DoublyLinkedStackNode<T>} */ (this.#top).Value;
            } else {
                current_value = this.#traversing_node.Value;
            }

            return current_value;
        }
            
        /**
         * Returns and removes the element on the top of the stack. 
         * @returns {T | null}
         */
        Pop() {
            if (this.#top === null) return null;

            if (this.#bottom === this.#top) {
                this.#bottom = null;
            }

            let current_value = this.#top.Value;        

            this.#top = this.#top?.Next ?? null;
            this.#size--;

            if (this.#top != null) {
                this.#top.Prev = null;
            }

            return current_value;
        }

        /**
         * Returns whether the stack is empty or not.
         * @returns {boolean}
         */
        IsEmpty() {
            return this.#top === null;
        }

        /**
         * Returns whether the stack is full or not.
         * @returns {boolean}
         */
        IsFull() {
            return this.#size === this.#buffer_size;
        }

        /**
         * Returns the available space before the stack starts droping older content.
         * @type {number}
         */
        get Space() {
            return this.#buffer_size - this.#size;
        }

        /**
         * Starts traversing the stack from the top.
         * @returns {void}
         */
        StartTraversal() {
            this.#traversing_node = null;
        }

        /**
         * Returns the Size of the buffer
         * @type {number}
         */
        get Size() {
            return this.#size;
        }

        /**
         * Returns a human readable string of the stack. Mainly intended for debugging.
         * @returns {string}
         */
        toString() {
            let current_node = this.#top;
            let stack_str = "";

            while (current_node !== null) {
                stack_str += `${current_node.Value} -> `;
                current_node = current_node.Next;
            }

            stack_str += "null";

            return stack_str;
        }

        /**
         * Returns the value of the current traversing node and moves to the next node.
         * Returns null if the traversal has reached the end or StartTraversal was not called before traversing.
         * @returns {T | null}
         */
        Traverse() {
            if (this.#top === null) {
                throw Error("Cannot traverse an empty stack.");
            }

            if (this.TraversalEnd) {
                return null;
            }

            if (this.#traversing_node == null) {
                this.#traversing_node = this.#top;
            } else {
                this.#traversing_node = /** @type {DoublyLinkedStackNode<T>} */ (this.#traversing_node).Next;
            }

            const current_value = /** @type {DoublyLinkedStackNode<T>} */ (this.#traversing_node).Value;

            return current_value;
        }

        /**
         * Whether the traversal has reach the end of the stack.
         * @type {boolean}
         */
        get TraversalEnd() {
            return this.#traversing_node !== null && this.#traversing_node.Next === null; 
        }
    }
/*=====  End of Data structures  ======*/


/*=============================================
=            Environment capabilities detection            =
=============================================*/

    export const hasWindowContext = () => typeof window === 'object' && globalThis === window;


    export const canUseDOMEvents = () => hasWindowContext() && typeof window.addEventListener === 'function';

/*=====  End of Environment capabilities detection  ======*/


/*=============================================
=            Feature Permission            =
=============================================*/

    /**
     * Returns whether the current context has permission to read from the clipboard. If the permission state is 'prompt', it will treat that 
     * the same way as 'denied' and return false.
     * @requires Permissions
     * @returns {Promise<boolean>} 
     */
    export const hasClipboardReadPermission = async () => {
        // @ts-ignore
        let permission = await navigator.permissions.query({name: 'clipboard-read'}); 

        return permission.state === 'granted';
    }

    /**
     * Returns whether the current context has permission to write to the clipboard. If the permission state is 'prompt', it will treat that 
     * the same way as 'denied' and return false.
     * @requires Permissions
     * @returns {Promise<boolean>} 
     */
    export const hasClipboardWritePermission = async () => {
        // @ts-ignore
        let permission = await navigator.permissions.query({name: 'clipboard-write'});
        return permission.state === 'granted';
    }

/*=====  End of Feature Permission  ======*/



