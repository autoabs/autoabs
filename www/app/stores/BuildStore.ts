/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as Events from 'events';
import * as BuildTypes from '../types/BuildTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class BuildStore extends Events.EventEmitter {
	_builds: BuildTypes.Builds = [];
	_index: number;
	_count: number;
	_map: {[key: string]: number} = {};
	_loadingState: boolean;
	_token = Dispatcher.register((this._callback).bind(this));

	get builds(): BuildTypes.Builds {
		return this._builds;
	}

	get index(): number {
		return this._index || 0;
	}

	get count(): number {
		return this._count;
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

	_update(data: BuildTypes.Build): void {
		let n = this._map[data.id];
		if (n === undefined) {
			return;
		}
		this._builds[n] = data;
		this.emitChange();
	}

	_sync(builds: BuildTypes.Builds, index: number, count: number): void {
		this._index = index;
		this._count = count;

		this._map = {};
		for (let i = 0; i < builds.length; i++) {
			this._map[builds[i].id] = i;
		}
		this._builds = builds;

		this.emitChange();
	}

	_remove(id: string): void {
		let n = this._map[id];
		if (n === undefined) {
			return;
		}
		delete this._map[id];

		this._builds.splice(n, 1);

		for (let i = n; i < this._builds.length; i++) {
			this._map[this._builds[i].id] = i;
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

			case BuildTypes.UPDATE:
				this._update(action.data.build);
				break;

			case BuildTypes.SYNC:
				this._sync(action.data.builds, action.data.index, action.data.count);
				break;

			case BuildTypes.REMOVE:
				this._remove(action.data.id);
				break;
		}
	}
}

export default new BuildStore();
