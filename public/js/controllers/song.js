import { Controller } from 'https://unpkg.com/@hotwired/stimulus@3.2.2/dist/stimulus.js';

export default class extends Controller {
	static targets = ['container'];

	static values = {
		streamer: String,
	};

	pullLoop = null;

	connect() {
		this.socket = this.#newSocket();
		window.Turbo.session.connectStreamSource(this.socket);
		this.socket.addEventListener('open', this.#socketConnected.bind(this));
	}

	disconnect() {
		if (this.pullLoop) {
			clearInterval(this.pullLoop);
		}
		this.socket.close();
	}

	#socketConnected() {
		const observer = new IntersectionObserver(this.#pullSongs.bind(this), {
			root: null, // Use viewport as the root
			rootMargin: '0px', // No margin
			threshold: 1.0, // Trigger when 100% visible
		});

		observer.observe(this.containerTarget);
		this.pullLoop = setInterval(() => {
			if (this.#isScrollable()) {
				clearInterval(this.pullLoop);
				return;
			}

			this.#pullSongs();
		}, 5_000);

		// Add an event listener for messages from the WebSocket
		this.socket.addEventListener('message', () => {
			var loadingText = document.getElementById('loading-text');
			if (loadingText) {
				loadingText.remove();
			}
		});
	}

	#pullSongs() {
		console.debug('Pulling...');
		return this.socket.send('PULL');
	}

	#isScrollable() {
		return (
			this.containerTarget.scrollHeight > this.containerTarget.clientHeight
		);
	}

	#newSocket() {
		const proto = window.location.protocol === 'http:' ? 'ws' : 'wss';
		return new WebSocket(
			`${proto}://${window.location.host}/watch/${this.streamerValue}/ws`
		);
	}
}
