const letterHotkeys = () => {
    const characters_count = Array.from(Array(26).keys());
    const abc = [...characters_count.map(i => String.fromCharCode(i + 97)), ...characters_count.map(i => String.fromCharCode(i + 65))];

    return abc
}

export const letter_hotkeys = letterHotkeys();


/*=============================================
=            Hotkeys information            =
=============================================*/

    export const HOTKEYS_HIDDEN_GROUP = "hidden";
    export const HOTKEYS_GENERAL_GROUP = "general";

/*=====  End of Hotkeys information  ======*/


/*=============================================
=            Binder Config            =
=============================================*/

    export const MAX_PAST_EVENTS = 100;
    export const MAX_TIME_BETWEEN_SEQUENCE_KEYSTROKES = 1200; // Milliseconds

/*=====  End of Binder Config  ======*/





