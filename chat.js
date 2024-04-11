#!/usr/bin/env node

import tmi from 'tmi.js';

const client = new tmi.Client({
	channels: ['aspecticor'],
});

client.connect();

client.on('message', (channel, tags, message, self) => {
	console.log(`message:${tags['display-name']}:${message}`);
});

client.on('cheer', (channel, userstate, message) => {
	console.log(`bits:${userstate['display-name']}:${userstate.bits}:${message}`);
});
