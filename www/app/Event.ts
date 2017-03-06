/// <reference path="./References.d.ts"/>
import Dispatcher from './dispatcher/Dispatcher';

let connected = false;

function _connect(): void {
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
			_connect();
		}, 500);
	});

	socket.addEventListener('message', (evt) => {
		Dispatcher.dispatch(JSON.parse(evt.data).data);
	})
}

export function connect() {
	if (connected) {
		return;
	}
	connected = true;

	_connect();
}
