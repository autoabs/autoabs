/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as AlertTypes from '../types/AlertTypes';

export function create(level: Symbol, message: string): void {
	Dispatcher.dispatch({
		type: AlertTypes.CREATE,
		data: {
			level: level,
			message: message,
		},
	});
}

export function remove(id: string): void {
	Dispatcher.dispatch({
		type: AlertTypes.REMOVE,
		data: {
			id: id,
		},
	});
}
