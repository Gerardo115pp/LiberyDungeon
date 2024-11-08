/**
 * Returns the Jaro similarity between two strings. The Jaro similarity is a measure of how many edits are required to make two strings equal.
 * @see https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance#Jaro_similarity
 * @param {string} s1
 * @param {string} s2
 * @returns {number}
 */
export function jaro(s1, s2) {
	if (s1 == s2) return 1.0;

	let s1_length = s1.length;
	let s2_length = s2.length;

	const max_matching_character_distance = Math.floor(Math.max(s1_length, s2_length) / 2) - 1;

	let match_count = 0;

	const matches_s1 = Array(s1.length).fill(0);
	const matches_s2 = Array(s1.length).fill(0);

	for (let h = 0; h < s1_length; h++) {

		for (let k = Math.max(0, h - max_matching_character_distance); k < Math.min(s2_length, h + max_matching_character_distance + 1); k++) {
            
            if (s1[h] == s2[k] && matches_s2[k] == 0) {
                matches_s1[h] = 1;
                matches_s2[k] = 1;
                match_count++;
                break;
            }
        }

	}

	if (match_count == 0) return 0.0;

	let transpositions = 0;

	let s2_iterator = 0;

	// Count number of occurrences where two characters match but there is a third matched character in between the indices
	for (let h = 0; h < s1_length; h++) {
		if (matches_s1[h]) {

			// Find the next matched character
			// in second string
			while (matches_s2[s2_iterator] == 0) {
				s2_iterator++;
            }

			if (s1[h] != s2[s2_iterator++]) {
				transpositions++;
            }
		}
    }

	transpositions = transpositions / 2;

	return (1/3) * ((match_count / s1_length) + (match_count / s2_length) + ((match_count - transpositions) / (match_count)));
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