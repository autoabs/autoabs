/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as Events from 'events';
import * as BuildInfoTypes from '../types/BuildInfoTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class BuildInfoStore extends Events.EventEmitter {
	_id: string = '';
	_output: string = '';
	_token = Dispatcher.register((this._callback).bind(this));

	get id(): string {
		return this._id;
	}

	get output(): string {
		return this._output;
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

	_open(id: string): void {
		this._id = id;
		this._output = '';
		this.emitChange();
	}

	_close(): void {
		this._id = null;
		this._output = '';
		this.emitChange();
	}

	_update(output: string[]): void {
		this._output = output.join('\n');
		this.emitChange();
	}

	_callback(action: BuildInfoTypes.BuildInfoDispatch): void {
		switch (action.type) {
			case BuildInfoTypes.OPEN:
				this._open(action.data.id);
				break;

			case BuildInfoTypes.CLOSE:
				this._close();
				break;

			case BuildInfoTypes.UPDATE:
				this._update(action.data.output);
				break;
		}
	}
}

export default new BuildInfoStore();
