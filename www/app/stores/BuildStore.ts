/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from 'events';
import * as BuildTypes from '../types/BuildTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class _BuildStore extends EventEmitter {
	_state: BuildTypes.Builds = {};
	token: string;

	get builds(): BuildTypes.Builds {
		return this._state;
	}

	emitChange(): void {
		this.emit(GlobalTypes.CHANGE)
	}

	addChangeListener(callback: () => void): void {
		this.on(GlobalTypes.CHANGE, callback);
	}

	removeChangeListener(callback: () => void): void {
		this.removeListener(GlobalTypes.CHANGE, callback);
	}
}
let BuildStore = new _BuildStore();
export default BuildStore;

function loading(): void {
	BuildStore._state = {
		'loading': {
			'id': 'loading',
		},
	};
	BuildStore.emitChange();
}

function load(data: BuildTypes.Build[]): void {
	BuildStore._state = {};
	for (let item of data) {
		BuildStore._state[item.id] = item;
	}
	BuildStore.emitChange();
}

function remove(id: string): void {
	delete BuildStore._state[id];
	BuildStore.emitChange();
}

BuildStore.token = Dispatcher.register(function(
		action: BuildTypes.BuildDispatch): void {
	switch (action.type) {
		case BuildTypes.LOADING:
			loading();
			break;

		case BuildTypes.LOAD:
			load(action.data.builds);
			break;

		case BuildTypes.REMOVE:
			remove(action.data.id);
			break;
	}
});
