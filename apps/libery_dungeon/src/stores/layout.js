import { browser } from "$app/environment";
import { get, writable } from "svelte/store";

/**
 * If there is a resize on the viewport but is not bigger than this threshold, the layout is not considered changed
 * @type {number} - a value between 0 and 1
 */
const layout_change_threshold = 0.05;

/**
 * The most common viewport sizes by device type
 */
const VIEWPORT_SIZES = {
    MOBILE: {
        WIDTH: 375,
        HEIGHT: 667,
    },
    TABLET: {
        WIDTH: 768,
        HEIGHT: 1024
    },
    DESKTOP: {
        WIDTH: 1920,
        HEIGHT: 1080
    }
}

let LAYOUT_PROPERTIES = {
    IS_MOBILE: false,
    MOBILE_BREAKPOINT: 768,
    TABLET_BREAKPOINT: 1024,
    VIEWPORT_WIDTH: 1920,
    VIEWPORT_HEIGHT: 1080,
    SPACING: {
        SPACING_1: 0,
        SPACING_2: 0,
        SPACING_3: 0,
        SPACING_4: 0,
        SPACING_5: 0,
        SPACING_6: 0,
        SPACING_7: 0,
        SPACING_8: 0,
        SPACING_9: 0
    }
}

if (browser) {
    LAYOUT_PROPERTIES.VIEWPORT_WIDTH = window.innerWidth;
    LAYOUT_PROPERTIES.VIEWPORT_HEIGHT = window.innerHeight;
}

/**
 * @type {import('svelte/store').Writable<typeof LAYOUT_PROPERTIES>}
 * @description the layout properties of the website
 */
export const layout_properties = writable(LAYOUT_PROPERTIES);


export const defineLayout = () => {
    if (!browser) return;

    const root_styles = getComputedStyle(document.documentElement);

    let new_layout_properties = {
        ...LAYOUT_PROPERTIES
    };

    new_layout_properties.SPACING = {
        SPACING_1: parseInt(root_styles.getPropertyValue("--spacing-1")),
        SPACING_2: parseInt(root_styles.getPropertyValue("--spacing-2")),
        SPACING_3: parseInt(root_styles.getPropertyValue("--spacing-3")),
        SPACING_4: parseInt(root_styles.getPropertyValue("--spacing-4")),
        SPACING_5: parseInt(root_styles.getPropertyValue("--spacing-5")),
        SPACING_6: parseInt(root_styles.getPropertyValue("--spacing-6")),
        SPACING_7: parseInt(root_styles.getPropertyValue("--spacing-7")),
        SPACING_8: parseInt(root_styles.getPropertyValue("--spacing-8")),
        SPACING_9: parseInt(root_styles.getPropertyValue("--spacing-9")),
    }

    new_layout_properties.IS_MOBILE = isMobile();
    console.log(`called, and the answer to isMobile is ${new_layout_properties.IS_MOBILE}`);

    new_layout_properties.VIEWPORT_WIDTH = window.innerWidth;
    new_layout_properties.VIEWPORT_HEIGHT = window.innerHeight;

    layout_properties.set(new_layout_properties);
    LAYOUT_PROPERTIES = new_layout_properties;
}

export const hasChangedLayout = () => {
    const current_width = window.innerWidth;
    const current_height = window.innerHeight;

    const width_change = Math.abs(current_width - LAYOUT_PROPERTIES.VIEWPORT_WIDTH) / LAYOUT_PROPERTIES.VIEWPORT_WIDTH;
    const height_change = Math.abs(current_height - LAYOUT_PROPERTIES.VIEWPORT_HEIGHT) / LAYOUT_PROPERTIES.VIEWPORT_HEIGHT;

    return width_change > layout_change_threshold || height_change > layout_change_threshold;
}

export function isMobile() {
    let is_mobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    
    if (!is_mobile && window.innerWidth < LAYOUT_PROPERTIES.MOBILE_BREAKPOINT) {
        is_mobile = true;
    }

    
    return is_mobile;
}

/**
 * @param {boolean} is_mobile 
 */
export function setVirtualMobile(is_mobile) {

    LAYOUT_PROPERTIES.IS_MOBILE = is_mobile;
    LAYOUT_PROPERTIES.VIEWPORT_WIDTH = LAYOUT_PROPERTIES.IS_MOBILE ? VIEWPORT_SIZES.MOBILE.WIDTH : VIEWPORT_SIZES.DESKTOP.WIDTH;
    LAYOUT_PROPERTIES.VIEWPORT_HEIGHT = LAYOUT_PROPERTIES.IS_MOBILE ? VIEWPORT_SIZES.MOBILE.HEIGHT : VIEWPORT_SIZES.DESKTOP.HEIGHT;

    layout_properties.set(LAYOUT_PROPERTIES);
}


/*=============================================
=            Layout elements            =
=============================================*/

/**
 * @type {import('svelte/store').Writable<boolean>} whether the navbar is hidden or not, means the navbar css visibility property is set to hidden or visible
 */
export const navbar_hidden = writable(false);

/**
 * whether to show the hotkeys help modal or not. THIS IS NOT THE HOTKEYS SHEET. this is the little floating button labeled '?' that TRIGGERS the hotkeys sheet on hover.
 * @type {import("svelte/store").Writable<boolean>} 
 * @default true
 */
export const hotkeys_modal_visible = writable(true);

/**
 * @type {import("svelte/store").Writable<boolean>} whether to show the hotkeys sheet table or not
 * @default false
 */
export const hotkeys_sheet_visible = writable(false);

export const toggleHotkeysSheet = () => hotkeys_sheet_visible.set(!get(hotkeys_sheet_visible))


/*=============================================
=            Utils            =
=============================================*/

/**
 * Dark mode is not a dark theme, it just hints at all components that they should use black as max dark color instead of
 * var(--grey)(which is hsl(0, 0%, 0%)). if force_enable is false, it will toggle the dark mode, if true, it will set it 
 * to true regardless of the current state.
 * @param {boolean} force_enable
 * @returns {void}
 */
export const toggleDarkMode = force_enable => {
    if (!browser) return;

    let root = document.documentElement;

    const dark_mode_enabled = root.classList.contains("dark-mode");

    if (dark_mode_enabled && !force_enable) {
        root.classList.remove("dark-mode");
    } else if (!dark_mode_enabled) {
        root.classList.add("dark-mode");
    } 
}

/**
 * Returns the state of dark mode
 * @returns {boolean}
 */
export const inDarkMode = () => {
    if (!browser) return false;

    let root = document.documentElement;

    return root.classList.contains("dark-mode");
}

/**
 * Cinema mode is a hint to all components that they should remove any non-essential elements from the ui, like the navbar and prioritize the media content visibility.
 * @param {boolean} [force_enable]
 */
export const toggleCinemaMode = force_enable => {
    if (!browser) return;

    let root = document.documentElement;

    const cinema_mode_enabled = root.classList.contains("cinema-mode");

    if (cinema_mode_enabled && !force_enable) {
        root.classList.remove("cinema-mode");
    } else if (!cinema_mode_enabled) {
        root.classList.add("cinema-mode");
    }
}

/**
 * Returns the state of cinema mode
 * @returns {boolean}
 */
export const inCinemaMode = () => {
    if (!browser) return false;

    let root = document.documentElement;

    return root.classList.contains("cinema-mode");
}

/*=====  End of Utils  ======*/

