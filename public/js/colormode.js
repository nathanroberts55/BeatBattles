/*!
 * Color mode toggler for Bootstrap's docs (https://getbootstrap.com/)
 * Copyright 2011-2024 The Bootstrap Authors
 * Licensed under the Creative Commons Attribution 3.0 Unported License.
 */

(() => {
	'use strict';

	const getStoredTheme = () => localStorage.getItem('theme');
	const setStoredTheme = (theme) => localStorage.setItem('theme', theme);

	const getPreferredTheme = () => {
		const storedTheme = getStoredTheme();
		if (storedTheme) {
			return storedTheme;
		}

		return window.matchMedia('(prefers-color-scheme: dark)').matches
			? 'dark'
			: 'light';
	};

	const setTheme = (theme) => {
		document.documentElement.setAttribute('data-bs-theme', theme);
	};

	setTheme(getPreferredTheme());

	window
		.matchMedia('(prefers-color-scheme: dark)')
		.addEventListener('change', () => {
			const storedTheme = getStoredTheme();
			if (storedTheme !== 'light' && storedTheme !== 'dark') {
				setTheme(getPreferredTheme());
			}
		});

	window.addEventListener('DOMContentLoaded', () => {
		const colorModeToggle = document.querySelector('#colormode-toggle');
		colorModeToggle.addEventListener('click', (event) => {
			event.preventDefault();
			const currentTheme = getPreferredTheme();
			const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
			setStoredTheme(newTheme);
			setTheme(newTheme);
		});
	});
})();
