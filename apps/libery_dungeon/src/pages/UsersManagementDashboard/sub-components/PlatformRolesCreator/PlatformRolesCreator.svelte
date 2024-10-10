<script>
    import { getAllGrants, createGrant } from "@models/Users";
    import { role_mode_enabled, all_grants, all_roles } from "../../app_page_store";
    import { onMount, onDestroy, createEventDispatcher } from "svelte";
    import GrantCreation from "./sub-components/GrantCreation.svelte";
    import RoleCreator from "./sub-components/RoleCreator.svelte";

    /*=============================================
    =            Properties            =
    =============================================*/


        const dispatch = createEventDispatcher();

    /*=====  End of Properties  ======*/

    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Cleans the all_grants store.
         * @returns {void}
         */
        const cleanAllGrants = () => {
            all_grants.set([]);
        }

        /**
         * Handles the role created event emitted by the RoleCreator component.
         */
        const handleRoleCreated = () => {
            dispatch("platform-role-created");
        }
    
    /*=====  End of Methods  ======*/

</script>


{#if $role_mode_enabled}
    <article id="roles-creator-tools">
        <header id="rct-header">
            <h2>Create new roles and grants</h2>
        </header>
        <div id="grant-creator-wrapper"
            class="rct-tool"
        >
            <GrantCreation />
        </div>    
        <hr/>
        <div id="role-creator-wrapper"
            class="rct-tool"
        >
            <RoleCreator 
                on:role-created={handleRoleCreated}
            />
        </div>    
    </article>
{/if}

<style>
    article#roles-creator-tools {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-4);
    }

    header#rct-header {
        & > h2 {
            font-size: var(--font-size-h3);
            line-height: 1;
            text-align: center;
        }
    }

    .rct-tool {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-3);
    }

    /*=============================================
    =            Grant creation            =
    =============================================*/

    
    /*=====  End of Grant creation  ======*/
    
    
</style>