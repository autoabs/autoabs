/// <reference path="./References.d.ts"/>
import Dispatcher from './dispatcher/Dispatcher';

let connected = false;

function connect(): void {
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
		Dispatcher.dispatch(JSON.parse(evt.data).data);
	})
}

export function init() {
	if (connected) {
		return;
	}
	connected = true;

	connect();
}
