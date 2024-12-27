/*=============================================
=            GLOBAL            =
=============================================*/

    export const global_hotkey_action_triggers = {
        QUITE_CONTEXT: ["q"],
        TOGGLE_HOTKEYS_SHEET: ["?"],
        NAVIGATION_UP: ["w"],
        NAVIGATION_DOWN: ["s"],
        NAVIGATION_LEFT: ["a"],
        NAVIGATION_RIGHT: ["d"],
        ITEM_SELECTION: ["e"],
        ITEM_DELETION_NON_IMPERATIVE: ["x"],
        ITEM_DELETION_IMPERATIVE: ["del"],
        ITEM_RENAMING: ["c c"],
        ITEM_YANKING: ["y y"],
        ITEM_PASTING: ["p"],
        REGISTER_CHANGE: ["` \\c"],
        NAVIGATION_ITEM_GOTO_FINALIZER: ["g"],
    }

    export const global_hotkey_movement_triggers = [
        ...global_hotkey_action_triggers.NAVIGATION_UP,
        ...global_hotkey_action_triggers.NAVIGATION_DOWN,
        ...global_hotkey_action_triggers.NAVIGATION_LEFT,
        ...global_hotkey_action_triggers.NAVIGATION_RIGHT,
    ]

    export const global_search_hotkeys = {
        SEARCH: ["/"],
        SEARCH_NEXT: ["n"],
        SEARCH_PREVIOUS: ["N"],
    }

/*=====  End of GLOBAL  ======*/

