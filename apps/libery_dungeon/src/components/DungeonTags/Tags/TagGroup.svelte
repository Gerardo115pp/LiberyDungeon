<script>
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { LabeledError } from "@libs/LiberyFeedback/lf_models";
    import { createEventDispatcher } from "svelte";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * A list of dungeon tags to display.
         * @type {import("@models/DungeonTags").DungeonTag[]}
         */ 
        export let dungeon_tags = [];

        
        /*----------  State  ----------*/

            /**
             * Whether to enable tag selection by keyboard.
             * @type {boolean}
             * @default false
             */
            export let enable_keyboard_selection = false;

            /**
             * The currently focused tag index.
             * @type {number}
             */
            export let focused_tag_index = 0;
        
        /*----------  Tag creator  ----------*/

            /**
             * Whether to enable the tag creator element.
             * @type {boolean}
             */
            export let enable_tag_creator = false;

            /**
             * The tag creator input element.
             * @type {HTMLInputElement}
             */
            let the_tag_creator_input;

            /**
             * whether the new tag name is ready.
             * @type {boolean}
             */
            let new_tag_name_is_ready = false;
        
        /*----------  Style  ----------*/
            /**
             * The tag group's color.
             * @type {string}
             * @default "var(--grey-9)"
             */
            export let tag_group_color = "var(--grey-9)";

            /**
             * The color to use for the tag creator element.
             * @type {string}
             * @default "var(--grey-8)"
             */
            export let tag_creator_color = "var(--grey-5)";
        
            /**
             * Whether to display the tags indexes. These been determined by the order in which the tags were passed.
             * @type {boolean}
             * @default true
             */
            export let expose_indexes = true;
        
        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Method            =
    =============================================*/

        /**
         * Checks if the given tag name exist among the TagGroup.dungeon_tags array.
         * @param {string} tag_name
         * @returns {boolean}
         */
        const checkTagExists = tag_name => {
            return dungeon_tags.some(dtag => dtag.Name === tag_name);
        }

        /**
         * Emits the tag-created event.
         * @param {string} tag_name
         */
        const emitTagCreated = (tag_name) => {
            dispatch("tag-created", {tag_name});
        }

        /**
         * Emits the tag-deleted event.
         * @param {number} tag_id
         */
        const emitTagDeleted = (tag_id) => {
            dispatch("tag-deleted", {tag_id});
        }

        /**
         * Emits the tag-selected event.
         * @param {number} tag_id
         */
        const emitTagSelected = (tag_id) => {
            dispatch("tag-selected", {tag_id});
        }

        /**
         * Focuses the tag creator input element.
         */
        export const focusTagCreator = () => {
            the_tag_creator_input.focus();
        }

        /**
         * Handles the keydown event on the tag creator input.
         * @param {KeyboardEvent} event
         */
        const handleTagCreatorKeyDown = (event) => {
            if (event.key === "Enter") {
                event.preventDefault();
                if (!new_tag_name_is_ready) {
                    the_tag_creator_input.reportValidity();
                    return;
                };

                let new_tag_name = the_tag_creator_input.value.toLowerCase().trim();

                let tag_name_available = !checkTagExists(new_tag_name);

                if (!tag_name_available) {
                    let labeled_error = new LabeledError("In TagGroup.handleTagCreatorKeyDown", `The value '${new_tag_name}' already exists for this attribute`);

                    labeled_error.alert();
                    return;
                }

                emitTagCreated(new_tag_name);
                resetTagCreatorState();
            }

            if (event.key === "Escape") {
                the_tag_creator_input.blur();
            }
        }

        /**
         * Handles the keyup event on the tag creator input.
         * @param {KeyboardEvent} event
         */
        const handleTagCreatorKeyUp = (event) => {
            event.preventDefault();

            if (the_tag_creator_input.validationMessage !== "") {
                the_tag_creator_input.setCustomValidity("");
            }

            new_tag_name_is_ready = the_tag_creator_input.checkValidity() && the_tag_creator_input.value !== "";
            console.log("Is tag name valid: ", new_tag_name_is_ready);
        }
    
        /**
         * Handles the tag selection event.
         * @param {CustomEvent<{item_id: string}>} event
         */    
        const handleTagSelection = (event) => {
            emitTagSelected(event.detail.item_id);
        }

        /**
         * Handles the item-deleted event triggered by the DeleteableItem component.
         * @param {CustomEvent<{item_id: string}>} event
         */
        const handleTagDeletion = (event) => {
            emitTagDeleted(event.detail.item_id);
        }

        /**
         * Resets the tag creator related state.
         */
        const resetTagCreatorState = () => {
            the_tag_creator_input.value = the_tag_creator_input.defaultValue;
            new_tag_name_is_ready = false;
        }

    
    /*=====  End of Method  ======*/
    
</script>

<ol class="dungeon-tag-group dungeon-tag-list">
    {#key dungeon_tags.length}
        {#each dungeon_tags as tag, h (tag.Name)}
            {@const is_keyboard_selected = enable_keyboard_selection && h === focused_tag_index}
            <DeleteableItem 
                item_color={!is_keyboard_selected ? tag_group_color : "var(--main-dark-transparent)"}
                item_id={tag.Id}
                on:item-selected={handleTagSelection}
                on:item-deleted={handleTagDeletion}
            >
                <p class="dtg-tag-name taxonomy-member">
                    {#if expose_indexes}
                        <i class="dtg-tag-index">
                            {h + 1}
                        </i>
                    {/if}
                    <span>
                        {tag.Name}
                    </span>
                </p>
            </DeleteableItem>
        {/each}
    {/key}
    {#if enable_tag_creator}
        <DeleteableItem
            item_color={tag_creator_color}
            is_protected
        >
            <input class="tag-creator taxonomy-member" 
                bind:this={the_tag_creator_input}
                type="text"
                on:keydown={handleTagCreatorKeyDown}
                on:keyup={handleTagCreatorKeyUp}
                placeholder="New tag"
                minlength="1"
                maxlength="64"
                pattern="{'[A-z_\\d][A-z_\\-\\s\\d]{1,64}'}"
                required
            />
        </DeleteableItem>
    {/if}
</ol>

<style>

    
    /*=============================================
    =            Current members            =
    =============================================*/
    
        .taxonomy-member {
            display: flex;
            font-size: var(--font-size-1);
            line-height: 1;
            column-gap: var(--spacing-2);
        }

        p.dtg-tag-name {

            padding: 0 var(--spacing-1);

            & > i.dtg-tag-index {
                color: var(--grey-3);
                line-height: 1;
            }

            & > span {
                line-height: 1;
            }
        }
    
    /*=====  End of Current members  ======*/
    
    
    /*=============================================
    =            Tag creator            =
    =============================================*/
    
        input.tag-creator.taxonomy-member:focus {
            &:valid {
                border-bottom: var(--thin-border-width) solid hsl(from var(--success) h s l / .8);
            }
            &:invalid {
                border-bottom: var(--thin-border-width) solid hsl(from var(--warning) h s l / .9);
            }
        }
    
    /*=====  End of Tag creator  ======*/
    
    
    

</style>