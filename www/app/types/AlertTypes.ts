/// <reference path="../References.d.ts"/>
export const CREATE = Symbol('alert.create');
export const REMOVE = Symbol('alert.remove');

export const INFO = Symbol('alert.info');
export const WARNING = Symbol('alert.warning');
export const ERROR = Symbol('alert.error');

export interface Alert {
	id?: string;
	level?: Symbol;
	message?: string;
}

export type Alerts = Alert[];

export interface AlertDispatch {
	type: Symbol;
	data?: {
		id?: string;
		level?: Symbol;
		message?: string;
	};
}
