<script>
    import { ui_core_dungeon_references } from "@app/common/ui_references/core_ui_references";
    import { ui_pandasworld_tag_references } from "@app/common/ui_references/dungeon_tags_references";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        
        /*----------  Style  ----------*/
        
            /**
             * A number between 0 and 1 that is used for the alpha channel of the tool's background color.
             * @type {number}
             * @default 1
             */
            export let background_alpha = 1;
    
    /*=====  End of Properties  ======*/
    
    
</script>

<dialog open id="dungeon-tagged-medias-tool"
    class="libery-dungeon-window"
    class:section-activated={true}
    style:--background-alpha={background_alpha}
>
    <header id="dtmt-header">
        <h2 id="dtmt-header-headline">
            Filter {ui_core_dungeon_references.MEDIA.EntityNamePlural} by {ui_pandasworld_tag_references.TAG_TAXONOMY.EntityNamePlural} with multiple {ui_pandasworld_tag_references.TAG.EntityNamePlural}
        </h2>
    </header>
    <div class="dtmt-section"></div>
    <div class="dtmt-section"></div>
</dialog>