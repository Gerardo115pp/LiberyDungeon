<script>
    import { current_cluster } from "@stores/clusters";
    import { getCategoriesClusterSignAccess } from "@models/AppClaims";
    import { current_category, resetCategoriesTreeStore } from "@stores/categories_tree";
    import { goto } from "$app/navigation";
    import { updateCluster } from "@models/CategoriesClusters";
    import { emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The cluster this item is representing.
         * @type {import("@models/CategoriesClusters").CategoriesCluster} 
         */    
        export let cluster = {};

        /**
         * Dungeon door animation duration.
         * @type {number}
         */
        const door_animation_duration = 1000;

        /**
         * The main dom element of this component.
         * @type {HTMLDivElement}
         */
        let this_dom_element;

        /**
         * The renaming input element.
         * @type {HTMLInputElement}
         */
        let rename_input_element;

        /**
         * Whether the cluster is been focused by keyboard selection.
         * @type {boolean}
         */
        export let keyboard_focused = false;
        $: if (keyboard_focused && this_dom_element != null) {
            ensureElementIsVisible(this_dom_element);
        }

        
        /*=============================================
        =            State            =
        =============================================*/
        
            /**
             * Whether the cluster is being renamed.
             * @type {boolean}
             */
            let cluster_renaming = false;

            /**
             * Whether the dungeon has been clicked/selected.
             * @type {boolean}
             */
            let dungeon_selected = false;

            /**
             * Has the fist keydown event been triggered on the input field? used to prevent the hotkey from been added as part of the category name.
             * @type {boolean}
             */ 
            let first_keydown_triggered = false;            

        /*=====  End of State  ======*/
    
    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Ensures the given element is visible in the viewport.
         * @param {HTMLElement} element 
         */
        const ensureElementIsVisible = (element) => {
            element.scrollIntoView({
                behavior: "smooth",
                block: "center",
                inline: "center"
            });
        }
    
        const handleDungeonSelected = () => {
            dungeon_selected = true;
            current_cluster.set(cluster);

            if ($current_category != null && $current_category.ClusterUUID != cluster.UUID) {
                resetCategoriesTreeStore();
            }

            setTimeout(redirectToDungeon, door_animation_duration * 0.9);
        }

        /**
         * Prevents the the hotkey keystroke from being added to the category name.
         * @param event
         */
        const handleKeyPrevent = (event) => {
            if (!first_keydown_triggered) {
                event.preventDefault();
                first_keydown_triggered = true;
            }
        }

        const handleDungeonRenameRequest = () => {
            console.log(`Dungeon rename request on '${cluster.Name}'`);

            cluster_renaming = true;
        }

        /**
         * Handles the keydown event on the rename input.
         * @param {KeyboardEvent} event 
         */
        const handleRenameInput = (event) => {
            if (!first_keydown_triggered) {
                event.preventDefault();
                first_keydown_triggered = true;
                return;
            }

            let name_valid = validateClusterName();

            if (event.key === "Enter") {
                event.preventDefault();
                event.stopPropagation();
                if (name_valid) {
                    cluster_renaming = false;
                    first_keydown_triggered = false;
                    renameCluster(rename_input_element.value);
                }
            }

            if (event.key === "Escape") {
                event.preventDefault();
                event.stopPropagation();
                cluster_renaming = false;
                first_keydown_triggered = false;
            }
        }

        const redirectToDungeon = async () => {
            const cluster_access = await getCategoriesClusterSignAccess(cluster.UUID);

            if (cluster_access == null || !cluster_access.granted) {
                console.error("Access denied to this cluster.");
                return;
            }

            // window.location.href = cluster_access.redirect_url;
            goto(cluster_access.redirect_url);
        }

        const renameCluster = async (new_name) => {
            let original_name = cluster.Name;
            console.log(`Renaming '${cluster.Name}' to '${new_name}'`);

            cluster.Name = new_name;

            let renamed = await updateCluster(cluster);

            if (!renamed) {
                console.error("Error renaming cluster.");
                cluster.Name = original_name;
                return;
            }
            
            emitPlatformMessage(`Cluster '${original_name}' renamed to '${new_name}'`);
        }

        /**
         * Checks the validity of the rename input. If rename_input_element is null, it will return false.
         * @returns {boolean} 
         */
        const validateClusterName = () => {
            if (rename_input_element == null) return false;

            return rename_input_element.checkValidity();
        }
    
    /*=====  End of Methods  ======*/

</script>

<div class="categories-cluster-item"
    bind:this={this_dom_element}
    class:is-keyboard-focused={keyboard_focused}
    class:adebug={false}
    on:click={handleDungeonSelected}
    on:cluster-rename={handleDungeonRenameRequest}
>
    <svg viewBox="0 0 191.54 235.04">
        <g 
            class="door"
            class:open={dungeon_selected}
            style:animation-duration="{door_animation_duration}ms"
        >
            <path class="door-background" d="M47 87A45.5 28 0.0 0 1 138 87L138 228L47 228Z"></path>
            <line class="door-bar door-reinforcement" x1="47.62" y1="123.13" x2="138.85" y2="123.13"/>
            <line class="door-bar door-reinforcement" x1="47.62" y1="132.59" x2="138.85" y2="132.59"/>
            <g class="door-bar-one">
                <line class="door-bar" x1="67.38" y1="73.06" x2="67.38" y2="120.36"/>
                <line class="door-bar" x1="62.04" y1="73.02" x2="62.04" y2="120.32"/>
                <line class="door-bar" x1="62.04" y1="137.36" x2="62.04" y2="172.91"/>
                <line class="door-bar" x1="67.38" y1="137.4" x2="67.38" y2="172.95"/>
                <line class="door-bar" x1="62.04" y1="187.31" x2="62.04" y2="222.86"/>
                <line class="door-bar" x1="67.38" y1="187.35" x2="67.38" y2="222.9"/>
            </g>
            <g class="door-bar-two">
                <line class="door-bar" x1="87.06" y1="64.43" x2="86.83" y2="120.41"/>
                <line class="door-bar" x1="81.72" y1="64.39" x2="81.49" y2="120.36"/>
                <line class="door-bar" x1="81.49" y1="137.4" x2="81.49" y2="172.95"/>
                <line class="door-bar" x1="86.83" y1="137.44" x2="86.83" y2="172.99"/>
                <line class="door-bar" x1="81.49" y1="187.35" x2="81.49" y2="222.9"/>
                <line class="door-bar" x1="86.83" y1="187.39" x2="86.83" y2="222.94"/>
            </g>
            <g class="door-bar-four">
                <line class="door-bar" x1="128.87" y1="73.14" x2="128.87" y2="120.45"/>
                <line class="door-bar" x1="123.54" y1="73.1" x2="123.54" y2="120.41"/>
                <line class="door-bar" x1="123.54" y1="137.44" x2="123.54" y2="172.99"/>
                <line class="door-bar" x1="128.87" y1="137.48" x2="128.87" y2="173.03"/>
                <line class="door-bar" x1="123.54" y1="187.39" x2="123.54" y2="222.94"/>
                <line class="door-bar" x1="128.87" y1="187.43" x2="128.87" y2="222.98"/>
            </g>
            <g class="door-bar-three">
                <line class="door-bar" x1="102.36" y1="64.39" x2="102.13" y2="120.32"/>
                <line class="door-bar" x1="109.54" y1="120.32" x2="109.54" y2="64.31"/>
                <line class="door-bar" x1="102.13" y1="137.36" x2="102.13" y2="172.91"/>
                <line class="door-bar" x1="110.81" y1="137.36" x2="110.81" y2="172.91"/>
                <line class="door-bar" x1="102.13" y1="187.31" x2="102.13" y2="222.86"/>
                <line class="door-bar" x1="110.81" y1="187.31" x2="110.81" y2="222.86"/>
            </g>
            <line class="door-bar door-reinforcement" x1="138.85" y1="175.73" x2="47.62" y2="175.73"/>
            <line class="door-bar door-reinforcement" x1="47.62" y1="183.44" x2="138.85" y2="183.44"/>
        </g>
        <path class="stone exalted-stone" d="m77.12.51h29.08s9.33-.59,6.58,7.7l-3.84,37.63-.82,2.07-.55,5.04s-1.37,3.56-2.19,2.67-23.32,0-23.32,0c0,0-2.74.59-3.02-3.56s-4.12-35.85-4.12-35.85l-1.77-12.29s.04-1.08.34-1.16c.9-.25,2.8-.92,3.63-2.25Z"/>
        <path class="stone" d="m70.84,9.8c.47-.05.89.32.95.82l5.04,42.44s.21.89-1.54,2c-1.48.94-8.96,3.16-11.2,3.82-.38.11-.77-.06-.97-.42l-1.47-2.61c-.03-.05-.06-.1-.1-.15l-2.31-2.74c-.07-.08-.12-.17-.16-.27l-10.69-27.85s-1.22-2.09-.56-4.19c.31-.98,1.35-1.77,2.53-2.69s6.09-3.77,11.04-6c3.08-1.39,7.6-1.96,9.45-2.15Z"/>
        <path class="stone" d="m115.51,9.73l-1.16,1.16-4.81,43.84s.72,2.33,1.54,2.11,9.26,4.22,9.26,4.22c0,0,3.09,1.56,3.6-1l13.38-39.67-.62-4.89s-19.13-7.93-21.2-5.78Z"/>
        <path class="stone" d="m141.84,21.71l-14.85,38.84c-.22.57-.04,1.23.42,1.59l13.89,10.63c.54.41,1.27.29,1.68-.26l20.64-31.86c.33-.46.35-1.1.04-1.58-1.09-1.66-3.05-6.09-7.4-10.01-3.74-3.37-7.01-6.23-12.93-8.12-.6-.19-1.24.15-1.48.78Z"/>
        <path class="stone exalted-stone" d="m143.76,75.78l23.9-29.43c.13-.16.32-.25.52-.25h1c.17,0,.34.07.47.2l12.78,16.19c.08.08.14.18.18.28l4.81,22.88s.05.19.05.29l-1.44,1.2c0,.39-.28.72-.64.76l-40.55.36c-.27.03-.53-.11-.67-.36-.86-1.5-.48-6.5-2.84-9.2-.71-.82,1.73-2.09,2.44-2.91Z"/>
        <path class="stone" d="m59.24,59.51s0,2.02-1.63,3.18c-1.49,1.06-9.19,8.44-9.19,8.44,0,0,.27.89-2.19,0l-27.82-22.32-.96-3.56s10.81-14.27,22.06-19.16l4.94.74,14.8,32.67Z"/>
        <path class="stone exalted-stone" d="m44.28,75.22l-31.36-22.6s-3.57-1.21-4.66,2.62c-.49,1.7-3.52,8.84-4.39,16.25-.87,7.37-.36,15.01-.11,17.87.06.66.59,1.14,1.2,1.1l35.68-2.7s1.65.41,1.65-3.44c0-2.13.47-4.4.89-6,.33-1.25,2.21-2.57,1.12-3.11Z"/>
        <path class="stone" d="m40.22,91.95l-27.03,3.54c-.23.02-.45.09-.65.2l-4.83,2.56c-.45.24-.73.73-.73,1.27l-.09,16.22c0,.72.54,1.3,1.2,1.3l31.13,1.81c.62,0,1.19-.36,1.5-.93l.94-1.73v-21.79c0-.42-.14-.83-.4-1.14l-1.04-1.29Z"/>
        <path class="stone" d="m8.58,120.33l32.52,1.52c.19,0,.32.19.29.39l-1.05,6.24c-.07.44,0,.89.21,1.27.28.52.66,1.38.9,2.56.41,2-.46,5.5-.31,8.72s-.62,3.83-.62,3.83l-31.87-1.91c-.32-.02-.59-.26-.68-.59l-.87-3.32c-.08-.3-.12-.62-.12-.93v-15.47s-.11-1.1.66-1.94c.24-.26.59-.39.94-.37Z"/>
        <path class="stone" d="m41.56,149.5l-29.58-2.65c-.1,0-.19.05-.24.14l-1.29,2.51h-3.11c-.2,0-.36.18-.36.39v19.11c0,.93.51,1.78,1.31,2.14l2.78,1.29c.68.31,1.41.48,2.14.48h26.32l2.47-6.22-.44-17.19Z"/>
        <path class="stone" d="m42,179.08s1.37-3.56-5.08-2.37l-25.24-.24h-3.7s-2.19-.05-1.78,2.61l-.05,18.57,3.89,4.3h31.14c.45,0,.82-.4.82-.89v-21.98Z"/>
        <path class="stone exalted-stone" d="m42,210.91l-3.03-4.91c-.25-.4-.66-.64-1.1-.64l-32.77-.28-1.51.74-2.09.25c-.57.07-1,.59-1,1.21l.21,22.02c0,.32.09.63.27.89l.78,1.16,3.54,2.79c.32.25.71.39,1.11.39l33.4-2.42c.44,0,.87-.19,1.18-.53l.61-.66c.35-.38.55-.89.54-1.43l.16-17.59c0-.36-.1-.7-.29-1Z"/>
        <path class="stone" d="m146.76,91.95h27.76l1.6,1.43,6.15,1.77c.61.18,1.04.77,1.04,1.46v17.34c0,.7-.45,1.31-1.09,1.46l-5.27,1.23-30.25-.42c-.76-.03-1.43-.55-1.69-1.32l-.84-2.43v-17.72c0-1.55,1.16-2.8,2.59-2.8Z"/>
        <path class="stone" d="m145.3,137.06l2.23-6.44-.1-11.22,33.68.33,1.37,3.26.82,4.96v10.59s-1.37,2.22-2.47,2.07c-1.04-.14-31.81,0-35.23,0-.16,0-.3-.14-.3-.32v-3.24Z"/>
        <path class="stone" d="m146.12,144.25l35.67-.52,1.71,5.76v15.28l-1.92,2.74-32.72-.07-3.47-3.19-.32-18.82c0-.63.46-1.15,1.04-1.18Z"/>
        <path class="stone" d="m146.25,172.29l17.48-1.75s17.78.74,18.2.59c.36-.13,2.13.43,3.1,1.39.21.21.33.51.32.82l-.61,20.84c-.03,1.1-.88,1.97-1.9,1.94l-16.95-.55-9.88.76-8.99-.5c-.42-.01-.76-.38-.76-.84l-.96-19.72v-1.93c0-.57.42-1.04.95-1.06Z"/>
        <path class="stone exalted-stone" d="m146.44,204.77l6.14-2.82,11.16-1.19,20.72.16s3.84.42,6.58,3.68l-.1,23.31-3.45,3.73c-.95,1.03-2.24,1.61-3.59,1.61l-37.07-2.45c-.85-.06-1.52-.82-1.53-1.74l-.21-22.48c0-.87.56-1.63,1.35-1.81Z"/>
    </svg>
    <div class="cct-cluster-label">
        {#if !cluster_renaming}
            <h4 class="cci-cluster-name">
                {cluster.Name}
            </h4>
        {:else}
            <input class="ce-ic-rename-input" 
                type="text"
                bind:this={rename_input_element}
                on:keydown={handleRenameInput} 
                on:keypress={handleKeyPrevent}
                on:click|stopPropagation={() => {}}
                value="{cluster.Name}"
                maxlength="200"
                pattern="{'[\\p{L}a-zA-Z\\d\\s_\\-]{2,}'}"
                required
                autofocus
            />
        {/if}
    </div>
</div>

<style>
    .categories-cluster-item {
        display: flex;
        flex-direction: column;
        gap: var(--vspacing-2);
        perspective: 1000px;
    }
    
    
    /*=============================================
    =            Dungeon door icon            =
    =============================================*/

        .categories-cluster-item > svg path, .categories-cluster-item > svg line {
            stroke: var(--main);
            stroke-width: 1.2;
            stroke-linecap: round;
            stroke-linejoin: round;
        }

        .categories-cluster-item .stone {
            stroke: var(--main-dark);
            background: var(--body-bg-color);
        }

        .categories-cluster-item .stone.exalted-stone {
            stroke: var(--main-dark-color-7);
            stroke-width: 1.9;
        }

        .categories-cluster-item .door-bar {
            stroke: var(--main-8);
            stroke-width: 1;
        }

        .categories-cluster-item .door-bar.door-reinforcement{
            stroke: var(--main-7);
        } 

        .categories-cluster-item .door.open {
            animation-name: door-open;
            animation-duration: 1s;
            animation-timing-function: ease-out;
            animation-fill-mode: forwards;
            animation-iteration-count: 1;
            transform-box: fill-box;
            transform-style: preserve-3d;
            transform: matrix3d(
                1, 0.2, 0, 0.01,
                0, 1, 0, 0,
                0, 0, 1, 0.003,
                0, 0, 0, 1
            );
        }

        @keyframes door-open {
            0% {
                transform: matrix3d(
                    1, 0.1, 0, 0,
                    0, 1, 0, 0,
                    0, 0, 1, 0,
                    0, 0, 0, 1
                );
            }
            20% {
                transform: matrix3d(
                    1, 0.2, 0, 0.002,
                    0, 1, 0, 0,
                    0, 0.2, 1, 0,
                    0, 0, 0, 1
                );
            }
            100% {
                transform: matrix3d(
                    1, 0.2, 0, 0.01,
                    0, 1, 0, 0,
                    0, 0, 1, 0.003,
                    0, 0, 0, 1
                );
            }
        }

        .categories-cluster-item .door .door-background {
            opacity: 0;
        }

        @supports (color: rgb(from white r g b)) {
            .categories-cluster-item .stone {
                fill: hsl(from var(--main) h s l / 0.06);
                transition: all 500ms ease-out;
            }

            .categories-cluster-item .stone.exalted-stone {
                fill: hsl(from var(--main) h s l / 0.09);
            } 

            .categories-cluster-item .door .door-background {
                opacity: 1;
                fill: transparent;
                stroke: transparent;
                transition: all 800ms ease-out;
            }

            .categories-cluster-item:hover .door:not(.open) .door-background {
                opacity: 1;
                fill: hsl(from var(--danger-8) h s l / 0.15);
                stroke: none;
            }

            /* .categories-cluster-item:hover:has(.door:not(.open)) .stone {
                fill: hsl(from var(--main) h s l / 0.05);
            } */
        }
        
    
    
    /*=====  End of Dungeon door icon  ======*/

    
    /*=============================================
    =            Label            =
    =============================================*/

        .categories-cluster-item .cct-cluster-label {
            width: 100%;
        }
    
        .cct-cluster-label h4.cci-cluster-name {
            font-size: var(--font-size-p-small);
            font-family: var(--font-read);
            color: var(--grey-2);
            text-align: center;
        }

        .cct-cluster-label input.ce-ic-rename-input {
            width: 100%;
            background: transparent;
            color: var(--main-3);
            outline: none;

            &:invalid {
                color: var(--danger) !important;
            }
        }
    
    
    /*=====  End of Label  ======*/
    
    
    
    
</style>