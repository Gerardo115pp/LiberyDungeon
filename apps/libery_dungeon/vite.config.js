import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import path from 'path';
import fs from 'fs';

const ENABLE_FULL_RELOAD = true;

/**
 * @type {import('vite').PluginOption}
 */
const fullReloadAlways = {
	handleHotUpdate(hmr_context) {
		console.log(`Full reload requested for '${hmr_context.file}'`);
		if (hmr_context.file.includes('libs/LiberyHotkeys')) {
			hmr_context.server.ws.send({
				type: 'full-reload'
			});
		}
		return [];
	}
}

export default defineConfig(async ({ command, mode, isSsrBuild, isPreview }) => {
	let is_production = command === 'build';

	let build_config = {
		PW_SERVER: JSON.stringify(process.env.PW_SERVER),
		BASE_DOMAIN: JSON.stringify(process.env.BASE_DOMAIN),
		JD_SERVER: JSON.stringify(process.env.JD_SERVER),
		CATEGORIES_SERVER: JSON.stringify(process.env.CATEGORIES_SERVER),
		MEDIAS_SERVER: JSON.stringify(process.env.MEDIAS_SERVER),
		COLLECT_SERVER: JSON.stringify(process.env.COLLECT_SERVER),
		DOWNLOADS_SERVER: JSON.stringify(process.env.DOWNLOADS_SERVER),
		METADATA_SERVER: JSON.stringify(process.env.METADATA_SERVER),
		USERS_SERVER: JSON.stringify(process.env.USERS_SERVER),
	}

	/**
	 * @type {import('vite').UserConfig}
	 */
	let config = {
		server: {
			open: false, 
			host: "0.0.0.0",
			port: 5005,
		},
		define: {
			...build_config,
		},
		resolve: {
			alias: {
				'@errors': path.resolve(__dirname, 'src/errors'),
				'@libs': path.resolve(__dirname, 'src/libs'),
				'@databases': path.resolve(__dirname, 'src/databases'),
				'@components': path.resolve(__dirname, 'src/components'),
				'@pages': path.resolve(__dirname, 'src/pages'),
				'@svg': path.resolve(__dirname, 'src/svg'),
				'@models': path.resolve(__dirname, 'src/models'),
				"@actions": path.resolve(__dirname, 'src/actions'),
				"@events": path.resolve(__dirname, 'src/events'),
				"@stores": path.resolve(__dirname, 'src/stores'),
				"@app": path.resolve(__dirname, 'src'),
				"@config": path.resolve(__dirname, 'src/config'),
			}
		},
		plugins: [
			sveltekit()
		],
		clearScreen: true,
	}

	if (ENABLE_FULL_RELOAD) {
		config.plugins.push(fullReloadAlways);
	}

	if (!is_production) {
		config.server.https = {
			key: fs.readFileSync(process.env.SSL_KEY_PATH),
			cert: fs.readFileSync(process.env.SSL_CERT_PATH)
		}
	}
	
	return config;
});
