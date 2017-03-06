/// <reference path="./References.d.ts"/>

export function connect(): void {
	let url = '';
	let location = window.location;

	if (location.protocol === 'https:') {
		url += 'wss';
	} else {
		url += 'ws';
	}

	url += '://' + location.host + '/event';

	let socket = new WebSocket(url);

	socket.addEventListener('close', () => {
		setTimeout(() => {
			connect();
		}, 500);
	});

	socket.addEventListener('message', (evt) => {
		console.log(evt);
	})
}
