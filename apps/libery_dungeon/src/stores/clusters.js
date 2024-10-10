import { get, writable } from "svelte/store";
import { CategoriesCluster } from "@models/CategoriesClusters";
import { browser } from "$app/environment";


/**
 * @type {import('svelte/store').Writable<CategoriesCluster>}
 */
export const current_cluster = writable(null);

/**
 * The key to be used in the local storage to store the current cluster.
 */
const CURRENT_CLUSTER_KEY = 'current_cluster';


/**
 * Persist the current cluster in the local storage.
 * @param {CategoriesCluster} cluster 
 */
const persistCluster = cluster => {
    if ((get(current_cluster) === null || cluster.UUID == null) || !browser) {
        return;
    }

    localStorage.setItem(CURRENT_CLUSTER_KEY, JSON.stringify(cluster));
};

/**
 * Load the current cluster from the local storage. if a cluster is found then it will set it to the current_cluster store. returns true if a cluster was found, false otherwise.
 * @returns {boolean}
 */
export const loadCluster = () => {
    if (!browser) return false;

    const stored_value = localStorage.getItem(CURRENT_CLUSTER_KEY);

    if (stored_value == null) {
        return false;
    }

    const cluster = JSON.parse(stored_value);

    try {
        const new_cluster = new CategoriesCluster(cluster);
        current_cluster.set(new_cluster);
        console.log("Loaded cluster from local storage: ", new_cluster);
    } catch (error) {
        console.error("Error parsing cluster from local storage: ", error);
        return false;
    }

    return true;
}

current_cluster.subscribe(persistCluster);
