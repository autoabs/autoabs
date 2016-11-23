/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as BuildTypes from '../types/BuildTypes';

export function loading(): void {
	Dispatcher.dispatch({
		type: BuildTypes.LOADING,
	});
}

export function load(builds: BuildTypes.Build[]): void {
	Dispatcher.dispatch({
		type: BuildTypes.LOAD,
		data: {
			builds: builds,
		},
	});
}

export function remove(id: string): void {
	Dispatcher.dispatch({
		type: BuildTypes.REMOVE,
		data: {
			id: id,
		},
	});
}
