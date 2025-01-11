<script>
    import { ui_core_dungeon_references } from "@app/common/ui_references/core_ui_references";
    import { ui_pandasworld_tag_references } from "@app/common/ui_references/dungeon_tags_references";
    import SettingEntry from "@components/Informative/DataEntries/SettingEntry.svelte";
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import { LabeledError } from "@libs/LiberyFeedback/lf_models";
    import { emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { getPathBasename } from "@libs/utils";
    import { getMediaIdentityByUUID } from "@models/Medias";
    import { current_cluster } from "@stores/clusters";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * the current list of billboard media uuids.
         * @type {string[]}
         */ 
        export let the_billboard_media_uuids = [];

        /**
         * A map of readable media labels -> MediaIdentities. Used to allow the user to understand what medias he/she is adding.
         * @type {Map<string, import('@models/Medias').MediaIdentity>}
         */
        let media_label_to_media_identity = new Map();

        /**
         * the setting editor used to add new media uuids.
         * @type {SettingEntry}
         */
        let the_media_aggregator;

        $: if (the_billboard_media_uuids.length != 0) {
            generateLabelMap(the_billboard_media_uuids);
        }
    
    /*=====  End of Properties  ======*/
    
    
    
    /*=============================================
    =            Methods            =
    =============================================*/


        /**
         * Adds a billboard media to the media_label_to_media_identity map.
         * @param {import('@models/Medias').MediaIdentity} new_media
         * @returns {void}
         */
        const addMediaIdentityToMap = new_media => {
            const human_readable_label = generateIdentityHRLabel(new_media);

            emitPlatformMessage(`Added ${human_readable_label}`);

            const media_labels_map_copy = new Map(media_label_to_media_identity.entries());

            media_labels_map_copy.set(human_readable_label, new_media);

            media_label_to_media_identity = media_labels_map_copy;
        }

        /**
         * generates the media_label_to_media_identity map.
         * @param {string[]} new_media_uuids
         * @returns {Promise<void>}
         */
        async function generateLabelMap(new_media_uuids) {
            const media_identity_list = await $current_cluster.getClusterMedias(new_media_uuids);

            let new_media_labels_map = new Map();

            for (let media_identity of media_identity_list) {
                const human_readable_label = generateIdentityHRLabel(media_identity);

                new_media_labels_map.set(human_readable_label, media_identity);
            }

            media_label_to_media_identity = new_media_labels_map
        }

        /**
         * Generates a human readable label out of a media identity.
         * @param {import('@models/Medias').MediaIdentity} media_identity
         * @returns {string}
         */
        const generateIdentityHRLabel = media_identity => {
            const likely_category_name = getPathBasename(media_identity.CategoryPath);
            const MAX_MEDIA_NAME_LENGTH = 15;

            let human_readable_label = media_identity.Media.MediaName;

            if (human_readable_label.length > MAX_MEDIA_NAME_LENGTH) {
                human_readable_label = human_readable_label.slice(0, MAX_MEDIA_NAME_LENGTH) + "...";
            }

            if (likely_category_name !== "") {
                human_readable_label += `(${likely_category_name})`;
            }

            return human_readable_label;
        }
    
        /**
         * Handles the aggregation of the billboard medias.
         * @type {import('@components/Informative/DataEntries/data_entries').SettingChangeHandler}
         */
        const handleNewBillboardMedia = async (setting_id, new_value) => {
            const new_billboard_media_uuid = new_value;

            if (new_billboard_media_uuid === "") return;

            the_media_aggregator.clearSettingValue();

            const new_media = await getMediaIdentityByUUID(new_billboard_media_uuid);

            if (new_media === null || new_media.ClusterUUID !== $current_cluster.UUID) {
                const labeled_err = new LabeledError(
                    "In @componentss/Categories/CategoryConfiguration/sub-components/Setting__BillboardMedias.handleNewBillboardMedia",
                    `The ${ui_core_dungeon_references.MEDIA.EntityName} uuid  doesn't seem to exist.`,
                    lf_errors.ERR_HUMAN_ERROR
                )

                labeled_err.alert();
                return;
            }

            addMediaIdentityToMap(new_media);
        }
    
    /*=====  End of Methods  ======*/
    
    
</script>

<div id="cacow-billboard-medias-setting">
    <div id="cacow-billboard-medias-input-wrapper">
        <SettingEntry
            bind:this={the_media_aggregator}
            id_selector="billboard-medias-aggregator"
            font_size="var(--cacow-font-size)"
            information_entry_label="New media (uuid)"
            on_setting_change={handleNewBillboardMedia}
        />
    </div>
    <ol id="cacow-current-billboard-medias"
        class="dungeon-tag-group dungeon-tag-list" 
    >
        {#each media_label_to_media_identity as [media_uuid, identity]}
            <DeleteableItem
                item_id={media_uuid}    
                id_selector="cacow-billboard-media-{media_uuid}"
            >
                <p class="cacow-billboard-media">
                    {media_uuid}
                </p>
            </DeleteableItem>
        {:else}
            <p class="dungeon-generic-text">
                No {ui_core_dungeon_references.MEDIA.EntityNamePlural} added for this {ui_core_dungeon_references.CATEGORY.EntityName}'s billboard
            </p>
        {/each}
    </ol>
</div>

<style>
    #cacow-billboard-medias-setting {
        display: contents;
    }

    ol#cacow-current-billboard-medias {
        & p {
            font-size: var(--cacow-font-size);
        }
    }
</style>