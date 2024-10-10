<script>
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * Whether to show the path elegibility rules dialog or not
         * @type {boolean}
         */
        export let show_path_elegibility_rules = false;
    
    /*=====  End of Properties  ======*/
    
</script>

<dialog open={show_path_elegibility_rules} id="cct-directory-eligibility-rules-dialog">
    <ol id="cct-directory-eligibility-rules-list">
        <li class="cct-derl-law law-title">
            <p>The path cannot be the root directory aka <code>'/'</code> or <code>'C:\'</code> in windows.</p>
        </li>
        <li class="cct-derl-law law-title">
            <p>It must be under but not equal to the directory defined by the app_config.SERVICE_CLUSTERS_ROOT setting.</p>
        </li>
        <li class="cct-derl-law law-title">
            <p>A non-existing directory can be passed but it's direct parent must exist.</p>
        </li>
        <li class="cct-derl-law law-title">
            <p>If the path exists, it must be a directory.</p>
        </li>
        <li class="cct-derl-law law-title">
            <p>The path cannot be used by another dungeon. This means:</p>
            <ol class="cct-derl-minor-laws">
                <li class="minor-law law-title">
                    The path is not equal to another dungeon's path.
                </li>
                <li class="minor-law law-title">
                    The path is not a parent of another dungeon's path.
                </li>
                <li class="minor-law law-title">
                    The path is not a child of another dungeon's path.
                </li>
            </ol>
        </li>
    </ol>
</dialog>

<style>
    dialog#cct-directory-eligibility-rules-dialog {
        position: absolute;
        inset: 100% 0 auto auto;
        width: min(100%, 500px);
        background: var(--grey);
        color: var(--grey-1);
        padding: var(--spacing-2) var(--spacing-3);
        z-index: var(--z-index-t-2);

        & ol {
            counter-reset: cct-derl-law;
            display: flex;
            list-style-type: none;
            flex-direction: column;
        }

        & ol li.minor-law::before, & ol li > p::before {
            counter-increment: cct-derl-law;
            content: counters(cct-derl-law, ".") " - ";
            color: var(--grey-4);
        }

        & code {
            color: var(--main);
        }
    }

    ol#cct-directory-eligibility-rules-list {
        gap: var(--spacing-1);
    }

    ol.cct-derl-minor-laws {
        gap: calc(0.5 * var(--spacing-1));
        padding-inline-start: var(--spacing-2);
    }

    li.law-title {
        font-size: calc(var(--font-size-fineprint) * 1.2);

        & > p {
            font-size: calc(var(--font-size-1) * .9);
        }
    }
    
</style>