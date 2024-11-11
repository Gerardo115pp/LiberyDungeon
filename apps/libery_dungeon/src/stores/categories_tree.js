import { get, writable } from "svelte/store";
import { CategoryLeaf, CategoriesTree } from "@models/Categories";
import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
import { replaceState } from "$app/navigation";
import { lf_errors } from "@libs/LiberyFeedback/lf_errors";


/** 
 * Category currently selected by the user
 * @type {import('svelte/store').Writable<CategoryLeaf | null>} 
 */
export const current_category = writable(null);

/**
 * A category UUID that is meant to be pasted/moved to another category.
 * @type {import('svelte/store').Writable<string>}
 * @default ""
 */
export const yanked_category = writable("");

/**
 * @type {import('svelte/store').Writable<CategoriesTree | null>} 
 * the categories tree
 */
export const categories_tree = writable(null);

export const resetCategoriesTreeStore = () => {
    current_category.set(null);
    categories_tree.set(null);
}

/**
 * @returns {Promise<void>}
 */
export const navigateToParentCategory = async () => {
    let category = get(current_category);

    if (category == null) {
        throw new Error("Cannot use navigateToParentCategory if current_category is null");
    }

    if (category.parent === "") return;

    const the_tree = get(categories_tree);

    if (the_tree == null) {
        throw new Error("Cannot use navigateToParentCategory if categories_tree is null");
    }

    // @ts-ignore
    replaceState(`/dungeon-explorer/${category.parent}`);

    the_tree.navigateToParent();

    return;
}

/**
 * Navigates to unconnected category by changing the root of the tree to that category. this also means
 * loosing all the loaded categories. Returns true if successful, false otherwise.
 * @param {string} category_uuid
 * @returns {Promise<boolean>}
 */
export const navigateToUnconnectedCategory = async (category_uuid) => {
    /**
     * @type {LabeledError | null}
     */
    let labeled_error = null;

    let tree = get(categories_tree);

    if (tree == null) {
        throw new Error("Cannot use navigateToUnconnectedCategory if categories_tree is null");
    }

    try {
        await tree.changeRootCategory(category_uuid);
    } catch (error) {
        let variable_context = new VariableEnvironmentContextError("In categories_tree.js/navigateToUnconnectedCategory")
        variable_context.addVariable("error", error)
        variable_context.addVariable("category_uuid", category_uuid)

        labeled_error = new LabeledError(variable_context, "The server seems to think this category does not exist.", lf_errors.ERR_PROCESSING_ERROR);

        labeled_error.alert();

        return false;
    }

    // @ts-ignore
    replaceState(`/dungeon-explorer/${category_uuid}`)

    return true;
}
