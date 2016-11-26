/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as Events from 'events';
import * as BuildTypes from '../types/BuildTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class BuildStore extends Events.EventEmitter {
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

let buildStore = new BuildStore();
export default buildStore;

function loading(): void {
	buildStore._state = {
		'loading': {
			'id': 'loading',
		},
	};
	buildStore.emitChange();
}

function load(data: BuildTypes.Build[]): void {
	buildStore._state = {};
	for (let item of data) {
		buildStore._state[item.id] = item;
	}
	buildStore.emitChange();
}

function remove(id: string): void {
	delete buildStore._state[id];
	buildStore.emitChange();
}

buildStore.token = Dispatcher.register(function(
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
