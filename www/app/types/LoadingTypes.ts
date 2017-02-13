/// <reference path="../References.d.ts"/>
export const ADD = Symbol('loading.add');
export const DONE = Symbol('loading.done');

export interface LoadingDispatch {
	type: Symbol;
	data?: {
		id?: string;
	};
}
