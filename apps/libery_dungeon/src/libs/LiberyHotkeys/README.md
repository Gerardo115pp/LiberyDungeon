# Libery Hotkeys

  *Documentation for LiberyHotkeys is a work in progress...*

A Library for handling Vim-like hotkeys, supports a robust hotkey keybind mini-language. Quickly changing keybinds using HotkeysContexts. this allows sub-components to have their own set of hotkeys when they and return the control of the hotkeys to the parent, with out needing to know who the parent is.
<br/>
##### Terminology
- **Hotkey**: A method or action meant to be triggered by a Keybind.
- **Keybind**: A key or combination of keys that trigger a Hotkey.

## Example usage
For simplicity this example will be using Svelte, but this library is compatible with any framework. or
even vanilla JS.

```html
<script>
    import { setupHotkeysManager, getHotkeysManager } from '@libs/LiberyHotkeys/libery_hotkeys';
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";

    /**
     * @type {import('@libs/LiberyHotkeys/libery_hotkeys').HotkeyContextManager | null}
     */
    let global_hotkeys_manager = null; 

    const hotkeys_context_name = "taxonomy_tags";

    onMount(() => {
        
        if (global_hotkeys_manager === null) {
            setupHotkeysManager();
            global_hotkeys_manager = getHotkeysManager();
        }
        
        if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) {
            const hotkeys_context = new HotkeysContext();

            hotkeys_context.register(["q"], () => global_hotkeys_manager.dropContext(hotkeys_context_name), {
                description: `Here you can put a description of the hotkey.`,
            });

            hotkeys_context.register(["w", "a", "s", "d"], handleNavigation, {
                description: "<navigation>Navigates the value.", 
            });


            global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);
        }

        global_hotkeys_manager.loadContext(hotkeys_context_name);
    });

    /**
     * Handles the navigation on a set of items.
     * @param {KeyboardEvent} event
     * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
     */
    const handleTagNavigation = (event, hotkey) => {
        // you'r navigation logic here.
    }
</script>
```
**Explanation:**
- `setupHotkeysManager`: function initializes both a singleton instance of a `HotkeyContextManager`(directly) and a `HotkeyBinder`(indirectly) accessible through the `getHotkeysManager` function. If called outside of a Window context, Sets a HotkeyContextManager but not the HotkeysBinder. It also defines the HotkeysContextManager as a global on globalThis. But you should never access it that way, as it's not guaranteed to be there. It only set's it up on globalThis for testing/debugging purposes. and although the behavior is not implemented, the idea is that on production builds it will be automatically removed.
### `getHotkeysManager`
function returns a singleton instance of a `HotkeyContextManager` or null.
### `HotkeysContext`

Is a collection of hotkeys that are meant to be 'turned on/off' together. It also has documentation capabilities to add information about a hotkey that can then be used to create components that show the user what hotkeys are available and what they do.

### `hotkeys_context.register`
Registers a hotkey on a the `hotkeys_context`. The hotkey the will be binded when the global hotkeys manager loads the context(by calling `.loadContext`). The hotkey can have a single trigger(one string) or multiple triggers(an array of strings). The second argument is the callback that will be triggered when a sequence of keys matches the hotkey. The third argument is an object that can have a `description`, `mode`, `await_execution`, `consider_time_in_sequence` and `can_repeat` properties. these values are optional but the description shouldn't be omitted.

### `hotkeys_context_name`: 

This is an identifier for the context. 

### `global_hotkeys_manager.dropContext`:

Unbinds all hotkeys in the context and removes it from the record of contexts. Should be called when the methods that handle a hotkey rely on references that about to te destroyed(e.g when a component is about to be unmounted/destroyed).

If you want to unbind the context but keep a reference to it(so it can be loaded again later) you can use `global_hotkeys_manager.loadPreviousContext` instead. which will unbind the current context and load the one that was loaded before it.

### `global_hotkeys_manager.declareContext vs. global_hotkeys_manager.loadContext`: 

`declareContext` is links the context to the given name, but doesn't bind the hotkeys. `loadContext` binds the hotkeys so they can be triggered by the user. But the context must already be declared.

## Hotkey keybind examples

| Case | Hotkey | Description |
| --- | --- | --- |
| `a` | triggers when only `a` is pressed |
| `shift+a` | triggers when `shift` and `a` are pressed at the same time |
| `a b` | triggers when `a` is pressed and then `b` is pressed |
| `a+b` | triggers when `a` and `b` are pressed at the same time |
| `\\l` | triggers when any letter is pressed. Matches only one letter not a sequence of letters |
| `\\l \\l` | triggers when two letters are pressed one after the other |
| `\\s enter` | triggers a sequence of letters is pressed and then the enter key is pressed |
| `\\d g` | triggers when a sequence of numbers is pressed and then the `g` key is pressed |
| `A` | triggers when the `shift+a` is pressed or `A`(with caps lock) is pressed |
| `space` | triggers when the space key is pressed |
| `tab` | triggers when the tab key is pressed |
| `alt` | WILL NOT TRIGGER. Modifier keys are dropped and added as requirements for hotkey 'fragments' |

### Some notes:

keys like `\\d`, `\\l`, `\\s` are called Metakey, and when they represent a sequence, they requirer a sequence finalizar(e.g `\\d enter`) otherwise they will match a single key press, meaning for example that `\\s` would behave as `\\s`.

hotkeys are never triggered if the current target has a `contenteditable` attribute set to true or is an input-like element(e.g input, textarea, select, etc).