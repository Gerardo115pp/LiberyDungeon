import adapter from '@sveltejs/adapter-auto';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	onwarn: (warning, handler) => {
        if (warning.code.startsWith('a11y-')) return;
        handler(warning);
    },
	kit: {
		// adapter-auto only supports some environments, see https://kit.svelte.dev/docs/adapter-auto for a list.
		// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
		// See https://kit.svelte.dev/docs/adapters for more information about adapters.
		adapter: adapter(),
		alias: {
			"@errors/*": "./src/errors/*",
			"@libs/*": "./src/libs/*",
			"@components/*": "./src/components/*",
			"@models/*": "./src/models/*",
			"@pages/*": "./src/pages/*",
			"@svg/*": "./src/svg/*",
			"@actions/*": "./src/actions/*",
			"@events/*": "./src/events/*",
			"@stores/*": "./src/stores/*",
			"@databases/*": "./src/databases/*",
			"@themes/*": "./src/themes/*",
			"@app/*": "./src/*",
			"@config/*": "./src/config/*",
		}
	}
};

export default config;
