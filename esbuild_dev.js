import esbuild from 'esbuild';
import { sassPlugin } from 'esbuild-sass-plugin';

const context = await esbuild.context({
	entryPoints: ['public/scss/custom.scss'],
	outdir: 'public/css',
	bundle: true,
	minify: true,
	plugins: [sassPlugin()],
});

// Enable watch mode
await context.watch();
