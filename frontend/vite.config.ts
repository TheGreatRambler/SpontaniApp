import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
        ssr: {
                noExternal: ['@googlemaps/js-api-loader'],
        },
        plugins: [sveltekit()],
        server: {
                fs: {
                        // allow serving assets in the static directory
                        allow: ['..'],
                }
        }
});
