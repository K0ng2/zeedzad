import tailwindcss from '@tailwindcss/vite'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

// Read package.json to get version
const packageJson = JSON.parse(readFileSync(resolve(__dirname, 'package.json'), 'utf-8'))

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	compatibilityDate: '2025-07-15',
	devtools: { enabled: true },
	css: [
		'@/assets/css/main.css',
		'@fortawesome/fontawesome-svg-core/styles.css',
	],
	ssr: false,
	vite: {
		define: {
			__APP_VERSION__: JSON.stringify(packageJson.version),
		},
		plugins: [
			tailwindcss(),
		],
	},
	app: {
		head: {
			title: 'Zeedzad',
			meta: [
				{ name: 'description', content: 'Zeedzad - Video Game Database' },
				{ name: 'viewport', content: 'width=device-width, initial-scale=1' },
				{ name: 'theme-color', content: '#ffffff' },
				{ charset: 'utf-8' },
			],
		}
	},
	pwa: {
		manifest: {
			name: 'GDB',
			short_name: 'GDB',
			start_url: '/',
			display: 'standalone',
			background_color: '#ffffff',
			orientation: "portrait",
			description: "GDB",
			scope: '/',
			icons: [
				{
					src: "favicon.webp",
					sizes: "256x256",
					type: 'image/webp'
				}
			]
		},
		workbox: {
			globPatterns: ['**/*.{js,css,html,svg,png,ico}'],
			cleanupOutdatedCaches: true,
			clientsClaim: true,
			maximumFileSizeToCacheInBytes: 5 * 1024 * 1024, // 5MB
		},
		registerType: 'autoUpdate',
	},
	runtimeConfig: {
		API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:8088',
	},
	modules: [
		'@vite-pwa/nuxt',
	],
})
