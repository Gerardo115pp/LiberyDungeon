export const platform_well_known_events = {
    FS_CHANGED: 'cluster_fs_change',
    MEDIA_DELETED: "media_deleted",
    MEDIA_ADDED: "media_added",
}

/**
 * A message for the FS_CHANGED event.
* @typedef {Object} ClusterFsChangeEvent
 * @property {string} cluster_uuid
 * @property {number} medias_added
 * @property {number} medias_deleted
 * @property {number} medias_updated
*/