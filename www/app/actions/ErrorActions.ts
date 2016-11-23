/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as ErrorTypes from '../types/ErrorTypes';

export function create(level: Symbol, message: string): void {
	Dispatcher.dispatch({
		type: ErrorTypes.CREATE,
		data: {
			level: level,
			message: message,
		},
	});
}

export function remove(id: string): void {
	Dispatcher.dispatch({
		type: ErrorTypes.REMOVE,
		data: {
			id: id,
		},
	});
}
