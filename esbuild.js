import esbuild from 'esbuild';
import { sassPlugin } from 'esbuild-sass-plugin';

if (process.env['NODE_ENV'] !== 'production') {
	await esbuild.build({
		entryPoints: ['public/scss/custom.scss'],
		outdir: 'public/css',
		bundle: true,
		minify: true,
		plugins: [sassPlugin()],
	});
}
