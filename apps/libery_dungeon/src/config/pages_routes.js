export const LOGIN_PAGE_PATH = "/login";
export const INITIAL_SETUP_PAGE_PATH = "/initial-setup";

export function isPublicPage(path) {
    return path == LOGIN_PAGE_PATH || path == INITIAL_SETUP_PAGE_PATH;
}