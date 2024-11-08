/**
 * Returns the Jaro similarity between two strings. The Jaro similarity is a measure of how many edits are required to make two strings equal.
 * @see https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance#Jaro_similarity
 * @param {string} s1
 * @param {string} s2
 * @returns {number}
 */
export function jaro(s1, s2) {
    let s1_length = s1.length;
	let s2_length = s2.length;

    if (s1_length == 0 && s2_length == 0) return 1; // If both strings are zero-length, they are completely equal.

    if (s1_length == 0 || s2_length == 0) return 0; // If one of the strings is zero-length, they are completely unequal.

    let match_distance = Math.max(0, Math.floor(Math.max(s1_length, s2_length) / 2) - 1);

    const matches_s1 = new Array(s1_length).fill(false);
    const matches_s2 = new Array(s2_length).fill(false);

    let matching_characters = 0;

    // Step 1: Matching characters - Find all characters in the first string that match a character in the second string.
    for (let h = 0; h < s1_length; h++) {
        let start = Math.max(0, h - match_distance);
        let end = Math.min(s2_length - 1, h + match_distance);

        for (let k = start; k <= end; k++) {
            if (matches_s2[k]) continue;

            if (s1[h] == s2[k]) {
                matches_s1[h] = true;
                matches_s2[k] = true;
                matching_characters++;
                break;
            }
        }
    }

    if (matching_characters == 0) return 0; // If there are no matching characters, the strings are completely unequal.

    // Step 2: Transpositions - count the number of unaligned matches.
    let unaligned_matches = 0;
    let matches_s2_iterator = 0;

    for (let h = 0; h < s1_length; h++) {
        if (!matches_s1[h]) continue;

        while (!matches_s2[matches_s2_iterator]) matches_s2_iterator++;

        if (s1[h] != s2[matches_s2_iterator]) {
            unaligned_matches++;
        }

        matches_s2_iterator++;
    }

    const transpositions = Math.floor(unaligned_matches / 2);

    return (1 / 3) * ((matching_characters / s1_length) + (matching_characters / s2_length) + ((matching_characters - transpositions) / matching_characters));
}

/**
 * Returns the Jaro-Winkler similarity between two strings. Like Jaro but boosts scores for strings that share common prefixes.
 * @see https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance#Jaro%E2%80%93Winkler_similarity
 * @param {string} s1
 * @param {string} s2
 * @returns {number}
 */
export function jaroWinkler(s1, s2) {
    const boost_threshold = 0.7; 
    const prefix_scale = 0.1; // Scaling factor for how much the score is adjusted for having common prefixes. Recommended value is 0.1
    let prefix_length = 4; // The number of characters at the start of the string to consider for the common prefix adjustment up to a maximum of 4 characters. Recommended value is 4

    let jaro_similarity = jaro(s1, s2);

    if (jaro_similarity <= boost_threshold) return jaro_similarity;

    prefix_length = Math.min(s1.length, Math.min(prefix_length, s2.length));

    let prefix_match = 0;

    for (let h = 0; h < prefix_length; h++) {
        if (s1[h] == s2[h]) {
            prefix_match++;
        } else {
            break;
        }
    }

    return jaro_similarity + prefix_scale * prefix_match * (1 - jaro_similarity);
}