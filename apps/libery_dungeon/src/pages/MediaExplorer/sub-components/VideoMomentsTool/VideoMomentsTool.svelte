<script>
    import { current_cluster } from "@stores/clusters";
    import { onMount } from "svelte";
    import { VideoMomentIdentity } from "@models/Metadata";
    import VideoMomentIdentityView from "./sub-components/VideoMomentIdentityView.svelte";
    import VideoMomentSearchTool from "./sub-components/VideoMomentSearchTool.svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * All the video moments in the cluster.
         * @type {VideoMomentIdentity[]}
         */
        let cluster_video_moments = []; 
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        console.log("VideoMomentsTool mounted");
       
        loadClusterVideoMoments();
    })
    
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Loads all the video moments from the current cluster 
         * into the cluster_video_moments array 
         * @returns {Promise<void>} 
         */ 
        async function loadClusterVideoMoments() {
            const cluster = $current_cluster;
            if (cluster === null) {
                console.error("@pages/MediaExplorer/sub-components/VideoMomentsTool/VideoMomentsTool.loadClusterVideoMoments: cluster is null");
                return;
            }

            if (cluster.isVideoMomentsCacheEmpty()) {
                await cluster.loadAllVideoMoments();
            }

            const video_moments = cluster.getVideoMomentsCacheItems();
            const videos_uuids = new Set(video_moments.map(vm => vm.VideoUUID));
            console.log("videos_uuids", videos_uuids);
            console.log("video_moments", video_moments);

            const video_identities = await cluster.getClusterMedias(Array.from(videos_uuids)); 
            
            const video_identities_map = new Map();

            for (let identity of video_identities) {
                video_identities_map.set(identity.Media.uuid, identity);
            }

            const video_moment_identities = [];

            for (let moment of video_moments) {
                const video_identity = video_identities_map.get(moment.VideoUUID);

                if (video_identity === undefined) {
                    throw new Error(`@pages/MediaExplorer/sub-components/VideoMomentsTool/VideoMomentsTool.loadClusterVideoMoments: video identity not found for video uuid ${moment.VideoUUID}`);
                }

                const new_video_moment_identity = new VideoMomentIdentity(video_identity, moment);

                video_moment_identities.push(new_video_moment_identity);
            }

            cluster_video_moments = video_moment_identities;
        }

        /**
         * Handles search results from the VideoMomentSearchTool.
         * @param {import('@models/Metadata').VideoMoment[]} search_results
         * @returns {void}
         */
        const handleMomentSearchResults = search_results => {
            return;
        }
    
    /*=====  End of Methods  ======*/
    
    
    
</script>

<dialog open id="dungeon-video-moments-explorer-tool"
    class="libery-dungeon-window"
>
    <header id="diviext-header">
        <h4 id="diviext-tool-title">
            Video moments
        </h4>
        <p class="dungeon-description">
            This is a list of all the video moments you have <strong>save in this cluster</strong>. Click on any of them to load the video on that video moment.
        </p>
    </header>
    <div id="diviext-search-moments">
        <VideoMomentSearchTool 
            the_video_moments={cluster_video_moments}
            onSearchResults={handleMomentSearchResults}
        />
    </div>
    <div class="diviext-moments-wrapper">
        <ul class="diviext-moments dungeon-scroll">
            {#each cluster_video_moments as video_moment}
                <li class="diviext-moment-wrapper">
                    <VideoMomentIdentityView 
                        the_video_moment_identity={video_moment}
                    />
                </li>
            {/each}
        </ul>
    </div>
</dialog>

<style>
    #dungeon-video-moments-explorer-tool {
        display: flex;
        flex-direction: column;
        container-type: size;
        background: var(--grey-t);
        width: 100cqw;
        aspect-ratio: 7 / 8;
        row-gap: var(--spacing-3);
        padding: var(--spacing-3);
        backdrop-filter: var(--backdrop-filter-blur);
    }

    
    /*=============================================
    =            header            =
    =============================================*/
    
        header#diviext-header {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-1);
            color: var(--grey-2);
        } 
    
    /*=====  End of header  ======*/

    ul.diviext-moments {
        display: flex;
        flex-direction: column;
        overflow-y: auto;
        overflow-x: hidden;
        overscroll-behavior: contain;
        background: hsl(from var(--grey-t) h s l / 0.7);
        border-radius: var(--rounded-box-border-radius);
        height: 70cqh;

        & > li.diviext-moment-wrapper:not(:last-child) {
            border-bottom: 1px solid var(--main-dark-color-7);
        }
    }
</style>