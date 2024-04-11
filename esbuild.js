import esbuild from 'esbuild';
import { sassPlugin } from 'esbuild-sass-plugin';

const context = await esbuild.context({
	entryPoints: ['frontend/src/Index.tsx', 'frontend/static/scss/App.scss'],
	outdir: 'public/assets',
	bundle: true,
	minify: true,
	plugins: [sassPlugin()],
});

// Enable watch mode
await context.watch();
