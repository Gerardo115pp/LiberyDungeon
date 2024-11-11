<script>
    import { role_mode_enabled, all_grants, all_roles } from "../../app_page_store";
    import { createEventDispatcher } from "svelte";
    import RolesEditor from "./sub-components/RoleGrantLinker.svelte";
    import { deleteRole, getRoleTaxonomy } from "@models/Users";
    import RolesEditorDangerZone from "./sub-components/RolesEditorDangerZone.svelte";

    /*=============================================
    =            Properties            =
    =============================================*/

        /** 
         * The given role's taxonomy.
         * @type {import("@models/Users").RoleTaxonomy | null}
         */
        let the_role_taxonomy;

        export let subject_role_label = ""; 
        $: if (subject_role_label !== "" && subject_role_label != null) {
            refreshRoleTaxonomy(subject_role_label);
        } 

        const dispatch = createEventDispatcher();

    /*=====  End of Properties  ======*/

    /*=============================================
    =            Methods            =
    =============================================*/
        
        /** 
         * Fetches the given role's taxonomy and stores it in the_role_taxonomy.
         * @param {string} role_label 
         */
        async function refreshRoleTaxonomy(role_label) {
            let role_taxonomy = await getRoleTaxonomy(role_label);

            if (role_taxonomy == null) return;

            the_role_taxonomy = role_taxonomy;
        }

        /**
         * Handles the role-grants-changed event emitted by the RoleGrantLinker component.
        */
        const handleRoleGrantsChanged = () => {
            if (the_role_taxonomy == null) return;

            refreshRoleTaxonomy(the_role_taxonomy.RoleLabel);
        }

        /**
         * Handles the role-deletion-requested event emitted by the RolesEditorDangerZone component.
         */
        const handleRoleDeletionRequested = async () => {
            if (the_role_taxonomy == null) return;

            let deleted = await deleteRole(the_role_taxonomy.RoleLabel);

            if (deleted) {
                dispatch("role-deleted");
                the_role_taxonomy = null;
            }
        }
    
    /*=====  End of Methods  ======*/

</script>

{#if $role_mode_enabled && the_role_taxonomy != null}
    <article id="roles-editor-tools">
        <header id="ret-header">
            <h2>
                Modify existing roles
            </h2>
            <article id="ret-header-instructions">
                <p class="dungeon-instructions">
                    Use this tool to remove or add grants to a role. removing grants with this tool, does not remove them from the system. all other roles that have the grants you remove will remain unaffected.
                </p>
                <p class="dungeon-instructions">
                    <span>Important</span> - When you add a grant to a role, all roles with higher hierarchy will automatically have that same grant added to them and will keep it even if you remove it from the original role.
                </p>
            </article>
        </header>
        <div id="role-editor-tool"
            class="ret-tool"
        >
            <RolesEditor 
                subject_role_taxonomy={the_role_taxonomy}
                on:role-grants-changed={handleRoleGrantsChanged}
            />
        </div>
        <div id="ret-danger-zone-wrapper">
            <RolesEditorDangerZone 
                on:role-deletion-requested={handleRoleDeletionRequested}
                {the_role_taxonomy}
            />
        </div>
    </article>
{/if}

<style>
    article#roles-editor-tools {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-4);
    }

    /*=============================================
    =            Header            =
    =============================================*/
    
        header#ret-header {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-3);

            & > h2 {
                font-size: var(--font-size-h3);
                line-height: 1;
                text-align: center;
            }

            & > article#ret-header-instructions {
                display: flex;
                flex-direction: column;
                row-gap: var(--spacing-2);
                padding: var(--spacing-1);
            }
        }
    
    /*=====  End of Header  ======*/

    .ret-tool {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-3);
    }

</style>