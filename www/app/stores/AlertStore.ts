/// <reference path="../References.d.ts"/>
import * as Events from 'events';
import Dispatcher from '../dispatcher/Dispatcher';
import * as AlertTypes from '../types/AlertTypes';
import * as GlobalTypes from '../types/GlobalTypes';
import * as MiscUtils from '../utils/MiscUtils';

class AlertStore extends Events.EventEmitter {
	_state: AlertTypes.Alerts = [];
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get alerts(): AlertTypes.Alerts {
		return this._state;
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

	_create(level: Symbol, message: string): void {
		let id = MiscUtils.uuid();
		this._map[id] = this._state.push({
			id: id,
			level: level,
			message: message,
		}) - 1;
		this.emitChange();
	}

	_remove(id: string): void {
		let n = this._map[id];
		if (n === undefined) {
			return;
		}
		delete this._map[id];

		this._state.splice(n, 1);

		for (let i = n; i < this._state.length; i++) {
			this._map[this._state[i].id] = i;
		}

		this.emitChange();
	}

	_callback(action: AlertTypes.AlertDispatch): void {
		switch (action.type) {
			case AlertTypes.CREATE:
				this._create(action.data.level, action.data.message);
				break;

			case AlertTypes.REMOVE:
				this._remove(action.data.id);
				break;
		}
	}
}

export default new AlertStore();
