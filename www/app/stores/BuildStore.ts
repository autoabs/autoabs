/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as Events from 'events';
import * as BuildTypes from '../types/BuildTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class BuildStore extends Events.EventEmitter {
	_state: BuildTypes.Builds = [];
	_loadingState: boolean;
	_token = Dispatcher.register((this._callback).bind(this));

	get builds(): BuildTypes.Builds {
		return this._state;
	}

	get loading(): boolean {
		return this._loadingState;
	}

	emitChange(): void {
		this.emit(GlobalTypes.CHANGE);
	}

	addChangeListener(callback: () => void): void {
		this.on(GlobalTypes.CHANGE, callback);
	}

	removeChangeListener(callback: () => void): void {
		this.removeListener(GlobalTypes.CHANGE, callback);
	}

	_loading(): void {
		if (this._loadingState !== true) {
			this._loadingState = true;
			this.emitChange();
		}
	}

	_loaded(): void {
		if (this._loadingState !== false) {
			this._loadingState = false;
			this.emitChange();
		}
	}

	_sync(data: BuildTypes.Build[]): void {
		this._state = data;
		this.emitChange();
	}

	_remove(id: string): void {
		for (let i = 0; i < this._state.length; i++) {
			if (this._state[i].id === id) {
				this._state.splice(i, 1);
				break;
			}
		}
		this.emitChange();
	}

	_callback(action: BuildTypes.BuildDispatch): void {
		switch (action.type) {
			case BuildTypes.LOADING:
				this._loading();
				break;

			case BuildTypes.LOADED:
				this._loaded();
				break;

			case BuildTypes.SYNC:
				this._sync(action.data.builds);
				break;

			case BuildTypes.REMOVE:
				this._remove(action.data.id);
				break;
		}
	}
}

export default new BuildStore();
